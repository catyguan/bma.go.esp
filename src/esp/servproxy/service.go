package servproxy

import (
	"bmautil/valutil"
	"context"
	"esp/goluaserv"
	"fmt"
	"golua"
	"logger"
	"sync"
)

const (
	tag = "servproxy"
)

type Service struct {
	name    string
	gls     *goluaserv.Service
	config  *configInfo
	lock    sync.RWMutex
	ports   map[string]*PortObj
	targets map[string]*TargetObj
}

func NewService(n string, gls *goluaserv.Service) *Service {
	r := new(Service)
	r.name = n
	r.gls = gls
	r.ports = make(map[string]*PortObj)
	r.targets = make(map[string]*TargetObj)
	return r
}

func (this *Service) GetGoLuaService() *goluaserv.Service {
	return this.gls
}

func (this *Service) _createPort(n string, cfg *PortConfigInfo) error {
	if _, ok := this.ports[n]; ok {
		return nil
	}
	h, err0 := AssertPortHandler(cfg.Type)
	if err0 != nil {
		return err0
	}
	p := NewPortObj(this, n, cfg, h)
	this.ports[n] = p
	err1 := p.Start()
	if err1 != nil {
		p.Stop()
		delete(this.ports, n)
		return err1
	}
	return nil
}

func (this *Service) _createTarget(n string, cfg *TargetConfigInfo) error {
	if _, ok := this.targets[n]; ok {
		return nil
	}
	h, err0 := AssertRemoteHandler(cfg.Type)
	if err0 != nil {
		return err0
	}
	o := NewTargetObj(this, n, cfg, h)
	this.targets[n] = o
	err1 := o.Start()
	if err1 != nil {
		o.Stop()
		delete(this.targets, n)
		return err1
	}
	return nil
}

func (this *Service) _removePort(n string) error {
	p, ok := this.ports[n]
	if !ok {
		return nil
	}
	delete(this.ports, n)
	p.Stop()
	return nil
}

func (this *Service) _removeTarget(n string) {
	o, ok := this.targets[n]
	if !ok {
		return
	}
	delete(this.targets, n)
	o.Stop()
}

func (this *Service) RemovePort(n string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	this._removePort(n)
	return nil
}

func (this *Service) RemoveTarget(n string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	this._removeTarget(n)
	return nil
}

func (this *Service) Execute(port *PortObj, o golua.VMTable, req ProxyRequest) (interface{}, error) {
	cfg := port.cfg

	ri := golua.NewRequestInfo()
	ri.Script = cfg.Script

	gl, errG := this.gls.GetGoLua(cfg.GoLua)
	if errG != nil {
		return nil, errG
	}
	if gl == nil {
		return nil, fmt.Errorf("invalid GoLua App - %s", cfg.GoLua)
	}

	ctx := context.Background()
	ctx, _ = context.CreateExecId(ctx)
	ctx = golua.CreateRequest(ctx, ri)

	locals := make(map[string]interface{})
	locals["request"] = o
	r, errE := gl.DoExecute(ctx, locals)
	if errE != nil {
		return nil, port.handler.AnswerError(port, req, errE)
	}
	if r != nil {
		res, ok := r.(map[string]interface{})
		if !ok {
			return r, nil
		}
		act := ""
		if action, ok := res["Action"]; ok {
			act = valutil.ToString(action, "")
		}
		if act == "" {
			return r, nil
		}

		switch act {
		case "forward":
			tar := ""
			if target, ok := res["Target"]; ok {
				tar = valutil.ToString(target, "")
			}
			if tar == "" {
				return nil, fmt.Errorf("invalid forward Target - %v", res)
			}
			write := valutil.ToBool(res["Write"], false)
			errF := this.DoForward(port, tar, req, write)
			if errF != nil {
				return nil, port.handler.AnswerError(port, req, errF)
			}
			return true, nil
		}
	}
	return r, nil
}

func (this *Service) Select(tar string, write bool) (*RemoteObj, error) {
	var robj *RemoteObj
	this.lock.RLock()
	defer this.lock.RUnlock()
	tobj, ok := this.targets[tar]
	if ok {
		p := 0
		for _, ro := range tobj.remotes {
			if !ro.Ping() {
				continue
			}
			if write && ro.cfg.ReadOnly {
				continue
			}
			if robj == nil {
				robj = ro
				p = ro.cfg.Priority
			} else {
				if ro.cfg.Priority > p {
					robj = ro
					p = ro.cfg.Priority
				}
			}
		}
	}
	return robj, nil
}

func (this *Service) DoForward(port *PortObj, tar string, req ProxyRequest, write bool) error {
	for {
		robj, err := this.Select(tar, write)
		if err != nil {
			return err
		}
		if robj == nil {
			return fmt.Errorf("no usable Remote for target(%s)", tar)
		}
		logger.Debug(tag, "forward '%s' to Remote(%s)", req, robj.name)
		session, err1 := robj.handler.Begin(robj)
		if err1 != nil {
			robj.Fail()
			logger.Debug(tag, "Remote(%s) begin session fail - %s", robj.name, err1)
			continue
		}
		defer session.Finish()

		err2 := req.BeginRead()
		if err2 != nil {
			return err2
		}
		writed := false
		wdone := false
		for {
			ok, data, err3 := req.Read()
			if err3 != nil {
				if writed {
					session.Fail()
				}
				return err3
			}
			if !ok {
				wdone = true
				break
			}
			err4 := session.Write(data)
			if err4 != nil {
				session.Fail()
				logger.Debug(tag, "Remote(%s) session write fail - %s", robj.name, err4)
				break
			}
			writed = true
		}
		if !wdone {
			continue
		}
		if !req.HasResponse() {
			return nil
		}

		for {
			ok, data, err3 := session.Read()
			if err3 != nil {
				session.Fail()
				return err3
			}
			if !ok {
				break
			}
			err4 := port.handler.Write(port, req, data)
			if err4 != nil {
				return err4
			}
		}
		return nil
	}
}
