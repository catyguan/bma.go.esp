package n2n

import (
	"bmautil/goo"
	"container/list"
	"esp/cluster/nodebase"
	"esp/espnet/auth"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"net"
	"time"
)

const (
	tag = "n2n"
)

type RemoteListener interface {
	OnRemoteConnect(nodeId nodebase.NodeId, sock *espsocket.Socket)

	OnRemoteDisconnect(nodeId nodebase.NodeId)
}

type Service struct {
	config *ConfigInfo

	goo        goo.Goo
	connectors map[string]*connector
	remotes    map[nodebase.NodeId]*list.List
	listeners  map[string]RemoteListener
}

func NewService(qsize int) *Service {
	this := new(Service)
	this.connectors = make(map[string]*connector)
	this.remotes = make(map[nodebase.NodeId]*list.List)
	this.goo.InitGoo(tag, qsize, this.doExit)
	return this
}

func (this *Service) doExit() {
	this.doCloseAllConnectors()
	this.doCloseAllRemote()
}

func (this *Service) doCheckConnector(k string, host string, code string) {
	ctor, ok := this.connectors[k]
	if ok {
		return
	}
	ctor = new(connector)
	this.connectors[k] = ctor
	ctor.InitConnector(this, k, host, code)
}

func (this *Service) doCloseConnector(k string) {
	ctor, ok := this.connectors[k]
	delete(this.connectors, k)
	if ok {
		ctor.Close()
	}
}

func (this *Service) doCloseAllConnectors() {
	tmp := this.connectors
	this.connectors = make(map[string]*connector)

	for _, ctor := range tmp {
		ctor.Close()
	}
}

func (this *Service) doRemoteClosed(ri *remoteInfo) {
	l, ok := this.remotes[ri.nodeId]
	if !ok {
		return
	}
	find := false
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == ri {
			l.Remove(e)
			find = true
			break
		}
	}
	if !find {
		return
	}
	logger.Debug(tag, "'%s' RemoteClosed", ri)
	if l.Len() == 0 {
		delete(this.remotes, ri.nodeId)
		if this.listeners != nil {
			for _, lis := range this.listeners {
				go lis.OnRemoteDisconnect(ri.nodeId)
			}
		}
	}
}

func (this *Service) doCloseAllRemote() {
	for nid, l := range this.remotes {
		for e := l.Front(); e != nil; e = e.Next() {
			ri := e.Value.(*remoteInfo)
			ri.close(true)
		}
		delete(this.remotes, nid)
		if this.listeners != nil {
			for _, lis := range this.listeners {
				go lis.OnRemoteDisconnect(nid)
			}
		}
	}
}

func (this *Service) doSocketAccept(n string, host string, code string, sock *espsocket.Socket) {
	go func() {
		// do login
		errA := this.PostAuth(sock, code)
		if errA != nil {
			logger.Warn(tag, "%s auth fail - %s", sock, errA)
			time.Sleep(10 * time.Second)
			sock.AskClose()
			return
		}

		// do join
		logger.Debug(tag, "send joinReq -> (%s : %s)", n, sock)

		req := this.makeJoinReq()

		msg := esnp.NewRequestMessage()
		addr := msg.GetAddress()
		addr.SetHost(host)
		addr.SetService(SN_N2N)
		addr.SetOp(OP_JOIN)
		req.Write(msg)

		rmsg, err := sock.Call(msg, time.Duration(this.config.TimeoutMS)*time.Millisecond)
		if err != nil {
			logger.Warn(tag, "%s call fail - %s", sock, err)
			return
		}
		is, _ := auth.CheckAuth(rmsg)
		if is {
			logger.Warn(tag, "need login")
			return
		}
		err = this.handleJoin(sock, rmsg, false)
		if err != nil {
			logger.Warn(tag, "%s handle join resp fail - %s", sock, err)
			sock.AskClose()
		}
		return
	}()
}

func ValidHost(host string, shost interface{}) string {
	h1, p1, _ := net.SplitHostPort(host)
	if h1 == "_" || h1 == "" {
		if shost != nil {
			str := fmt.Sprintf("%s", shost)
			h2, _, _ := net.SplitHostPort(str)
			return net.JoinHostPort(h2, p1)
		}
	}
	return host
}

func (this *Service) doJoin(req *joinReq, sock *espsocket.Socket) error {
	ri := new(remoteInfo)
	shost, _ := sock.GetProperty(espsocket.PROP_SOCKET_REMOTE_ADDR)
	host := ValidHost(req.Host, shost)
	logger.Debug(tag, "doJoin(%v, %s)", req, host)
	err := ri.InitRemoteInfo(this, req.Id, req.Name, host, sock)
	if err != nil {
		sock.AskClose()
		return err
	}
	l, ok := this.remotes[req.Id]
	if !ok {
		logger.Debug(tag, "create remoteInfo(%v, %s)", req, host)
		l = list.New()
		this.remotes[req.Id] = l
	} else {
		e := l.Front()
		if e != nil {
			old := e.Value.(*remoteInfo)
			if old.nodeName != req.Name {
				return logger.Warn(tag, "node(%d) name not same(%s , %s), refuse join", req.Id, old.nodeName, req.Name)
			}
		}
	}
	l.PushBack(ri)
	if this.listeners != nil && !ok {
		for _, lis := range this.listeners {
			go lis.OnRemoteConnect(req.Id, sock)
		}
	}
	return nil
}

func (this *Service) SetListener(n string, lis RemoteListener) error {
	return this.goo.DoSync(func() {
		if this.listeners == nil {
			this.listeners = make(map[string]RemoteListener)
		}
		if lis == nil {
			delete(this.listeners, n)
		} else {
			this.listeners[n] = lis
		}
	})
}

func (this *Service) GetSocket(nid nodebase.NodeId) (*espsocket.Socket, error) {
	var rsock *espsocket.Socket
	err := this.goo.DoSync(func() {
		l, ok := this.remotes[nid]
		if !ok {
			return
		}
		e := l.Front()
		if e != nil {
			rsock = e.Value.(*remoteInfo).sock
		}
	})
	if err != nil {
		return nil, err
	}
	return rsock, nil
}
