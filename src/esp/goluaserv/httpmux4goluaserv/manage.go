package httpmux4goluaserv

import (
	"fmt"
	"net/http"
)

func (this *Service) InitManageMux(mux *http.ServeMux, path string) {
	if path == "" {
		path = "/"
	}
	mux.HandleFunc(path+"reset", func(w http.ResponseWriter, req *http.Request) {
		this.doReset(w, req)
	})
}

func (this *Service) doReset(w http.ResponseWriter, req *http.Request) {
	g := req.FormValue("g")
	if g == "" {
		this.Error(w, "empty param", http.StatusBadRequest, nil)
		return
	}
	ok := this.gls.ResetGoLua(g)
	resp := "ok"
	if !ok {
		resp = "fail"
	}
	fmt.Fprintf(w, resp)
}
