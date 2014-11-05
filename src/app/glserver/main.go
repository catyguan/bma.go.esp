package main

import (
	"boot"
	"esp/acclog"
	"esp/aclserv"
	"esp/goluaserv"
	"esp/goluaserv/httpmux4goluaserv"
	"fileloader"
	"golua"
	"golua/vmmacclog"
	"golua/vmmclass"
	"golua/vmmhttp"
	"golua/vmmjson"
	"golua/vmmsql"
	"httpserver"
	"httpserver/aclmux"
	"net/http"
	"os"
	"smmapi/httpmux4smmapi"
	_ "smmapi/smmapi4server"

	_ "fileloader/http4fileloader"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tag = "glserver"
)

func main() {
	cfile := "config/glserver-config.json"

	acls := aclserv.NewService("acl")
	boot.AddService(acls)

	acclog := acclog.NewService("acclog")
	boot.AddService(acclog)

	fl := fileloader.NewService("fileloader")
	boot.AddService(fl)

	service := goluaserv.NewService("goluaServ", func(gl *golua.GoLua) {
		myInitor(gl, acclog)
	})
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()

	// httpmux4boot.InitMuxReload(mux, "/smm.api/boot/reload")

	httpmux4smmapi.InitMuxInvoke(mux, "/smm.api/invoke")

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4goluaserv.NewService("goluaMux", service)
	mux4gl.SetupAcclog(acclog, "httpserv")
	mux4gl.InitMux(mux, "/")
	boot.AddService(mux4gl)

	mux4gl.InitMuxReset(mux, "/smm.api/golua/reset")

	amux := aclmux.NewAclServerMux("http", mux)

	httpService := httpserver.NewHttpServer("httpPoint", amux)
	boot.AddService(httpService)

	boot.Go(cfile)
}

func myInitor(gl *golua.GoLua, acclog *acclog.Service) {
	golua.InitCoreLibs(gl)
	vmmhttp.InitGoLuaWithHttpServ(gl)
	vmmhttp.InitGoLuaWithHttpClient(gl, acclog, "httpclient")
	vmmacclog.InitGoLua(gl)
	vmmjson.InitGoLua(gl)
	vmmsql.InitGoLua(gl)
	vmmclass.InitGoLua(gl)
}
