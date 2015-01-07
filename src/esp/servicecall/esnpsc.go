package servicecall

import (
	"bmautil/conndialpool"
	"bmautil/valutil"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"objfac"
	"time"
)

type ESNPServiceCaller struct {
	name       string
	pool       *conndialpool.DialPool
	maxPackage int
	timeoutMS  int
}

func (this *ESNPServiceCaller) SetName(n string) {
	this.name = n
	if this.pool != nil {
		this.pool.SetName("serviceCall" + n)
	}
}

func (this *ESNPServiceCaller) Ping() bool {
	if this.pool.GetInitSize() > 0 {
		return this.pool.ActiveConn() > 0
	}
	return true
}

func (this *ESNPServiceCaller) Start() error {
	if !this.pool.StartAndRun() {
		return fmt.Errorf("serviceCall pool start fail")
	}
	return nil
}

func (this *ESNPServiceCaller) Stop() {
	this.pool.Close()
}

func ESNPCall(sock espsocket.Socket, serviceName string, method string, params map[string]interface{}, deadline time.Time) (interface{}, error) {
	msg := esnp.NewRequestMessageWithId()
	addr := msg.GetAddress()
	addr.SetCall(serviceName, method)
	dt := msg.Datas()
	dt.Set("p", params)

	ts := time.Now()
	rmsg, err1 := espsocket.CallTimeout(sock, msg, deadline)
	te := time.Now()
	if err1 != nil {
		if rmsg == nil {
			sock.AskClose()
		}
		logger.Debug(tag, "[%s:%s] esnp call(%f) fail '%s'", serviceName, method, te.Sub(ts).Seconds(), err1)
		return nil, err1
	}
	dt = rmsg.Datas()
	st, errX1 := dt.GetInt("s", 0)
	if errX1 != nil {
		return nil, errX1
	}
	val, errX2 := dt.Get("r")
	if errX2 != nil {
		return nil, errX2
	}

	logger.Debug(tag, "[%s:%s] esnp call(%f) end '%d'", serviceName, method, te.Sub(ts).Seconds(), st)
	if st != 200 {
		msg, _ := dt.GetString("m", "")
		if msg == "" {
			msg = fmt.Sprintf("invalid esnp status(%d)", st)
		}
		return nil, fmt.Errorf(msg)
	}
	return val, nil
}

func (this *ESNPServiceCaller) Call(serviceName, method string, params map[string]interface{}, deadline time.Time) (interface{}, error) {
	tm := this.timeoutMS
	if tm <= 0 {
		tm = 5000
	}
	tmd := time.Now().Add(time.Duration(tm) * time.Millisecond)
	if !deadline.IsZero() && tmd.After(deadline) {
		tmd = deadline
	}
	conn, err0 := this.pool.GetConn(tmd, true)
	if err0 != nil {
		return nil, err0
	}
	sock := espsocket.NewConnSocket(conn, this.maxPackage)
	defer sock.AskFinish()
	sn := serviceName
	if sn == "" {
		sn = this.name
	}
	return ESNPCall(sock, sn, method, params, tmd)
}

type esnpConfig struct {
	Net        string
	Address    string
	MaxPackage int
	TimeoutMS  int
	InitSize   int
	MaxSize    int
	IdleTimeMS int
}

type ESNPServiceCallerFactory int

func (o ESNPServiceCallerFactory) Valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) error {
	var co esnpConfig
	if valutil.ToBean(cfg, &co) {
		if co.Net == "" {
			co.Net = "tcp"
		}
		if co.Address == "" {
			return fmt.Errorf("Address empty")
		}
		return nil
	}
	return fmt.Errorf("invalid ESNPServiceCaller config")
}

func (o ESNPServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}, ofp objfac.ObjectFactoryProvider) (same bool) {
	var co, oo esnpConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.Net != oo.Net {
		return false
	}
	if co.Address != oo.Address {
		return false
	}
	if co.TimeoutMS != oo.TimeoutMS {
		return false
	}
	if co.InitSize != oo.InitSize {
		return false
	}
	if co.MaxSize != oo.MaxSize {
		return false
	}
	if co.IdleTimeMS != oo.IdleTimeMS {
		return false
	}
	if co.MaxPackage != oo.MaxPackage {
		return false
	}
	return true
}

func (o ESNPServiceCallerFactory) Create(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (interface{}, error) {
	err := o.Valid(cfg, ofp)
	if err != nil {
		return nil, err
	}
	var co esnpConfig
	valutil.ToBean(cfg, &co)

	pcfg := new(conndialpool.DialPoolConfig)
	pcfg.Net = co.Net
	pcfg.Address = co.Address
	pcfg.InitSize = co.InitSize
	pcfg.MaxSize = co.MaxSize
	if pcfg.MaxSize <= 0 {
		pcfg.MaxSize = 128
	}
	pcfg.IdleMS = co.IdleTimeMS
	pcfg.Valid()
	pool := conndialpool.NewDialPool("serviceCall", pcfg)

	r := new(ESNPServiceCaller)
	r.pool = pool
	r.timeoutMS = co.TimeoutMS
	r.maxPackage = co.MaxPackage
	return ServiceCaller(r), nil
}
