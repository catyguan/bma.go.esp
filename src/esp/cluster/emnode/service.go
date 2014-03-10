package emnode

import (
	"boot"
	"esp/cluster/nodeinfo"
	"fmt"
	"logger"
	"sync"
)

type Service struct {
	name            string
	nodeInfoService *nodeinfo.Service

	lock  sync.RWMutex
	nodes map[string]*ExecuteManageNode
}

func NewService(name string, nis *nodeinfo.Service) *Service {
	this := new(Service)
	this.name = name
	this.nodeInfoService = nis
	this.nodes = make(map[string]*ExecuteManageNode)
	return this
}

func (this *Service) Name() string {
	return this.name
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
