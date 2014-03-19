package n2n

import (
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"logger"
	"sync"
)

const (
	tag = "n2n"
)

type Service struct {
	name   string
	ninfo  *nodeinfo.Service
	config *configInfo

	lock       sync.RWMutex
	connectors map[string]*connector
	remotes    map[nodeinfo.NodeId]*remoteInfo
}

func NewService(n string, ninfo *nodeinfo.Service) *Service {
	this := new(Service)
	this.name = n
	this.ninfo = ninfo
	this.connectors = make(map[string]*connector)
	this.remotes = make(map[nodeinfo.NodeId]*remoteInfo)
	return this
}

func (this *Service) onChannelClose(ri *remoteInfo) {
	logger.Debug(tag, "node[%d] break", ri.nodeId)
}

func (this *Service) checkConnector(k string, url *esnp.URL) {
	this.lock.Lock()
	ctor, ok := this.connectors[k]
	if ok {
		this.lock.Unlock()
		return
	}
	ctor = new(connector)
	this.connectors[k] = ctor
	this.lock.Unlock()
	ctor.InitConnector(this, k, url)
}

func (this *Service) closeConnector(k string) {
	this.lock.Lock()
	ctor, ok := this.connectors[k]
	delete(this.connectors, k)
	this.lock.Unlock()
	if ok {
		ctor.Close()
	}
}

func (this *Service) closeAllConnectors() {
	this.lock.Lock()
	tmp := this.connectors
	this.connectors = make(map[string]*connector)
	this.lock.Unlock()

	for _, ctor := range tmp {
		ctor.Close()
	}
}

func (this *Service) closeAllRemote() {
	this.lock.Lock()
	tmp := this.remotes
	this.remotes = make(map[nodeinfo.NodeId]*remoteInfo)
	this.lock.Unlock()

	for _, ri := range tmp {
		ri.Close()
	}
}

func (this *Service) doJoin(req *joinReq, ch espchannel.Channel) error {
	logger.Debug(tag, "doJoin(%v)", req)
	this.lock.Lock()
	defer this.lock.Unlock()
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
	return nil
}
