package httpmux4smmapi

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"logger"
	"net/http"
	"smmapi"
	"strings"
)

type Service struct {
	name   string
	config *configInfo
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	return r
}

func (this *Service) InitMuxInvoke(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		err0 := req.ParseForm()
		if err0 != nil {
			http.Error(w, err0.Error(), http.StatusInternalServerError)
			return
		}

		param := make(map[string]interface{})
		id := req.FormValue("id")
		aid := req.FormValue("aid")
		strparam := req.FormValue("param")
		code := strings.ToLower(req.FormValue("code"))
		mycode := ""
		cfgcode := this.config.Code
		if cfgcode != "" {
			tmp := fmt.Sprintf("%s/%s/%s/%s", id, aid, strparam, cfgcode)
			h := md5.New()
			h.Write([]byte(tmp))
			mycode = fmt.Sprintf("%x", h.Sum(nil))
		}
		if code != mycode {
			logger.Warn(tag, "'%s' code invalid(in=%s, my=%s)", req.RemoteAddr, code, mycode)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
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
