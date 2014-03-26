package httpmux4glua

import (
	"bmautil/valutil"
	"esp/glua"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func InitMux(mux *http.ServeMux, path string, service *glua.Service) {
	if path == "" {
		path = "/"
	}
	mux.HandleFunc(path+"restart", func(w http.ResponseWriter, req *http.Request) {

	})
	mux.HandleFunc(path+"reload", func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("parse form - %s", err), http.StatusInternalServerError)
			return
		}
		g := req.FormValue("g")
		l := req.FormValue("l")
		if g == "" || l == "" {
			http.Error(w, "empty param", http.StatusBadRequest)
			return
		}
		gl := service.GetGLua(g)
		if gl == nil {
			http.Error(w, "invalid g", http.StatusBadRequest)
			return
		}
		err2 := gl.ReloadScript(l)
		if err2 != nil {
			http.Error(w, fmt.Sprintf("reload fail - %s", err2), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "ok")
	})
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("parse form - %s", err), http.StatusInternalServerError)
			return
		}
		qpath := strings.TrimPrefix(req.URL.Path, path)
		qlist := strings.SplitN(qpath, "/", 2)
		if len(qlist) != 2 {
			http.Error(w, fmt.Sprintf("invalid request uri - %s", req.RequestURI), http.StatusBadRequest)
			return
		}
		g := qlist[0]
		f := qlist[1]
		gl := service.GetGLua(g)
		if gl == nil {
			http.Error(w, fmt.Sprintf("invalid gl - %s", g), http.StatusBadRequest)
			return
		}

		to := valutil.ToInt(req.FormValue("_to"), 0)
		if to <= 0 {
			to = 5000
		}

		ctx := gl.NewContext("service_" + f)
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
			http.Error(w, ctx.Error.Error(), http.StatusInternalServerError)
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

		content := "<empty response>"
		if true {
			if v, ok := res["Content"]; ok {
				content = fmt.Sprintf("%v", v)
			}
		}
		fmt.Fprint(w, content)
	})
}
