package router

import (
	"bmautil/qexec"
	"boot"
	"config"
	"errors"
	"esp/espnet"
	"fmt"
	"logger"
	"time"
)

const (
	tag = "Router"
)

type RouterDispatcher interface {
	Dispatch(n string) (espnet.ChannelFactory, error)
	Close()
}

type Router struct {
	name string

	executor    qexec.QueueExecutor
	fastCache   map[string]espnet.ChannelFactory
	factorys    map[string]espnet.ChannelFactory
	dispatchers []RouterDispatcher
	tasks       map[string]*dispatchTask
}

func NewRouter(name string) *Router {
	this := new(Router)
	this.name = name
	this.executor.InitQueueExecutor(tag, 16, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.fastCache = make(map[string]espnet.ChannelFactory)
	this.factorys = make(map[string]espnet.ChannelFactory)
	this.dispatchers = make([]RouterDispatcher, 0)
	this.tasks = make(map[string]*dispatchTask)
	return this
}

func (this *Router) requestHandler(o interface{}) (bool, error) {
	switch v := o.(type) {
	case func() error:
		return true, v()
	}
	return true, nil
}

func (this *Router) stopHandler() {
	for k, task := range this.tasks {
		delete(this.tasks, k)
		for _, resp := range task.respList {
			this.safe(resp, nil)
		}
	}
	for _, dis := range this.dispatchers {
		dis.Close()
	}
	this.fastCache = make(map[string]espnet.ChannelFactory)
	for k, _ := range this.factorys {
		delete(this.factorys, k)
	}
}

func (this *Router) Name() string {
	return this.name
}

func (this *Router) copyOnWrite() {
	m := make(map[string]espnet.ChannelFactory, len(this.fastCache))
	for k, f := range this.factorys {
		m[k] = f
	}
	this.fastCache = m
}

func (this *Router) InitFactory(n string, cf espnet.ChannelFactory) {
	this.factorys[n] = cf
	this.fastCache[n] = cf
}
func (this *Router) InitDispatcher(dis RouterDispatcher) {
	this.dispatchers = append(this.dispatchers, dis)
}

type configInfo struct {
	QueueSize int
	Factory   map[string]map[string]interface{}
}

func (this *Router) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		sz := cfg.QueueSize
		if sz > 0 {
			this.executor.InitRequests(sz)
		}
		if cfg.Factory != nil {
			for key, info := range cfg.Factory {
				if key == "" {
					logger.Error(tag, "factory name empty")
					return false
				}
				if info == nil {
					logger.Error(tag, "factory config empty")
					return false
				}
				typ := ""
				if true {
					v, ok := info["Type"]
					if ok {
						typ = fmt.Sprintf("%s", v)
					}
				}
				if typ == "" {
					logger.Error(tag, "factory[%s] type empty", key)
					return false
				}
				var nlist []string
				if true {
					v, ok := info["Address"]
					if ok {
						if l, ok2 := v.([]interface{}); ok2 {
							nlist = make([]string, len(l))
							for _, v2 := range l {
								nlist = append(nlist, fmt.Sprintf("%s", v2))
							}
						}
					}
				}
				if nlist == nil || len(nlist) == 0 {
					logger.Error(tag, "factory[%s] address empty", key)
					return false
				}

				fac, err := NewChannelFactory(typ, key, info)
				if err != nil {
					logger.Error(tag, "factory[%s] create fail - %s", key, err)
					return false
				}
				logger.Debug(tag, "factory[%s] - %s install", key, typ)
				boot.QuickDefine(fac, "fac_"+key, false)
				for _, al := range nlist {
					this.InitFactory(al, fac)
				}
			}
		}
		return true
	}
	logger.Error(tag, "GetBeanConfig(%s) fail", this.name)
	return false
}

func (this *Router) Start() bool {
	return this.executor.Run()
}

func (this *Router) Stop() bool {
	this.executor.Stop()
	return true
}

