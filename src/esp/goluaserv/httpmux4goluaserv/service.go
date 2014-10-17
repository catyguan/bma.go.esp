package httpmux4goluaserv

import (
	"bmautil/valutil"
	"boot"
	"context"
	"esp/acclog"
	"esp/goluaserv"
	"fmt"
	"golua"
	"logger"
	"net/http"
	"strings"
	"time"
)

const (
	tag = "httpmux4glua"
)

type ServRequest struct {
	golua.RequestInfo
	App     string
	Timeout int
}

type Service struct {
	name    string
	gls     *goluaserv.Service
	accLog  *acclog.Service
	accName string
}

func NewService(n string, s *goluaserv.Service) *Service {
	this := new(Service)
	this.name = n
	this.gls = s
	return this
}

func (this *Service) SetupAcclog(al *acclog.Service, n string) {
	if n == "" {
		n = this.name
	}
	this.accLog = al
	this.accName = n
}

func (this *Service) InitMux(mux *http.ServeMux, path string) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	mux.Handle(path, http.StripPrefix(path, this))
}

func (this *Service) Error(w http.ResponseWriter, err string, code int, adt map[string]interface{}) {
	if adt != nil {
		adt["st"] = code
		adt["err"] = err
	}
	if boot.DevMode {
		http.Error(w, err, code)
	} else {
		http.Error(w, "Sorry, Server Error", code)
	}
}

func (this *Service) Dispatch(ctx context.Context, w http.ResponseWriter, req *http.Request, data map[string]interface{}, hs map[string]interface{}, adt map[string]interface{}) (*ServRequest, bool) {
	gl := this.gls.GetGoLua(this.name)
	if gl == nil {
		this.Error(w, fmt.Sprintf("miss GoLua(%s)", this.name), http.StatusInternalServerError, adt)
		return nil, false
	}

	glreq := golua.NewRequestInfo()
	glreq.Script = "dispatch.lua"
	glreq.Data = make(map[string]interface{})

	host := req.Host
	if !strings.Contains(host, ":") {
		host = host + ":80"
	}
	glreq.Data["host"] = host
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	glreq.Data["path"] = path
	glreq.Data["request"] = data

	glreq.Context = hs

	if boot.DevMode {
		glreq.Reload = true
	}

	ctx = golua.CreateRequest(ctx, glreq)

	r, err2 := gl.Execute(ctx)
	if err2 != nil {
		this.Error(w, fmt.Sprintf("dispatch fail: %s", err2), http.StatusInternalServerError, adt)
		return nil, false
	}

	res, ok := r.(map[string]interface{})
	if !ok {
		this.Error(w, fmt.Sprintf("dispatch invalid: %v", r), http.StatusBadRequest, adt)
		return nil, false
	}
	var sq ServRequest
	if !valutil.ToBean(res, &sq) {
		this.Error(w, fmt.Sprintf("dispatch invalid: %v", res), http.StatusBadRequest, adt)
		return nil, false
	}
	logger.Debug(tag, "dispatch %s:%s => %s:%s", host, path, sq.App, sq.Script)
	if adt != nil {
		adt["app"] = sq.App
		adt["scr"] = sq.Script
	}

	return &sq, true
}

func (this *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var adt map[string]interface{}
	rst := http.StatusOK
	start := time.Now()
	if this.accLog != nil {
		adt = make(map[string]interface{})
		adt["host"] = req.Host
		adt["uri"] = req.RequestURI
	}
	defer func() {
		if adt != nil {
			du := time.Since(start)
			l := acclog.NewCommonLog(adt, du.Seconds())
			this.accLog.Write(this.name, l)
		}
	}()

	errP := req.ParseForm()
	if errP != nil {
		this.Error(w, errP.Error(), http.StatusBadRequest, adt)
		return
	}

	dt := make(map[string]interface{})
	hs := make(map[string]interface{})
	if true {
		for k, _ := range req.Form {
			v := req.FormValue(k)
			dt[k] = v
		}
		for k, _ := range req.Header {
			v := req.Header.Get(k)
			hs[k] = v
		}
	}
	ctx := context.Background()
	ctx, _ = context.CreateExecId(ctx)
	sq, ok := this.Dispatch(ctx, w, req, dt, hs, adt)
	if !ok {
		return
	}

	if sq.Data != nil {
		for k, v := range dt {
			if _, ok := sq.Data[k]; !ok {
				sq.Data[k] = v
			}
		}
	} else {
		sq.Data = dt
	}
	if sq.Context != nil {
		for k, v := range hs {
			if _, ok := sq.Context[k]; !ok {
				sq.Context[k] = v
			}
		}
	} else {
		sq.Context = hs
	}

	if boot.DevMode {
		sq.Reload = true
	}

	gl := this.gls.GetGoLua(sq.App)
	if gl == nil {
		this.Error(w, fmt.Sprintf("%s invalid App - %s", req.URL.Path, sq.App), http.StatusBadRequest, adt)
		return
	}

	ctx = golua.CreateRequest(ctx, &sq.RequestInfo)
	if sq.Timeout > 0 {
		nctx, cancel := context.WithTimeout(ctx, time.Duration(sq.Timeout)*time.Millisecond)
		defer cancel()
		ctx = nctx
	}
	r, errE := gl.Execute(ctx)

	if errE != nil {
		this.Error(w, errE.Error(), http.StatusInternalServerError, adt)
		return
	}

	res, ok := r.(map[string]interface{})
	if !ok {
		if str, ok2 := r.(string); ok2 {
			res = make(map[string]interface{})
			res["Content"] = str
		} else {
			this.Error(w, fmt.Sprintf("response invalid: %v", r), http.StatusInternalServerError, adt)
			return
		}
	}
	if v, ok := res["Status"]; ok {
		rst = valutil.ToInt(v, http.StatusOK)
	}

	rcontent := ""
	if true {
		if v, ok := res["Content"]; ok {
			rcontent = valutil.ToString(v, "")
		}
	}

	if logger.EnableDebug(tag) {
		logger.Debug(tag, "dispatch %s:%s => Status=%d:Size=%d", sq.App, sq.Script, rst, len(rcontent))
	}
	// fmt.Println(res)

	whs := w.Header()
	ctype := "text/plain; charset=utf-8"
	if ctypev, ok := res["Content-Type"]; ok {
		ctype = valutil.ToString(ctypev, ctype)
	}
	if true {
		if v, ok := res["Header"]; ok {
			if hs, ok2 := v.(map[string]interface{}); ok2 {
				for k, rv := range hs {
					sv := valutil.ToString(rv, "")
					if k == "Content-Type" {
						ctype = sv
						continue
					}
					if sv == "" {
						whs.Set(k, sv)
					}
				}
			}
		}
	}
	whs.Set("Content-Type", ctype)

	w.WriteHeader(rst)

	if adt != nil {
		adt["status"] = rst
	}

	fmt.Fprint(w, rcontent)
}
