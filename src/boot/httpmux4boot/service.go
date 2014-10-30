package httpmux4boot

import (
	"boot"
	"encoding/json"
	"net/http"
)

func InitMuxReload(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		r := make(map[string]interface{})
		r["Status"] = 200
		r["Result"] = "ok"
		if !boot.Restart() {
			r["Result"] = "fail"
		}
		jbs, _ := json.Marshal(r)
		w.Write(jbs)
	})
}
