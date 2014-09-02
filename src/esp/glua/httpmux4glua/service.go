package httpmux4glua

import (
	"bmautil/valutil"
	"esp/acclog"
	"esp/glua"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	tag = "httpmux4glua"
)

type Service struct {
	name    string
	config  *configInfo
	glua    *glua.Service
	AccLog  *acclog.Service
	AccName string
}

func NewService(n string, s *glua.Service) *Service {
	this := new(Service)
	this.name = n
	this.glua = s
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
	path := req.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	cfg := this.config
	if cfg.Location != nil {
		if s, ok := cfg.Location[path]; ok {
			path = s
		}
	}
	if path == "/reload" {
		this.doReload(w, req)
	} else {
		this.doInvoke(w, req, path)
	}
}

func (this *Service) doReload(w http.ResponseWriter, req *http.Request) {
	g := req.FormValue("g")
	l := req.FormValue("l")
	if g == "" || l == "" {
		this.Error(w, -1, "empty param", http.StatusBadRequest)
		return
	}
	gl := this.glua.GetGLua(g)
	if gl == nil {
		this.Error(w, -2, fmt.Sprintf("invalid gl - %s", g), http.StatusBadRequest)
		return
	}
	err2 := gl.ReloadScript(l)
	if err2 != nil {
		this.Error(w, -9, err2.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok")
}

func (this *Service) doInvoke(w http.ResponseWriter, req *http.Request, path string) {
	qpath := strings.TrimPrefix(path, "/")
	qlist := strings.SplitN(qpath, "/", 2)
	if len(qlist) != 2 {
		this.Error(w, -1, fmt.Sprintf("invalid request location - %s", path), http.StatusBadRequest)
		return
	}
	g := qlist[0]
	f := qlist[1]
	if g == "" || f == "" {
		this.Error(w, -1, "empty param", http.StatusBadRequest)
		return
	}
	gl := this.glua.GetGLua(g)
	if gl == nil {
		this.Error(w, -2, fmt.Sprintf("invalid gl - %s", g), http.StatusBadRequest)
		return
	}

	to := valutil.ToInt(req.FormValue("_to"), 0)
	if to <= 0 {
		to = this.config.TimeoutMS
	}
	if to <= 0 {
		to = 5000
	}

	ctx := gl.NewContext(this.config.FuncPrefix + f)
	ctx.Acclog = this.AccLog
	ctx.AccName = this.AccName
	ctx.Timeout = time.Duration(to) * time.Millisecond
	if true {
		dt := make(map[string]interface{})
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
		ctx.Data = dt
	}

	gl.ExecuteSync(ctx)

	if ctx.Error != nil {
		this.Error(w, -9, ctx.Error.Error(), http.StatusInternalServerError)
		return
	}

	res := ctx.Result
	fmt.Println(res)

	whs := w.Header()
	ctype := "text/plain; charset=utf-8"
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

	content := this.config.EmptyContent
	if true {
		if v, ok := res["Content"]; ok {
			content = fmt.Sprintf("%v", v)
		}
	}
	fmt.Fprint(w, content)
}
