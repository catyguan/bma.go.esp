package httpmux4glua

import (
	"esp/glua"
	"fmt"
	"net/http"
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

	})
}
