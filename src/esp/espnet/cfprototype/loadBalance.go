package cfprototype

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"esp/espnet"
	"fmt"
	"logger"
	"sync"
	"sync/atomic"
	"time"
	"uprop"
)

// Config
type LoadBalanceNodeConfig struct {
	ChannelName string
	Priority    int
	FailOver    bool
}
type LoadBalanceConfig struct {
	Nodes          []*LoadBalanceNodeConfig
	FailSkipTimeMS int
}

func (this *LoadBalanceConfig) Valid() error {
	if this.Nodes == nil || len(this.Nodes) == 0 {
		return errors.New("channels empty")
	}
	if this.Nodes != nil {
		for _, node := range this.Nodes {
			if node.ChannelName == "" {
				return errors.New("node empty name")
			}
			if node.Priority <= 0 {
				node.Priority = 1
			}
		}
	}
	return nil
}

func (this *LoadBalanceConfig) AddName(n string) *LoadBalanceNodeConfig {
	node := new(LoadBalanceNodeConfig)
	node.ChannelName = n
	this.AddNode(node)
	return node
}

func (this *LoadBalanceConfig) AddNode(n *LoadBalanceNodeConfig) {
	if this.Nodes == nil {
		this.Nodes = make([]*LoadBalanceNodeConfig, 0)
	}
	this.Nodes = append(this.Nodes, n)
}

func (this *LoadBalanceConfig) Remove(n string) bool {
	if this.Nodes != nil {
		for i, node := range this.Nodes {
			if node.ChannelName == n {
				c := len(this.Nodes)
				copy(this.Nodes[i:c-1], this.Nodes[i+1:c])
				this.Nodes[c-1] = nil
				this.Nodes = this.Nodes[:c-1]
				return true
			}
		}
	}
	return false
}

func (this *LoadBalanceConfig) GetProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)

	r = append(r, uprop.NewUProperty("failskip", this.FailSkipTimeMS, true, "fail skip duration, MS", func(v string) error {
		this.FailSkipTimeMS = valutil.ToInt(v, this.FailSkipTimeMS)
		return nil
	}))
	if this.Nodes != nil {
		c := len(this.Nodes)
		for i, node := range this.Nodes {
			var n string
			if i == c-1 {
				n = "last"
			} else {
				n = fmt.Sprintf("%d", i)
			}
			r = append(r, uprop.NewUProperty(n+".name", node.ChannelName, false, n+") channel name", func(v string) error {
				node.ChannelName = v
				return nil
			}))
			r = append(r, uprop.NewUProperty(n+".failover", node.FailOver, true, n+") use for fail over", func(v string) error {
				node.FailOver = valutil.ToBool(v, node.FailOver)
				return nil
			}))
			r = append(r, uprop.NewUProperty(n+".priority", node.Priority, true, n+") robin priority, defalut 1", func(v string) error {
				node.Priority = valutil.ToInt(v, 1)
				return nil
			}))
			r = append(r, uprop.NewUProperty(n+".delete", false, true, n+") delete node", func(v string) error {
				if !valutil.ToBool(v, true) {
					this.Remove(node.ChannelName)
				}
				return nil
			}))
		}
	}
	return r
}

// ChannleFactory
type lbItem struct {
	id          int
	pos         int
	channelName string
	factory     espnet.ChannelFactory
	priority    int
	failover    bool
}
type LoadBalanceChannelFactory struct {
	storage ChannelFactoryStorage
	config  *LoadBalanceConfig

	version       uint64
	lock          sync.RWMutex
	allItems      map[int]*lbItem
	items         []*lbItem
	foItems       []*lbItem
	roundRobin    uint32
	totalPriority int
	itemId        int
}

func NewLoadBalanceChannelFactory(storage ChannelFactoryStorage, cfg *LoadBalanceConfig) *LoadBalanceChannelFactory {
	if err := cfg.Valid(); err != nil {
		panic(err.Error())
	}

	this := new(LoadBalanceChannelFactory)
	this.storage = storage
	this.config = cfg
	return this
}

