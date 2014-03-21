package n2n

import (
	"bmautil/goo"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/espterminal"
	"logger"
	"time"
)

const (
	tag = "n2n"
)

type RemoteListener interface {
	OnRemoteConnect(nodeId nodeinfo.NodeId)

	OnRemoteDisconnect(nodeId nodeinfo.NodeId)
}

type Service struct {
	name   string
	ninfo  *nodeinfo.Service
	config *configInfo

	goo        goo.Goo
	connectors map[string]*connector
	remotes    map[nodeinfo.NodeId]*remoteInfo
	listeners  map[string]RemoteListener
}

func NewService(n string, ninfo *nodeinfo.Service) *Service {
	this := new(Service)
	this.name = n
	this.ninfo = ninfo
	this.connectors = make(map[string]*connector)
	this.remotes = make(map[nodeinfo.NodeId]*remoteInfo)
	this.goo.InitGoo(tag, 128, this.doExit)
	return this
}

func (this *Service) doExit() {
	this.doCloseAllConnectors()
	this.doCloseAllRemote()
}

func (this *Service) doCheckConnector(k string, url *esnp.URL) {
	ctor, ok := this.connectors[k]
	if ok {
		return
	}
	ctor = new(connector)
	this.connectors[k] = ctor
	ctor.InitConnector(this, k, url)
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
	old, ok := this.remotes[ri.nodeId]
	if !ok {
		return
	}
	if old != ri {
		return
	}
	if !ri.tunnel.IsBreak() {
		return
	}
	logger.Debug(tag, "'%s' RemoteClosed", ri)
	delete(this.remotes, ri.nodeId)
	if this.listeners != nil {
		for _, lis := range this.listeners {
			go lis.OnRemoteDisconnect(ri.nodeId)
		}
	}
}

func (this *Service) doCloseAllRemote() {
	tmp := this.remotes
	this.remotes = make(map[nodeinfo.NodeId]*remoteInfo)

	for _, ri := range tmp {
		ri.tunnel.SetCloseListener("this", nil)
		ri.Close()
		if this.listeners != nil {
			for _, lis := range this.listeners {
				go lis.OnRemoteDisconnect(ri.nodeId)
			}
		}
	}
}

func (this *Service) doChannelAccept(n string, url *esnp.URL, ch espchannel.Channel) {
	go func() {
		logger.Debug(tag, "send joinReq -> (%s : %s)", n, ch)

		req := this.makeJoinReq()

		msg := esnp.NewRequestMessage()
		addr := msg.GetAddress()
		url.BindAddress(addr, false)
		addr.SetOp(OP_JOIN)
		req.Write(msg)

		tm := new(espterminal.Terminal)
		tm.InitTerminal(n)
		tm.SetMessageListner(func(msg *esnp.Message) error {
			return this.Serve(ch, msg)
		})
		tm.Connect(ch)

		to := time.NewTimer(url.GetTimeout(3 * time.Second))
		rmsg, err := tm.Call(ch, msg, to)
		if err != nil {
			logger.Debug(tag, "%s call fail - %s", ch, err)
			return
		}
		tm.Disconnect(ch)
		err = this.handleJoin(ch, rmsg, false)
		if err != nil {
			logger.Debug(tag, "%s handle join resp fail - %s", ch, err)
			ch.AskClose()
		}
		return
	}()
}

func (this *Service) doJoin(req *joinReq, ch espchannel.Channel) error {
	logger.Debug(tag, "doJoin(%v)", req)
	ri, ok := this.remotes[req.Id]
	if !ok {
		logger.Debug(tag, "create remoteInfo(%v)", req)

		ri = new(remoteInfo)
		err := ri.InitRemoteInfo(this, req.Id, req.Name, req.URL)
		if err != nil {
			return err
		}
		this.remotes[req.Id] = ri
	}
	ri.Add(ch)
	if !ok {
		if this.listeners != nil {
			for _, lis := range this.listeners {
				go lis.OnRemoteConnect(ri.nodeId)
			}
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

func (this *Service) GetChannel(nid nodeinfo.NodeId) (espchannel.Channel, *espterminal.Terminal, error) {
	var rch espchannel.Channel
	var rtm *espterminal.Terminal
	err := this.goo.DoSync(func() {
		ri, ok := this.remotes[nid]
		if !ok {
			return
		}
		rch = ri.tunnel
		rtm = ri.tm
	})
	if err != nil {
		return nil, nil, err
	}
	return rch, rtm, nil
}
