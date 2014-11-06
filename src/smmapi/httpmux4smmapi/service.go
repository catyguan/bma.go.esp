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

		var id, aid, param string
		ctx := make(map[string]interface{})
		for k, _ := range req.Form {
			v := req.FormValue(k)
			switch k {
			case "_id":
				id = v
			case "_aid":
				aid = v
			case "_param":
				param = v
			default:
				ctx[k] = v
			}
		}

		result, err := smmapi.Invoke(id, aid, param, ctx)

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