func (this *LoadBalanceChannelFactory) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("LoadBalance[")
	if this.config != nil {
		if this.config.Nodes != nil {
			for i, n := range this.config.Nodes {
				if i > 0 {
					buf.WriteString(",")
				}
				buf.WriteString(n.ChannelName)
				if n.Priority >= 1 {
					buf.WriteString(fmt.Sprintf(":%d", n.Priority))
				}
				if n.FailOver {
					buf.WriteString("*")
				}
			}
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (this *LoadBalanceChannelFactory) LoadItems() ([]*lbItem, []*lbItem) {
	ver := this.storage.GetStorageVersion()

	this.lock.RLock()
	mver := this.version
	items := this.items
	foitems := this.foItems
	this.lock.RUnlock()

	if ver != mver || items == nil {
		this.lock.Lock()
		defer this.lock.Unlock()
		if ver != this.version || this.items == nil {
			flist := make([]*lbItem, 0)
			folist := make([]*lbItem, 0)
			this.allItems = make(map[int]*lbItem)
			this.totalPriority = 0
			if this.config.Nodes != nil {
				for i, n := range this.config.Nodes {
					cf, err := this.storage.GetChannelFactory(n.ChannelName)
					if err != nil {
						logger.Warn(tag, "getChannelFactory(%s) fail - %s", n.ChannelName, err)
						continue
					}
					item := new(lbItem)
					this.itemId++
					item.id = this.itemId
					item.pos = i
					item.channelName = n.ChannelName
					item.priority = n.Priority
					item.factory = cf
					if item.priority <= 0 {
						item.priority = 1
					}
					item.failover = n.FailOver
					if item.failover {
						folist = append(folist, item)
					} else {
						this.totalPriority += item.priority
						flist = append(flist, item)
					}
					this.allItems[item.id] = item
				}
			}
			if this.totalPriority <= 0 {
				this.totalPriority = 1
			}
			this.version = ver
			this.items = flist
			this.foItems = folist
			items = flist
			foitems = folist
		}
	}
	return items, foitems
}

func (this *LoadBalanceChannelFactory) restoreFail(iid int) {
	this.lock.RLock()
	item, ok := this.allItems[iid]
	this.lock.RUnlock()

	if !ok {
		logger.Debug(tag, "miss restore(%d)", iid)
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	item, ok = this.allItems[iid]
	if !ok {
		logger.Debug(tag, "miss2 restore(%d)", iid)
		return
	}
	done := false
	tmp := make([]*lbItem, 0, len(this.items)+1)
	for _, o := range this.items {
		if !done && item.pos < o.pos {
			tmp = append(tmp, item)
			done = true
		}
		tmp = append(tmp, o)
	}
	if !done {
		tmp = append(tmp, item)
	}
	this.totalPriority += item.priority
	this.items = tmp
	logger.Debug(tag, "restore %s", item.factory)
}

func (this *LoadBalanceChannelFactory) removeFail(iid int) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.items == nil || len(this.items) == 0 {
		return
	}
	done := false
	tmp := make([]*lbItem, 0, len(this.items)-1)
	for _, item := range this.items {
		if item.id == iid {
			this.totalPriority -= item.priority
			done = true
			logger.Debug(tag, "temp remove %s", item.factory)
			continue
		}
		tmp = append(tmp, item)
	}
	this.items = tmp

	if done {
		sktime := this.config.FailSkipTimeMS
		if sktime <= 0 {
			sktime = 5000
		}
		time.AfterFunc(time.Duration(sktime)*time.Millisecond, func() {
			this.restoreFail(iid)
		})
	}
}

func (this *LoadBalanceChannelFactory) proccess(ilist []*lbItem, ppos, dpos uint32) (espnet.Channel, bool) {
	c := len(ilist)
	for i := 0; i < c; i++ {
		pos := (int(ppos) + i) % c
		item := ilist[pos]
		f := item.factory
		if fb, ok := f.(espnet.ChannelFactoryBreakSupport); ok {
			bv := fb.IsBreak()
			if bv != nil && *bv {
				logger.Debug(tag, "%s break, skip", f)
				this.removeFail(item.id)
				continue
			}
		}
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "round %d => %s", dpos, f)
		}
		ch, err := f.NewChannel()
		if err != nil {
			logger.Warn(tag, "%s NewChannel fail %s", f, err)
			go this.removeFail(item.id)
			continue
		}
		return ch, true
	}
	return nil, false
}

func (this *LoadBalanceChannelFactory) NewChannel() (espnet.Channel, error) {
	ilist, folist := this.LoadItems()

	v := atomic.AddUint32(&this.roundRobin, 1)
	tp := this.totalPriority
	if tp <= 0 {
		tp = 1
	}
	pri := v % uint32(tp)
	pri2 := pri
	ppos := uint32(0)
	for i, item := range ilist {
		ipri := uint32(item.priority)
		if ipri > pri {
			ppos = uint32(i)
			break
		}
		pri -= ipri
	}
	ch, ok := this.proccess(ilist, ppos, pri2)
	if ok {
		return ch, nil
	}

	c := uint32(len(folist))
	if c > 0 {
		logger.Debug(tag, "process failover")
		ppos = v % c
		ch, ok = this.proccess(folist, ppos, ppos)
		if ok {
			return ch, nil
		}
	}
	return nil, errors.New("channels break")
}

func (this *LoadBalanceChannelFactory) IsBreak() *bool {
	ilist, _ := this.LoadItems()
	unknow := false
	for _, item := range ilist {
		f := item.factory
		if fb, ok := f.(espnet.ChannelFactoryBreakSupport); ok {
			vp := fb.IsBreak()
			if vp != nil {
				if *vp {
					r := false
					return &r
				}
			} else {
				unknow = true
			}
		} else {
			unknow = true
		}
	}
	if unknow {
		return nil
	}
	r := true
	return &r
}

// Prototype
type LoadBalancePrototype struct {
	config *LoadBalanceConfig
}

func (this *LoadBalancePrototype) Valid() error {
	if this.config == nil {
		return errors.New("config nil")
	}
	return this.config.Valid()
}

func (this *LoadBalancePrototype) ToMap() map[string]interface{} {
	if this.config != nil {
		m := valutil.BeanToMap(this.config)
		return m
	}
	return nil
}

func (this *LoadBalancePrototype) FromMap(data map[string]interface{}) error {
	if data != nil {
		cfg := new(LoadBalanceConfig)
		valutil.ToBean(data, cfg)
		this.config = cfg

		if err := this.Valid(); err != nil {
			return err
		}
	}
	return nil
}

func (this *LoadBalancePrototype) GetProperties() []*uprop.UProperty {
	// r := make([]*uprop.UProperty, 0)
	if this.config == nil {
		this.config = new(LoadBalanceConfig)
	}
	return this.config.GetProperties()
}

func (this *LoadBalancePrototype) CreateChannelFactory(storage ChannelFactoryStorage, name string, start bool) (espnet.ChannelFactory, error) {
	if err := this.Valid(); err != nil {
		return nil, err
	}
	return NewLoadBalanceChannelFactory(storage, this.config), nil
}
