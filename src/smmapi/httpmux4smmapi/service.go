package httpmux4smmapi

import (
	"encoding/json"
	"net/http"
	"smmapi"
)

func InitMuxInvoke(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		err0 := req.ParseForm()
		if err0 != nil {
			http.Error(w, err0.Error(), http.StatusInternalServerError)
		}

		id := req.FormValue("id")
		aid := req.FormValue("aid")
		param := req.FormValue("param")

		result, refresh, err := smmapi.Invoke(id, aid, param)
		var info *smmapi.SMInfo
		if refresh {
			o := smmapi.Get(id)
			info, err = o.GetInfo()
		}

		r := make(map[string]interface{})
		if err != nil {
			r["Status"] = 500
			r["Error"] = err.Error()
		} else {
			r["Status"] = 200
			r["Result"] = result
			r["Info"] = info
		}

		jbs, _ := json.Marshal(r)
		w.Write(jbs)
	})
}
