package acclog

import (
	"logger"
	"sync"
)

const (
	tag = "acclog"
)

type alNode struct {
	cfg  *fileConfig
	alog *AccessLoggerFile
}

type Service struct {
	name   string
	config *configInfo
	nodes  map[string]*alNode
	lock   sync.RWMutex
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.nodes = make(map[string]*alNode)
	return this
}

func (this *Service) _createNode(fc *fileConfig) {
	n := fc.FilePrex()
	fnc := NewFilenameCreator(n)
	alog := NewFile(fnc, fc.QueueSize, true, nil)
	if fc.MaxLines > 0 {
		alog.SetRotateLines(fc.MaxLines)
	}
	if fc.MaxSize > 0 {
		alog.SetRotateSize(fc.MaxSize)
	}
	alog.SetRotateDaily(!fc.NoDaily)

	node := new(alNode)
	node.alog = alog
	node.cfg = fc

	this.nodes[fc.Name] = node
}

func (this *Service) _closeNode(n string, node *alNode) {
	delete(this.nodes, n)
	if node.alog != nil {
		node.alog.Close()
	}
}

func (this *Service) Write(n string, info AccLogInfo) {
	var node *alNode
	ok := false

	this.lock.RLock()
	if node, ok = this.nodes[n]; !ok {
		node = this.nodes["*"]
	}
	this.lock.RUnlock()
	if node == nil {
		this.lock.Lock()
		node = new(alNode)
		this.nodes[n] = node
		this.lock.Unlock()
		logger.Warn(tag, "miss '%s' access log", n)
		return
	}
	if node.alog != nil {
		node.alog.Write(info)
	}
}
