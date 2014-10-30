package httpmux4goluaserv

import (
	"encoding/json"
	"net/http"
)

func (this *Service) InitMuxReset(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
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
	r := make(map[string]interface{})
	r["Status"] = 200
	r["Result"] = "ok"
	if !ok {
		r["Result"] = "fail"
	}
	jbs, _ := json.Marshal(r)
	w.Write(jbs)
}
