package httpmux4goluaserv

import (
	"bmautil/httputil"
	"bmautil/valutil"
	"boot"
	"context"
	"esp/acclog"
	"esp/goluaserv"
	"fmt"
	"golua"
	"golua/vmmhttp"
	"logger"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	tag = "httpmux4glua"
)

type ServRequest struct {
	golua.RequestInfo
	app *configApp
}

type Service struct {
	name    string
	config  *configInfo
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
		http.Error(w, http.StatusText(code), code)
	}
}

func (this *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var adt map[string]interface{}
	start := time.Now()
	if this.accLog != nil {
		adt = make(map[string]interface{})
		adt["host"] = req.Host
		adt["uri"] = req.RequestURI
		adt["status"] = http.StatusOK
	}
	defer func() {
		if adt != nil {
			du := time.Since(start)
			l := acclog.NewCommonLog(adt, du.Seconds())
			this.accLog.Write(this.name, l)
		}
	}()

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
		this.Error(w, fmt.Sprintf("can't dispatch '%s' to golua app", req.RequestURI), http.StatusBadRequest, adt)
		return
	}
	opath := path
	path = strings.TrimPrefix(path, mapp.Location)
	if path == "" || strings.HasSuffix(path, "/") {
		path = path + mapp.IndexName
	}
	logger.Debug(tag, "dispatch %s:%s => %s, %s, %s", host, opath, mapp.Name, mapp.Location, path)

	for _, p1 := range mapp.Skip {
		if p1 == path {
			logger.Debug(tag, "app skip %s", path)
			http.NotFound(w, req)
			return
		}
	}
	req.URL.Path = path

	this.doInvoke(w, req, mapp, path, adt)
}

func (this *Service) doInvoke(w http.ResponseWriter, req *http.Request, app *configApp, path string, adt map[string]interface{}) {

	errP := httputil.Prepare(req, this.config.ParseFormMaxMemory)
	if errP != nil {
		this.Error(w, errP.Error(), http.StatusBadRequest, adt)
		return
	}

	ri := golua.NewRequestInfo()
	ri.Script = app.Script

	gl := this.gls.GetGoLua(app.Name)
	if gl == nil {
		this.Error(w, fmt.Sprintf("%s invalid App - %s", req.URL.Path, app.Name), http.StatusBadRequest, adt)
		return
	}

	ctx := context.Background()
	ctx, _ = context.CreateExecId(ctx)
	ctx = acclog.CreateAcclogData(ctx, adt)
	ctx = vmmhttp.CreateServ(ctx, w, req)
	ctx = golua.CreateRequest(ctx, ri)

	if app.TimeoutMS > 0 {
		nctx, cancel := context.WithTimeout(ctx, time.Duration(app.TimeoutMS)*time.Millisecond)
		defer cancel()
		ctx = nctx
	}
	r, errE := gl.Execute(ctx)

	if errE != nil {
		this.Error(w, errE.Error(), http.StatusInternalServerError, adt)
		return
	}

	if r != nil {
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
		rst := http.StatusOK
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
			logger.Debug(tag, "execute %s:%s => Status=%d:Size=%d", app.Name, app.Script, rst, len(rcontent))
		}
		// fmt.Println(res)

		whs := w.Header()
		ctype := ""
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
		if ctype == "" {
			if cfile, ok := res["ContentFile"]; ok {
				scfile := valutil.ToString(cfile, "")
				if scfile != "" {
					ctype = mime.TypeByExtension(filepath.Ext(scfile))
				}
			}
		}
		if ctype == "" {
			ctype = "text/plain; charset=utf-8"
		}
		whs.Set("Content-Type", ctype)

		w.WriteHeader(rst)

		if adt != nil {
			adt["status"] = rst
		}

		fmt.Fprint(w, rcontent)
	}
}
