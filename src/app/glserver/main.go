package main

import (
	"boot"
	"esp/acclog"
	"esp/goluaserv"
	"esp/goluaserv/httpmux4goluaserv"
	"golua"
	"golua/vmmacclog"
	"golua/vmmhttp"
	"golua/vmmjson"
	"httpserver"
	"net/http"
	"os"
)

const (
	tag = "glserver"
)

func main() {
	cfile := "config/glserver-config.json"

	acclog := acclog.NewService("acclog")
	boot.AddService(acclog)

	service := goluaserv.NewService("goluaServ", myInitor)
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4goluaserv.NewService("goluaMux", service)
	mux4gl.SetupAcclog(acclog, "http")
	mux4gl.InitMux(mux, "/")
	boot.AddService(mux4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}

func myInitor(vmg *golua.VMG) {
	golua.CoreModule(vmg)
	golua.GoModule().Bind(vmg)
	golua.TypesModule().Bind(vmg)
	golua.TableModule().Bind(vmg)
	golua.StringsModule().Bind(vmg)
	vmmhttp.HttpServModule().Bind(vmg)
	vmmacclog.Module().Bind(vmg)
	vmmjson.Module().Bind(vmg)
}
