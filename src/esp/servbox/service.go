package servbox

import (
	"bmautil/conndialpool"
	"bmautil/connutil"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"esp/espnet/proxy"
	"fmt"
	"logger"
	"net"
	"sync"
	"time"
)

type nodeInfo struct {
	name     string
	net      string
	address  string
	services []string
	info     string
	skipKill bool
	pool     *conndialpool.DialPool
}

func (this *nodeInfo) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
}

func (this *nodeInfo) String() string {
	return fmt.Sprintf("%s,%s", this.name, this.info)
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
	logger.Debug(tag, "service[%s] --> %s, %s", sname, node, conn)
	sout := espsocket.NewConnSocket(conn, this.config.MaxPackage)
	defer sout.AskFinish()
	var fs proxy.ForwardSetting
	err2, _ := proxy.Forward(sin, sout, msg, &fs)
	if err2 != nil {
		// no failOver
		return err2
	}
	logger.Debug(tag, "service[%s] --> %s, %s done", sname, node, conn)
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

func (this *Service) AcceptManageConn(debug string) func(conn net.Conn) {
	goservice := espservice.NewGoService(this.name+"_manageService", this.ManageHandler)
	return goservice.AcceptConn(debug)
}
