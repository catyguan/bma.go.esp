package servicecall

import (
	"bmautil/conndialpool"
	"bmautil/valutil"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"time"
)

type ESNPServiceCaller struct {
	name       string
	pool       *conndialpool.DialPool
	maxPackage int
	timeoutMS  int
	runtime    bool
}

func NewESNPServiceCaller(n string, cfg map[string]interface{}, rt bool) (*ESNPServiceCaller, error) {
	fac := ESNPServiceCallerFactory(0)
	err := fac.Valid(cfg)
	if err != nil {
		return nil, err
	}
	sc, err1 := fac.Create(n, cfg)
	if err1 != nil {
		return nil, err1
	}
	r := sc.(*ESNPServiceCaller)
	r.runtime = rt
	return r, nil
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

func (this *ESNPServiceCaller) IsRuntime() bool {
	return this.runtime
}

func (this *ESNPServiceCaller) Call(method string, params map[string]interface{}, timeout time.Duration) (interface{}, error) {
	tm := this.timeoutMS
	if tm <= 0 {
		tm = 5000
	}
	tmd := time.Duration(tm) * time.Millisecond
	if timeout != time.Duration(0) && timeout < tmd {
		tmd = timeout
	}
	conn, err0 := this.pool.GetConn(tmd, true)
	if err0 != nil {
		return nil, err0
	}
	sock := espsocket.NewConnSocket(conn, this.maxPackage)
	defer sock.AskFinish()

	msg := esnp.NewRequestMessageWithId()
	addr := msg.GetAddress()
	addr.SetCall(this.name, method)
	dt := msg.Datas()
	dt.Set("p", params)

	ts := time.Now()
	rmsg, err1 := espsocket.CallTimeout(sock, msg, tmd)
	te := time.Now()
	if err1 != nil {
		if rmsg == nil {
			sock.AskClose()
		}
		logger.Debug(tag, "[%s:%s] esnp call(%f) fail '%s'", this.name, method, te.Sub(ts).Seconds(), err1)
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

	logger.Debug(tag, "[%s:%s] esnp call(%f) end '%d'", this.name, method, te.Sub(ts).Seconds(), st)
	if st != 200 {
		msg, _ := dt.GetString("m", "")
		if msg == "" {
			msg = fmt.Sprintf("invalid esnp status(%d)", st)
		}
		return nil, fmt.Errorf(msg)
	}
	return val, nil
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

func (o ESNPServiceCallerFactory) Valid(cfg map[string]interface{}) error {
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

func (o ESNPServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool) {
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

func (o ESNPServiceCallerFactory) Create(n string, cfg map[string]interface{}) (ServiceCaller, error) {
	err := o.Valid(cfg)
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
	pool := conndialpool.NewDialPool("serviceCall_"+n, pcfg)

	r := new(ESNPServiceCaller)
	r.name = n
	r.pool = pool
	r.timeoutMS = co.TimeoutMS
	r.maxPackage = co.MaxPackage
	return r, nil
}
