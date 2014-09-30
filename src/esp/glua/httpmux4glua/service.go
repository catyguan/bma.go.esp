package httpmux4glua

import (
	"context"
	"esp/acclog"
	"esp/glua"
	"fmt"
	"logger"
	"net/http"
	"strings"
	"time"
)

const (
	tag = "httpmux4glua"
)

type Service struct {
	name       string
	config     *configInfo
	glua       *glua.Service
	dispatcher Dispatcher
	AccLog     *acclog.Service
	AccName    string
}

func NewService(n string, s *glua.Service, dis Dispatcher) *Service {
	this := new(Service)
	this.name = n
	this.glua = s
	this.dispatcher = dis
	return this
}

func (this *Service) InitMux(mux *http.ServeMux, path string) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	mux.Handle(path, http.StripPrefix(path, this))
}

func (this *Service) Error(w http.ResponseWriter, ec int, err string, code int) {
	if this.config.DevMode {
		http.Error(w, fmt.Sprintf("%d: %s", ec, err), code)
	} else {
		http.Error(w, fmt.Sprintf("%d", ec), code)
	}
}

func (this *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if !strings.Contains(host, ":") {
		host = host + ":80"
	}
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	for _, p1 := range this.config.Skip {
		if p1 == path {
			logger.Debug(tag, "global skip %s:%s", host, path)
			http.NotFound(w, req)
			return
		}
	}

	cfg := this.config
	var mapp *configApp
	mlen := 0
	for _, app := range cfg.App {
		if app.Host == "" || app.Host == host {
			if strings.HasPrefix(path, app.Location) {
				if len(app.Location) > mlen {
					mlen = len(app.Location)
					mapp = app
				}
			}
		}
	}
	if mapp == nil {
		logger.Debug(tag, "miss %s:%s", host, path)
		this.Error(w, -1, "invalid glua app", http.StatusBadRequest)
		return
	}
	opath := path
	path = strings.TrimPrefix(path, mapp.Location)
	if path == "" || strings.HasSuffix(path, "/") {
		path = path + mapp.IndexName
	}
	logger.Debug(tag, "match %s:%s => %s:%s:%s", host, opath, mapp.Name, mapp.Location, path)

	for _, p1 := range mapp.Skip {
		if p1 == path {
			logger.Debug(tag, "app skip %s", path)
			http.NotFound(w, req)
			return
		}
	}

	this.doInvoke(w, req, mapp, path)
}

func (this *Service) doInvoke(w http.ResponseWriter, req *http.Request, app *configApp, path string) {
	greq, err0 := this.dispatcher(req, path)
	if err0 != nil {
		this.Error(w, -2, err0.Error(), http.StatusBadRequest)
		return
	}
	if greq.FuncName == "" {
		this.Error(w, -3, "glua func miss", http.StatusBadRequest)
		return
	}
	if app.FuncPrefix != "" {
		greq.FuncName = app.FuncPrefix + greq.FuncName
	}
	g := app.Name
	gl := this.glua.GetGLua(g)
	if gl == nil {
		this.Error(w, -4, fmt.Sprintf("invalid gl - %s", g), http.StatusBadRequest)
		return
	}

	to := greq.timeout
	if to <= 0 {
		to = app.TimeoutMS
	}
	if to <= 0 {
		to = 5000
	}

	ctx := gl.NewContext("", true)
	ainfo, _ := glua.GLuaContext.AcclogInfo(ctx)
	ainfo.Acclog = this.AccLog
	ainfo.AccName = this.AccName

	dt := make(map[string]interface{})
	if true {
		for k, _ := range req.Form {
			if !strings.HasPrefix(k, "_") {
				v := req.FormValue(k)
				dt[k] = v
			}
		}
		hs := make(map[string]string)
		for k, _ := range req.Header {
			v := req.Header.Get(k)
			hs[k] = v
		}
		dt["Header"] = hs
	}
	reload := false
	if this.config.DevMode && this.config.AutoReload {
		reload = true
	}
	lua := glua.NewLuaInfo(greq.Script, greq.FuncName, reload)
	glua.GLuaContext.SetExecuteInfo(ctx, greq.FuncName, lua, dt)

	errE := func() error {
		nctx, cancel := context.WithTimeout(ctx, time.Duration(to)*time.Millisecond)
		defer cancel()
		return gl.ExecuteSync(nctx)
	}()

	if errE != nil {
		glua.GLuaContext.End(ctx, errE)
		this.Error(w, -9, errE.Error(), http.StatusInternalServerError)
		return
	}

	res := glua.GLuaContext.GetResult(ctx)
	// fmt.Println(res)

	whs := w.Header()
	ctype := "text/plain; charset=utf-8"
	if ctypev, ok := res["Content-Type"]; ok {
		ctype = fmt.Sprintf("%v", ctypev)
	}
	if true {
		if v, ok := res["Header"]; ok {
			if hs, ok2 := v.(map[string]interface{}); ok2 {
				for k, rv := range hs {
					sv := fmt.Sprintf("%v", rv)
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

	st := http.StatusOK
	if true {
		if v, ok := res["StatusCode"]; ok {
			if sc, ok2 := v.(int); ok2 {
				st = sc
			}
		}
	}
	w.WriteHeader(st)

	content := app.EmptyContent
	if true {
		if v, ok := res["Content"]; ok {
			content = fmt.Sprintf("%s", v)
		}
	}
	fmt.Fprint(w, content)
}
