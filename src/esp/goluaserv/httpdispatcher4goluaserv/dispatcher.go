package httpmux4goluaserv

import (
	"bmautil/valutil"
	"boot"
	"context"
	"esp/goluaserv"
	"fmt"
	"golua"
	"net/http"
	"strings"
)

type disInfo struct {
	Path string
}

type Dispatcher struct {
	name    string
	gls     *goluaserv.Service
	handler *http.ServeMux
}

func NewDispatcher(n string, s *goluaserv.Service, mux *http.ServeMux) *Dispatcher {
	this := new(Dispatcher)
	this.name = n
	this.gls = s
	this.handler = mux
	return this
}

func (this *Dispatcher) InitMux(mux *http.ServeMux, path string) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	mux.Handle(path, http.StripPrefix(path, this))
}

func (this *Dispatcher) Error(w http.ResponseWriter, err string, code int) {
	if boot.DevMode {
		http.Error(w, err, code)
	} else {
		http.Error(w, "Sorry, Server Error", code)
	}
}

func (this *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	gl := this.gls.GetGoLua(this.name)
	if gl == nil {
		this.Error(w, fmt.Sprintf("miss GoLua(%s)", this.name), http.StatusInternalServerError)
		return
	}

	errP := req.ParseForm()
	if errP != nil {
		this.Error(w, errP.Error(), http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	host := req.Host
	if !strings.Contains(host, ":") {
		host = host + ":80"
	}
	data["host"] = host
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	data["path"] = path

	glreq := golua.NewRequestInfo()
	if true {
		dt := make(map[string]interface{})
		for k, _ := range req.Form {
			v := req.FormValue(k)
			dt[k] = v
		}
		data["request"] = dt
		glreq.Data = dt

		hs := make(map[string]interface{})
		for k, _ := range req.Header {
			v := req.Header.Get(k)
			hs[k] = v
		}
		glreq.Context = hs
	}
	if boot.DevMode {
		glreq.Reload = true
	}

	ctx := context.Background()
	ctx, _ = context.CreateExecId(ctx)
	ctx = golua.CreateRequest(ctx, glreq)

	r, err2 := gl.Execute(ctx)
	if err2 != nil {
		this.Error(w, fmt.Sprintf("dispatch fail: %s", err2), http.StatusInternalServerError)
		return
	}

	res, ok := r.(map[string]interface{})
	if !ok {
		this.Error(w, fmt.Sprintf("dispatch invalid: %v", r), http.StatusBadRequest)
		return
	}
	var disinfo disInfo
	if !valutil.ToBean(res, &disinfo) {
		this.Error(w, fmt.Sprintf("dispatch invalid: %v", res), http.StatusBadRequest)
		return
	}

	req.URL.Path = disinfo.Path
	this.handler.ServeHTTP(w, req)
}
