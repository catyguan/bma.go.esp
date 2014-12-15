package servproxy

import (
	"bmautil/valutil"
	"context"
	"esp/goluaserv"
	"fmt"
	"golua"
	"logger"
	"sync"
	"time"
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

	dl := req.Deadline()
	if dl.IsZero() {
		dl = time.Now().Add(time.Duration(port.cfg.TimeoutMS) * time.Millisecond)
		req.SetDeadline(dl)
	}
	tmdu := dl.Sub(time.Now())
	nctx, cancel := context.WithTimeout(ctx, tmdu)
	defer cancel()
	ctx = nctx

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
			errF := this.DoForward(port, tar, req)
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

func (this *Service) Port2Remote(port *PortObj, req ProxyRequest, rname string, session RemoteSession) (returnErr error, retry bool) {
	err01 := req.BeginRead()
	if err01 != nil {
		return err01, false
	}
	defer req.EndRead()

	err02 := session.BeginWrite()
	if err02 != nil {
		session.Fail()
		logger.Debug(tag, "Remote(%s) session beginWrite fail - %s", rname, err02)
		return err02, true
	}
	defer session.EndWrite()

	writed := false
	for {
		ok, data, err3 := req.Read()
		if err3 != nil {
			if writed {
				session.ForceClose()
			}
			return err3, false
		}
		if !ok {
			break
		}
		err4 := session.Write(data)
		if err4 != nil {
			session.Fail()
			logger.Debug(tag, "Remote(%s) session write fail - %s", rname, err4)
			return err4, true
		}
		writed = true
	}
	return nil, false
}

func (this *Service) Remote2Port(port *PortObj, req ProxyRequest, rname string, session RemoteSession) (returnErr error, retry bool) {
	err01 := session.BeginRead(req.Deadline())
	if err01 != nil {
		session.Fail()
		logger.Debug(tag, "Remote(%s) session beginRead fail - %s", rname, err01)
		return err01, true
	}
	defer session.EndRead()

	err02 := port.handler.BeginWrite(port, req)
	if err02 != nil {
		session.ForceClose()
		return err02, false
	}
	defer port.handler.EndWrite(port, req)

	for {
		ok, data, err3 := session.Read()
		if err3 != nil {
			session.Fail()
			return err3, true
		}
		if !ok {
			break
		}
		err4 := port.handler.Write(port, req, data)
		if err4 != nil {
			session.ForceClose()
			return err4, false
		}
	}
	return nil, false
}

func (this *Service) PortForwardRemote(port *PortObj, rname string, req ProxyRequest, session RemoteSession) (returnErr error, retry bool) {
	defer session.Finish()
	logger.Debug(tag, "'%s' -> '%s' port2remote ...", req, rname)
	err1, retry1 := this.Port2Remote(port, req, rname, session)
	if err1 != nil {
		logger.Debug(tag, "'%s' -> '%s' port2remote fail - %v, %s", req, rname, retry1, err1)
		return err1, retry1
	}
	if req.CheckFlag(PRF_NO_RESPONSE) {
		logger.Debug(tag, "'%s' -> '%s' no response, forward done", req, rname)
		return nil, false
	}

	logger.Debug(tag, "'%s' -> '%s' remote2port ...", req, rname)
	err2, retry2 := this.Remote2Port(port, req, rname, session)
	if err2 != nil {
		logger.Debug(tag, "'%s' -> '%s' remote2port fail - %v, %s", req, rname, retry2, err2)
		if retry2 {
			if req.CheckFlag(PRF_WRITE) {
				// don't retry on write operation
				logger.Debug(tag, "'%s' -> '%s' skip failover on write operate", req, rname)
				retry2 = false
			}
		}
		return err2, retry2
	}
	logger.Debug(tag, "'%s' -> '%s' forward done", req, rname)
	return nil, false
}

func (this *Service) DoForward(port *PortObj, tar string, req ProxyRequest) error {
	defer req.Finish()
	write := req.CheckFlag(PRF_WRITE)
	for {
		robj, err := this.Select(tar, write)
		if err != nil {
			return err
		}
		if robj == nil {
			return fmt.Errorf("no usable Remote for target(%s)", tar)
		}
		rname := robj.name
		logger.Debug(tag, "forward '%s' to Remote(%s)", req, rname)
		session, err1 := robj.handler.Begin(robj, RequestTimeout(req))
		if err1 != nil {
			robj.Fail()
			logger.Debug(tag, "Remote(%s) begin session fail - %s", rname, err1)
			continue
		}
		logger.Debug(tag, "Remote(%s) begin session(%s)", rname, session)
		err2, retry := this.PortForwardRemote(port, rname, req, session)
		if err2 != nil {
			if retry {
				continue
			}
			return err2
		}
		return nil
	}
}
