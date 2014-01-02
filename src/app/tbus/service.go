package tbus

import (
	"bmautil/socket"
	"bmautil/valutil"
	"boot"
	"config"
	"esp/espnet"
	"esp/espnet/cfprototype"
	"esp/sqlite"
	"fmt"
	"logger"
	"strings"
	"sync"
)

const (
	tag       = "TBus"
	tableName = "tbl_tbus_service"
)

type B3V byte

const (
	B3V_UNKNOW = B3V(0)
	B3V_TRUE   = B3V(1)
	B3V_FALSE  = B3V(2)
)

// ThriftServiceInfo
type ThriftMethodInfo struct {
	Name   string
	Oneway B3V
	Query  B3V
}

type ThriftServiceInfo struct {
	Methods map[string]*ThriftMethodInfo
	Remote  string
}

func (this *ThriftServiceInfo) Clone() *ThriftServiceInfo {
	r := new(ThriftServiceInfo)
	r.Remote = this.Remote
	r.Methods = make(map[string]*ThriftMethodInfo)
	for k, m := range this.Methods {
		r.Methods[k] = m
	}
	return r
}

type serviceInfo struct {
	Name   string
	Config *ThriftServiceInfo
	count  int
}

type remoteInfo struct {
	name      string
	kind      string
	prototype cfprototype.ChannelFactoryPrototype
	factory   espnet.ChannelFactory
}

// TBusService
type TBusService struct {
	name     string
	database *sqlite.SqliteServer
	config   configInfo

	lock    sync.RWMutex
	infos   map[string]*serviceInfo
	methods map[string]string
	remotes map[string]*remoteInfo
	version uint64
}

func NewTBusService(name string, db *sqlite.SqliteServer) *TBusService {
	this := new(TBusService)
	this.name = name
	this.database = db
	this.infos = make(map[string]*serviceInfo)
	this.methods = make(map[string]string)
	this.remotes = make(map[string]*remoteInfo)
	this.initDatabase()
	return this
}

func (this *TBusService) Name() string {
	return this.name
}

type configInfo struct {
	SafeMode      bool
	MaxFrame      string
	maxFrame      int
	Trace         int
	DefaultRemote string
	AdminWord     string
}

func (this *TBusService) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.MaxFrame != "" {
			mf, _ := valutil.ToSize(cfg.MaxFrame, 1024, valutil.SizeB)
			cfg.maxFrame = int(mf)
		}
		this.config = cfg
	}
	if this.config.maxFrame == 0 {
		this.config.maxFrame = 10 * 1024 * 1024
	}
	espnet.RegSocketChannelCoder("tbus", func() espnet.SocketChannelCoder {
		return NewChannelCoder(this.config.maxFrame)
	})
	return true
}

func (this *TBusService) Start() bool {
	cfg, ok := this.loadRuntimeConfig()
	if !ok {
		if !this.config.SafeMode {
			return false
		}
	}
	if !this.setupByConfig(cfg) {
		return false
	}

	return true
}

func (this *TBusService) Close() bool {
	tmp := func() []espnet.ChannelFactory {
		r := make([]espnet.ChannelFactory, 0)
		this.lock.Lock()
		defer this.lock.Unlock()
		for k, rinfo := range this.remotes {
			delete(this.remotes, k)
			if rinfo.factory != nil {
				r = append(r, rinfo.factory)
			}
		}
		return r
	}()
	for _, fac := range tmp {
		boot.RuntimeStopCloseClean(fac, false)
	}
	return true
}

func (this *TBusService) initDatabase() {
	this.database.InitRuntmeConfigTable(tableName, []int{1})
}

type runtimeConfig struct {
	Services map[string]*ThriftServiceInfo
	Remotes  map[string]map[string]interface{}
}

func (this *TBusService) loadRuntimeConfig() (*runtimeConfig, bool) {
	var cfg runtimeConfig
	err := this.database.LoadRuntimeConfig(tableName, 1, &cfg)
	if err != nil {
		return nil, false
	}
	return &cfg, true
}