func (this *Router) Cleanup() bool {
	this.executor.WaitStop()
	return true
}

type dispatchTask struct {
	count    int
	respList []chan espnet.ChannelFactory
}

func (this *Router) safe(resp chan espnet.ChannelFactory, cf espnet.ChannelFactory) {
	defer func() {
		recover()
	}()
	resp <- cf
}

func (this *Router) endDispatch(n string, cf espnet.ChannelFactory) {
	task, ok := this.tasks[n]
	if !ok {
		return
	}
	if cf == nil {
		task.count--
		if task.count <= 0 {
			delete(this.tasks, n)
			for _, resp := range task.respList {
				this.safe(resp, nil)
			}
			task.respList = nil
			logger.Debug(tag, "'%s' all dispatch fail", n)
		}
		return
	}

	if cf != nil {
		delete(this.tasks, n)
		this.factorys[n] = cf
		this.copyOnWrite()
		for _, resp := range task.respList {
			this.safe(resp, cf)
		}
		task.respList = nil
	}
}

func (this *Router) callDispatch(dis RouterDispatcher, n string) {
	defer func() {
		this.executor.DoNow("endDispatch", func() error {
			this.endDispatch(n, nil)
			return nil
		})
	}()
	f, err := dis.Dispatch(n)
	if err != nil {
		logger.Debug(tag, "%s dispatch '%s' fail - %s", dis, n, err)
		return
	}
	if f != nil {
		logger.Debug(tag, "%s dispatch '%s' done - %s", dis, n, f)
		this.executor.DoNow("dispachResult", func() error {
			this.endDispatch(n, f)
			return nil
		})
	}
}

func (this *Router) doDispatch(n string, resp chan espnet.ChannelFactory) error {
	if cf, ok := this.factorys[n]; ok {
		this.safe(resp, cf)
		return nil
	}
	task, ok := this.tasks[n]
	if ok {
		task.respList = append(task.respList, resp)
		return nil
	}
	if len(this.dispatchers) == 0 {
		this.safe(resp, nil)
		return nil
	}

	task = new(dispatchTask)
	task.count = len(this.dispatchers)
	task.respList = make([]chan espnet.ChannelFactory, 1)
	task.respList[0] = resp
	this.tasks[n] = task
	for _, dis := range this.dispatchers {
		go this.callDispatch(dis, n)
	}
	return nil
}

func (this *Router) Dispatch(n string, timeout time.Duration) (espnet.ChannelFactory, error) {
	resp := make(chan espnet.ChannelFactory, 1)
	defer close(resp)
	call := func() error {
		return this.doDispatch(n, resp)
	}

	ev := make(chan error, 1)
	defer close(ev)
	tm := time.NewTimer(timeout)
	defer tm.Stop()
	err2 := this.executor.Do("dispatch", call, qexec.SyncCallback(ev))
	if err2 != nil {
		return nil, err2
	}

	select {
	case err := <-ev:
		if err != nil {
			return nil, err
		}
	case ch := <-resp:
		if ch == nil {
			return nil, errors.New(fmt.Sprintf("can't found %s", n))
		}
		return ch, nil
	case <-tm.C:
		return nil, errors.New("timeout")
	}

	select {
	case ch := <-resp:
		if ch == nil {
			return nil, errors.New(fmt.Sprintf("can't found %s", n))
		}
		return ch, nil
	case <-tm.C:
		return nil, errors.New("timeout")
	}
}

func (this *Router) GetChannel(n string, timeout time.Duration) (espnet.Channel, error) {
	if this.executor.IsClosing() {
		return nil, errors.New("closed")
	}
	m := this.fastCache
	if cf, ok := m[n]; ok {
		return cf.NewChannel()
	}
	// dispatch
	cf, err := this.Dispatch(n, timeout)
	if err != nil {
		return nil, err
	}
	if cf == nil {

	}
	return cf.NewChannel()
}
