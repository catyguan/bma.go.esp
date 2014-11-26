package servicecall

import (
	"esp/espnet/esnp"
	"fmt"
	"logger"
	"time"
)

type ESNPServiceCaller struct {
	name      string
	provider  SocketProvider
	timeoutMS int
}

func (this *ESNPServiceCaller) Ping() bool {
	sock, err := this.provider.GetSocket()
	if err != nil {
		return false
	}
	if sock == nil {
		return false
	}
	return !sock.IsBreak()
}

func (this *ESNPServiceCaller) Start() error {
	return nil
}

func (this *ESNPServiceCaller) Stop() {

}

func (this *ESNPServiceCaller) Call(method string, params []interface{}, timeout time.Duration) (interface{}, error) {
	sock, err0 := this.provider.GetSocket()
	if err0 != nil {
		return nil, err0
	}
	msg := esnp.NewRequestMessage()
	addr := msg.GetAddress()
	addr.SetCall(this.name, method)
	dt := msg.Datas()
	dt.Set("p", params)

	tm := this.timeoutMS
	if tm <= 0 {
		tm = 5000
	}
	tmd := time.Duration(tm) * time.Millisecond
	if timeout != time.Duration(0) && timeout < tmd {
		tmd = timeout
	}
	ts := time.Now()
	rmsg, err1 := sock.Call(msg, tmd)
	te := time.Now()
	if err1 != nil {
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
