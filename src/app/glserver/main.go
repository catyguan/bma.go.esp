package main

import (
	"boot"
	"esp/acclog"
	"esp/goluaserv"
	"esp/goluaserv/httpmux4goluaserv"
	"fileloader"
	"golua"
	"golua/vmmacclog"
	"golua/vmmhttp"
	"golua/vmmjson"
	"golua/vmmsql"
	"httpserver"
	"net/http"
	"os"

	_ "fileloader/http4fileloader"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tag = "glserver"
)

func main() {
	cfile := "config/glserver-config.json"

	acclog := acclog.NewService("acclog")
	boot.AddService(acclog)

	fl := fileloader.NewService("fileloader")
	boot.AddService(fl)

	service := goluaserv.NewService("goluaServ", func(vmg *golua.VMG) {
		myInitor(vmg, acclog)
	})
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4goluaserv.NewService("goluaMux", service)
	mux4gl.SetupAcclog(acclog, "httpserv")
	mux4gl.InitMux(mux, "/")
	boot.AddService(mux4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}

func myInitor(vmg *golua.VMG, acclog *acclog.Service) {
	golua.CoreModule(vmg)
	golua.GoModule().Bind(vmg)
	golua.TypesModule().Bind(vmg)
	golua.TableModule().Bind(vmg)
	golua.StringsModule().Bind(vmg)
	golua.TimeModule().Bind(vmg)
	vmmhttp.HttpServModule().Bind(vmg)
	vmmhttp.HttpClientModule(acclog, "httpclient").Bind(vmg)
	vmmacclog.Module().Bind(vmg)
	vmmjson.Module().Bind(vmg)
	vmmsql.Module().Bind(vmg)
}
