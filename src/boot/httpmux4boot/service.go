package httpmux4boot

import (
	"boot"
	"fmt"
	"net/http"
)

func InitMux(mux *http.ServeMux, path string) {
	if path == "" {
		path = "/"
	}
	mux.HandleFunc(path+"reload", func(w http.ResponseWriter, req *http.Request) {
		msg := "ok"
		if !boot.Restart() {
			msg = "fail"
		}
		fmt.Fprintf(w, msg)
	})
}