func (this *TBusService) setupByConfig(cfg *runtimeConfig) bool {
	if cfg.Remotes != nil {
		for rname, robj := range cfg.Remotes {
			err := func(rname string, robj map[string]interface{}) error {
				kind := valutil.ToString(robj["_kind"], "")
				p := CreateChannelFactoryPrototype(kind)
				if p == nil {
					return logger.Error(tag, "RemoteInfo[%s] prototype invalid", kind)
				}
				err2 := p.FromMap(robj)
				if err2 != nil {
					return logger.Error(tag, "RemoteInfo[%s] config fail - %s", rname, err2)
				}
				err := this.SetRemote(rname, kind, p, false)
				if err != nil {
					return logger.Error(tag, "RemoteInfo[%s] setup fail - %s", rname, err)
				}
				return nil
			}(rname, robj)
			if err != nil {
				if this.config.SafeMode {
					continue
				}
				return false
			}
		}
	}
	if cfg.Services != nil {
		for sname, sobj := range cfg.Services {
			err := this.SetService(sname, sobj)
			if err != nil {
				logger.Error(tag, "ThriftService[%s] setup fail %s", sname, err)
				if this.config.SafeMode {
					continue
				}
				return false
			}
		}
	}
	return true
}

func (this *TBusService) storeRuntimeConfig(cfg *runtimeConfig) error {
	return this.database.StoreRuntimeConfig(tableName, 1, cfg)
}

func (this *TBusService) buildRuntimeConfig() *runtimeConfig {
	r := new(runtimeConfig)
	this.lock.RLock()
	defer this.lock.RUnlock()
	r.Services = make(map[string]*ThriftServiceInfo)
	for sname, sobj := range this.infos {
		r.Services[sname] = sobj.Config
	}
	r.Remotes = make(map[string]map[string]interface{})
	for rname, robj := range this.remotes {
		m := robj.prototype.ToMap()
		if m == nil {
			m = make(map[string]interface{})
		}
		m["_kind"] = robj.kind
		r.Remotes[rname] = m
	}
	return r
}

func (this *TBusService) save() error {
	cfg := this.buildRuntimeConfig()
	return this.storeRuntimeConfig(cfg)
}

func (this *TBusService) GetServiceInfo(name string) *ThriftServiceInfo {
	this.lock.RLock()
	defer this.lock.RUnlock()
	o, ok := this.infos[name]
	if ok {
		return o.Config.Clone()
	}
	return nil
}

func (this *TBusService) SetService(name string, info *ThriftServiceInfo) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	old, ok := this.infos[name]
	if ok {
		mlist := old.Config.Methods
		for k, _ := range mlist {
			if this.methods[k] == name {
				delete(this.methods, k)
			}
		}
	}

	if info != nil {
		sinfo := new(serviceInfo)
		sinfo.Name = name
		sinfo.Config = info
		this.infos[name] = sinfo

		if info.Methods == nil {
			info.Methods = make(map[string]*ThriftMethodInfo)
		}
		mlist := info.Methods
		for k, _ := range mlist {
			this.methods[k] = name
		}
	}
	// fmt.Println(this.methods)

	return nil
}

func (this *TBusService) DeleteService(name string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	old, ok := this.infos[name]
	if ok {
		delete(this.infos, name)
		mlist := old.Config.Methods
		for k, _ := range mlist {
			if this.methods[k] == name {
				delete(this.methods, k)
			}
		}
	}
	return nil
}

func (this *TBusService) GetRemotePrototype(name string) (string, cfprototype.ChannelFactoryPrototype) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	info, ok := this.remotes[name]
	if !ok {
		return "", nil
	}
	return info.kind, info.prototype
}

