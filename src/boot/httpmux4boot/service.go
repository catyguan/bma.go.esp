package httpmux4boot

import (
	"boot"
	"esp/glua"
	"fmt"
	"net/http"
)

func InitMux(mux *http.ServeMux, path string, service *glua.Service) {
	if path == "" {
		path = "/"
	}
	mux.HandleFunc(path+"reload", func(w http.ResponseWriter, req *http.Request) {
		msg := "ok"
		if boot.Restart() {
			msg = "fail"
		}
		fmt.Fprintf(w, msg)
	})
}
