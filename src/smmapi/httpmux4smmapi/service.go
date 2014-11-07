package httpmux4smmapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smmapi"
)

func InitMuxInvoke(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		err0 := req.ParseForm()
		if err0 != nil {
			http.Error(w, err0.Error(), http.StatusInternalServerError)
		}

		param := make(map[string]interface{})
		id := req.FormValue("id")
		aid := req.FormValue("aid")
		strparam := req.FormValue("param")
		if strparam != "" {
			err1 := json.Unmarshal([]byte(strparam), &param)
			if err1 != nil {
				http.Error(w, fmt.Sprintf("decode param fail - %s", err1), http.StatusBadRequest)
				return
			}
		}

		result, err := smmapi.Invoke(id, aid, param)

		r := make(map[string]interface{})
		if err != nil {
			r["Status"] = 500
			r["Error"] = err.Error()
		} else {
			r["Status"] = 200
			r["Result"] = result
		}

		jbs, _ := json.Marshal(r)
		w.Write(jbs)
	})
}
