package emnode

import (
	"boot"
	"esp/cluster/nodeid"
	"esp/sqlite"
	"fmt"
	"logger"
	"sync"
)

type Service struct {
	name       string
	database   *sqlite.SqliteServer
	nidService *nodeid.Service
	nodeId     nodeid.NodeId

	lock  sync.RWMutex
	nodes map[string]*ExecuteManageNode
}

func NewService(name string, db *sqlite.SqliteServer, nid *nodeid.Service) *Service {
	this := new(Service)
	this.name = name
	this.database = db
	this.nidService = nid
	this.groups = make(map[string]*NodeGroup)
	return this
}

func (this *Service) Name() string {
	return this.name
}

type configInfo struct {
	NodeId uint64
}

func (this *Service) onNodeIdChange(nodeId nodeid.NodeId) {
	// TODO
}

func (this *Service) Start() bool {
	this.nodeId = this.nidService.GetAndListen("nodegroup_servie_"+this.name, this.onNodeIdChange)
	return true
}

func (this *Service) Close() bool {
	nl := make([]string, 0)
	this.lock.RLock()
	for n, _ := range this.groups {
		nl = append(nl, n)
	}
	this.lock.RUnlock()

	for _, n := range nl {
		this.CloseGroup(n, false)
	}
	return true
}

func (this *Service) GetNodeId() nodeid.NodeId {
	return this.nodeId
}

func (this *Service) GetGroup(n string) *NodeGroup {
	this.lock.RLock()
	defer this.lock.RUnlock()
	o, ok := this.groups[n]
	if ok {
		return o
	}
	return nil
}

func (this *Service) CreateGroup(n string, cfg *NodeGroupConfig) (*NodeGroup, error) {
	ng, err := func() (*NodeGroup, error) {
		this.lock.Lock()
		defer this.lock.Unlock()
		_, ok := this.groups[n]
		if ok {
			return nil, fmt.Errorf("group(%s) exists", n)
		}
		r := newNodeGroup(n, this, cfg)
		this.groups[n] = r
		return r, nil
	}()
	if err != nil {
		return nil, err
	}
	if !boot.RuntimeStartRun(ng) {
		boot.RuntimeStopCloseClean(ng, false)
		this.lock.Lock()
		delete(this.groups, n)
		this.lock.Unlock()
		return nil, fmt.Errorf("group(%s) start fail", n)
	}
	logger.Debug(tag, "CreateGroup(%s) done", n)
	return ng, nil
}

func (this *Service) CloseGroup(n string, wait bool) *NodeGroup {
	this.lock.Lock()
	ng, ok := this.groups[n]
	if ok {
		delete(this.groups, n)
	}
	this.lock.Unlock()
	if ok {
		boot.RuntimeStopCloseClean(ng, false)
		logger.Debug(tag, "CloseGroup(%s) done", n)
		return ng
	}
	return nil
}
