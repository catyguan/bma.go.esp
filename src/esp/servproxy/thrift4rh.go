package servproxy

import (
	"bmautil/conndialpool"
	"bmautil/valutil"
	"fmt"
	"time"
)

type thriftRemoteParam struct {
	PoolMax  int
	PoolInit int
}

type thriftRemoteData struct {
	params *thriftRemoteParam
	pool   *conndialpool.DialPool
}

type ThriftRemoteHandler int

func init() {
	AddRemoteHandler("thrift", ThriftRemoteHandler(0))
}

func (this ThriftRemoteHandler) Ping(remote *RemoteObj) (bool, bool) {
	if remote.Data == nil {
		return true, false
	}
	tdata, ok := remote.Data.(*thriftRemoteData)
	if !ok {
		return true, false
	}
	if tdata.pool == nil {
		return true, false
	}
	if tdata.pool.GetInitSize() > 0 {
		return true, tdata.pool.ActiveConn() > 0
	}
	return false, false
}

func (this ThriftRemoteHandler) Valid(cfg *RemoteConfigInfo) error {
	if cfg.Host == "" {
		return fmt.Errorf("Host invalid")
	}
	return nil
}

func (this ThriftRemoteHandler) Compare(cfg *RemoteConfigInfo, old *RemoteConfigInfo) bool {
	p1 := new(thriftRemoteParam)
	valutil.ToBean(cfg.Params, p1)
	p2 := new(thriftRemoteParam)
	valutil.ToBean(cfg.Params, p2)

	if p1.PoolMax != p2.PoolMax {
		return false
	}
	if p2.PoolInit != p2.PoolInit {
		return false
	}
	return true
}

func (this ThriftRemoteHandler) Start(o *RemoteObj) error {
	rcfg := o.cfg
	p := new(thriftRemoteParam)
	valutil.ToBean(rcfg.Params, p)

	data := new(thriftRemoteData)
	o.Data = data

	data.params = p

	cfg := new(conndialpool.DialPoolConfig)
	cfg.Address = rcfg.Host
	tm := rcfg.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	cfg.TimeoutMS = tm
	if p.PoolInit < 0 {
		cfg.InitSize = 1
	} else {
		cfg.InitSize = p.PoolInit
	}
	if p.PoolMax <= 0 {
		cfg.MaxSize = 10
	} else {
		cfg.MaxSize = p.PoolMax
	}
	err := cfg.Valid()
	if err != nil {
		return err
	}
	pool := conndialpool.NewDialPool(fmt.Sprintf("%s_remote", o.name), cfg)
	data.pool = pool
	if !pool.StartAndRun() {
		return fmt.Errorf("start remote pool fail")
	}
	return nil
}

func (this ThriftRemoteHandler) Stop(o *RemoteObj) error {
	if o.Data == nil {
		return nil
	}
	data, ok := o.Data.(*thriftRemoteData)
	if !ok {
		return nil
	}
	if data.pool != nil {
		data.pool.AskClose()
	}
	return nil
}

func (this ThriftRemoteHandler) Begin(o *RemoteObj, timeout time.Duration) (RemoteSession, error) {
	if o.Data == nil {
		return nil, fmt.Errorf("invalid ThriftRemoteData")
	}
	data, ok := o.Data.(*thriftRemoteData)
	if !ok {
		return nil, fmt.Errorf("invalid ThriftRemoteData(%T)", o.Data)
	}
	pool := data.pool
	if pool == nil {
		return nil, fmt.Errorf("invalid ThriftRemoteData ConnPool")
	}
	if true {
		tm := o.cfg.TimeoutMS
		if tm <= 0 {
			tm = 5 * 1000
		}
		xtimeout := time.Duration(tm) * time.Millisecond
		if timeout == 0 {
			timeout = xtimeout
		} else if timeout > xtimeout {
			timeout = xtimeout
		}
	}
	for {
		conn, err := pool.GetConn(timeout, true)
		if err != nil {
			return nil, err
		}
		if conn.CheckBreak() {
			pool.CloseConn(conn)
		} else {
			if debugTraffic {
				conn.Debuger = ConnDebuger
			}
			sess := new(ThriftRemoteSession)
			sess.remote = o
			sess.pool = pool
			sess.conn = conn
			return sess, nil
		}
	}
}
