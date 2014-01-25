package nodegroup

import (
	"esp/cluster/nodeid"
	"esp/espnet"
	"sync"
)

type Action int

const (
	ACTION_JOIN   = Action(1)
	ACTION_LEAVE  = Action(2)
	ACTION_UPDATE = Action(3)
)

type Listener func(action Action, nid nodeid.NodeId, ch espnet.Channel, data interface{})

type ngItem struct {
	channel espnet.Channel
	data    interface{}
}

type NodeGroup struct {
	name      string
	lock      sync.RWMutex
	items     map[nodeid.NodeId]*ngItem
	listeners map[string]Listener
}

func NewNodeGroup(name string) *NodeGroup {
	this := new(NodeGroup)
	this.name = name
	this.items = make(map[nodeid.NodeId]*ngItem)
	this.listeners = make(map[string]Listener)
	return this
}

func (this *NodeGroup) Name() string {
	return this.name
}

func (this *NodeGroup) doAddListener(n string, lis Listener) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.listeners[n] = lis
}

func (this *NodeGroup) AddListener(n string, lis Listener, fireEvent bool) {
	this.doAddListener(n, lis)
	if fireEvent {
		this.lock.RLock()
		defer this.lock.RUnlock()
		for nid, item := range this.items {
			lis(ACTION_JOIN, nid, item.channel, item.data)
		}
	}
}

func (this *NodeGroup) RemoveListener(n string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.listeners, n)
}

func (this *NodeGroup) doJoin(nid nodeid.NodeId, ch espnet.Channel, data interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.items[nid]
	if ok {
		return false
	}
	item := new(ngItem)
	item.channel = ch
	item.data = data
	this.items[nid] = item
	return true
}

func (this *NodeGroup) Join(nid nodeid.NodeId, ch espnet.Channel, data interface{}) bool {
	if !this.doJoin(nid, ch, data) {
		return false
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, lis := range this.listeners {
		lis(ACTION_JOIN, nid, ch, data)
	}
	return true
}

func (this *NodeGroup) doLeave(nid nodeid.NodeId) *ngItem {
	this.lock.Lock()
	defer this.lock.Unlock()
	old, ok := this.items[nid]
	if !ok {
		return nil
	}
	delete(this.items, nid)
	return old
}

func (this *NodeGroup) Leave(nid nodeid.NodeId) bool {
	old := this.doLeave(nid)
	if old == nil {
		return false
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, lis := range this.listeners {
		lis(ACTION_LEAVE, nid, old.channel, old.data)
	}
	return true
}

func (this *NodeGroup) Set(nid nodeid.NodeId, data interface{}) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	item, ok := this.items[nid]
	if ok {
		item.data = data
		for _, lis := range this.listeners {
			lis(ACTION_UPDATE, nid, item.channel, data)
		}
		return true
	}
	return false
}

func (this *NodeGroup) LockAndSet(nid nodeid.NodeId, data interface{}) bool {
	r := false
	this.lock.Lock()
	item, ok := this.items[nid]
	if ok {
		item.data = data
		r = true
	}
	this.lock.Unlock()

	if !r {
		return false
	}

	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, lis := range this.listeners {
		lis(ACTION_UPDATE, nid, item.channel, data)
	}
	return r
}

func (this *NodeGroup) Get(nid nodeid.NodeId) (espnet.Channel, interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	item, ok := this.items[nid]
	if ok {
		return item.channel, item.data
	}
	return nil, nil
}
