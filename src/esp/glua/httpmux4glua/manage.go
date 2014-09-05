package httpmux4glua

import (
	"fmt"
	"net/http"
)

func (this *Service) InitManageMux(mux *http.ServeMux, path string) {
	if path == "" {
		path = "/"
	}
	mux.HandleFunc(path+"reload", func(w http.ResponseWriter, req *http.Request) {
		this.doReload(w, req)
	})
	mux.HandleFunc(path+"reset", func(w http.ResponseWriter, req *http.Request) {
		this.doReset(w, req)
	})
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

func (this *Service) doReset(w http.ResponseWriter, req *http.Request) {
	g := req.FormValue("g")
	if g == "" {
		this.Error(w, -1, "empty param", http.StatusBadRequest)
		return
	}
	ok := this.glua.ResetGLua(g)
	resp := "ok"
	if !ok {
		resp = "fail"
	}
	fmt.Fprintf(w, resp)
}