func (this *TBusService) SetRemote(name string, kind string, p cfprototype.ChannelFactoryPrototype, edit bool) error {
	if edit {
		this.lock.RLock()
		_, old := this.remotes[name]
		this.lock.RUnlock()
		if old {
			this.DeleteRemote(name)
		}
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	_, ok := this.remotes[name]
	if ok {
		return fmt.Errorf("Remote[%s] exists", name)
	}

	r := new(remoteInfo)
	r.name = name
	r.kind = kind
	r.prototype = p
	this.remotes[name] = r
	this.version++

	return nil
}

func (this *TBusService) DeleteRemote(name string) error {
	ri := func() *remoteInfo {
		this.lock.Lock()
		defer this.lock.Unlock()

		ri, ok := this.remotes[name]
		if !ok {
			return nil
		}
		delete(this.remotes, name)
		this.version++
		return ri
	}()

	if ri != nil && ri.factory != nil {
		boot.RuntimeStopCloseClean(ri, true)
	}
	return nil
}

func (this *TBusService) lGetRemote(n string) *remoteInfo {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.remotes[n]
}

func (this *TBusService) GetStorageVersion() uint64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.version
}

func (this *TBusService) GetChannelFactory(n string) (espnet.ChannelFactory, error) {
	r := this.lGetRemote(n)
	if r == nil {
		return nil, fmt.Errorf("not found remote[%s]", n)
	}
	if r.factory != nil {
		return r.factory, nil
	}
	fac, err := r.prototype.CreateChannelFactory(this, n, true)
	if err != nil {
		return nil, err
	}
	rfac, ok := func() (espnet.ChannelFactory, bool) {
		this.lock.Lock()
		defer this.lock.Unlock()
		if r.factory == nil {
			r.factory = fac
			return fac, true
		}
		return r.factory, false
	}()
	if !ok {
		boot.RuntimeStopCloseClean(fac, false)
	}
	return rfac, nil
}

func (this *TBusService) OnSocketAccept(sock *socket.Socket) error {
	ch := espnet.NewSocketChannel(sock, "tbus")
	if this.config.Trace > 0 {
		sock.Trace = this.config.Trace
	}
	NewProxy(ch, this.FindChannel)
	return nil
}

func (this *TBusService) OnOutgoAccept(sock *socket.Socket) error {
	if this.config.Trace > 0 {
		sock.Trace = this.config.Trace
	}
	return nil
}

func (this *TBusService) FindServiceAndMethod(module, method string) (*serviceInfo, *ThriftMethodInfo, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if module == "" {
		module = this.methods[method]
	}
	if module == "" {
		return nil, nil, fmt.Errorf("unknow module for '%s'", method)
	}
	info := this.infos[module]
	if info == nil {
		return nil, nil, fmt.Errorf("module(%s) not config", module)
	}
	var m *ThriftMethodInfo
	if info.Config.Methods != nil {
		m = info.Config.Methods[method]
	}
	if m == nil {
		m = new(ThriftMethodInfo)
		m.Name = method
	}
	return info, m, nil
}

func SplitThriftName(name string) (string, string) {
	np := strings.SplitN(name, ".", 2)
	var module, method string
	if len(np) > 1 {
		module = np[0]
		method = np[1]
	} else {
		module = ""
		method = np[0]
	}
	return module, method
}

func (this *TBusService) FindChannel(name string) (*CFInfo, error) {
	module, method := SplitThriftName(name)

	info, m, err := this.FindServiceAndMethod(module, method)
	rem := ""
	if err != nil {
		if this.config.DefaultRemote == "" {
			return nil, err
		}
		logger.Debug(tag, "%s defaultRemote %s", name, this.config.DefaultRemote)
		rem = this.config.DefaultRemote
		m = new(ThriftMethodInfo)
		m.Name = method
	} else {
		rem = info.Config.Remote
		logger.Debug(tag, "%s remote %s", name, rem)
	}

	cf, err2 := this.GetChannelFactory(rem)
	if err2 != nil {
		return nil, err2
	}
	ch, err3 := cf.NewChannel()
	if err3 != nil {
		return nil, err3
	}

	r := new(CFInfo)
	r.channel = ch
	r.method = m
	return r, nil
}
