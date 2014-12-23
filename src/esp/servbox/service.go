package servbox

import (
	"bmautil/conndialpool"
	"bmautil/connutil"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"sync"
	"time"
)

type nodeInfo struct {
	name     string
	net      string
	address  string
	services []string
	pool     *conndialpool.DialPool
}

func (this *nodeInfo) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
}

type servItem struct {
	node *nodeInfo
}

type Service struct {
	name   string
	config *configInfo
	lock   sync.RWMutex
	nodes  map[string]*nodeInfo
	servs  map[string]*servItem
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	r.nodes = make(map[string]*nodeInfo)
	r.servs = make(map[string]*servItem)
	return r
}

func (this *Service) Handler(sin espsocket.Socket, msg *esnp.Message) error {
	sname, err1 := esnp.MessageLineCoders.Address.Get(msg, esnp.ADDRESS_SERVICE)
	if err1 != nil {
		return err1
	}
	var node *nodeInfo
	this.lock.RLock()
	if si, ok := this.servs[sname]; ok {
		if si.node != nil {
			node = si.node
		}
	}
	this.lock.RUnlock()
	if node == nil {
		return fmt.Errorf("unknow service[%s]", sname)
	}
	var conn *connutil.ConnExt
	for {
		conn, err1 = node.pool.GetConn(time.Duration(this.config.TimeoutMS)*time.Millisecond, true)
		if err1 != nil {
			return err1
		}
		if conn.CheckBreak() {
			node.pool.CloseConn(conn)
			continue
		}
		break
	}
	logger.Debug(tag, "service[%s] --> %s", sname, conn)
	sout := espsocket.NewConnSocket(conn, this.config.MaxPackage)
	defer sout.AskFinish()
	err2 := sout.WriteMessage(msg)
	if err2 != nil {
		sout.AskClose()
		return err2
	}
	if !msg.IsRequest() {
		logger.Debug(tag, "not request, skip response")
		return nil
	}
	rmsg, err3 := sout.ReadMessage()
	if rmsg != nil {
		err4 := sin.WriteMessage(rmsg)
		if err4 != nil {
			sout.AskClose()
			return err4
		}
	} else {
		sout.AskClose()
		return err3
	}
	logger.Debug(tag, "service[%s] --> %s done", sname, conn)
	return nil
}

func (this *Service) ManageHandler(sock espsocket.Socket, msg *esnp.Message) error {
	op, err0 := esnp.MessageLineCoders.Address.Get(msg, esnp.ADDRESS_OP)
	if err0 != nil {
		return err0
	}
	switch op {
	case op_Join:
		return this.doJoin(sock, msg)
	}
	return espservice.Miss(msg)
}
