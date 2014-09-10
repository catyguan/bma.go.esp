package main

import (
	"boot"
	"boot/httpmux4boot"
	"esp/acclog"
	"esp/glua"
	"esp/glua/httpmux4glua"
	"httpserver"
	"net/http"
)

const (
	tag = "esnode"
)

func main() {
	cfile := "config/esnode-config.json"

	acclog := acclog.NewService("acclog")
	boot.AddService(acclog)

	service := glua.NewService("gluaService")
	boot.AddService(service)

	// var wd, _ = os.Getwd()

	// mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4glua.NewService("gluaMux", service, httpmux4glua.CommonDispatcher)
	mux4gl.AccLog = acclog
	mux4gl.AccName = "http"
	boot.AddService(mux4gl)

	if true {
		fileServer := httpserver.NewHttpFileServer("fileServer")
		boot.AddService(fileServer)

		httpService := httpserver.NewHttpServer("httpPoint", nil)
		httpService.Add(fileServer.BuildMux)
		httpService.Add(func(mux *http.ServeMux) {
			mux4gl.InitMux(mux, "/")
		})
		boot.AddService(httpService)
	}

	if true {
		fileServer := httpserver.NewHttpFileServer("manageFileServer")
		boot.AddService(fileServer)

		httpService := httpserver.NewHttpServer("manageHttpPoint", nil)
		httpService.Add(fileServer.BuildMux)
		httpService.Add(func(mux *http.ServeMux) {
			mux4gl.InitManageMux(mux, "/m/")
			httpmux4boot.InitMux(mux, "/sys/")
		})
		boot.AddService(httpService)

	}

	boot.Go(cfile)
}
