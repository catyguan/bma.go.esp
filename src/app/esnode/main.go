package main

import (
	"boot"
	"boot/httpmux4boot"
	"esp/acclog"
	"esp/glua"
	"esp/glua/httpmux4glua"
	"httpserver"
	"net/http"
	"os"
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

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4glua.NewService("gluaMux", service, httpmux4glua.CommonDispatcher)
	mux4gl.AccLog = acclog
	mux4gl.AccName = "http"
	mux4gl.InitMux(mux, "/")
	mux4gl.InitManageMux(mux, "/m/")
	boot.AddService(mux4gl)

	httpmux4boot.InitMux(mux, "/sys/")

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}
