package main

import (
	"boot"
	"esp/acclog"
	"esp/aclserv"
	"esp/espnet/vmmesnp"
	"esp/goluaserv"
	"esp/goluaserv/httpmux4goluaserv"
	"esp/memserv"
	"esp/memserv/vmmmemserv"
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
	_ "smmapi/smmapi4config"
	_ "smmapi/smmapi4server"

	_ "fileloader/http4fileloader"
	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-sqlite3"
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

	mems := memserv.NewMemoryServ()
	bwmems := new(boot.BootWrap)
	bwmems.SetCleanup(func() bool {
		mems.CloseAll(true)
		return true
	})
	boot.AddService(bwmems)
	mems.InitSMMAPI("go.memserv")

	service := goluaserv.NewService("goluaServ", func(gl *golua.GoLua) {
		myInitor(gl, acclog, mems)
	})
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux4app := http.NewServeMux()
	mux4app.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4goluaserv.NewService("goluaMux", service)
	mux4gl.SetupAcclog(acclog, "httpserv")
	mux4gl.InitMux(mux4app, "/")
	boot.AddService(mux4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux4app)
	boot.AddService(httpService)

	mux4smm := http.NewServeMux()
	smmapis := httpmux4smmapi.NewService("smmapiServ")
	boot.AddService(smmapis)
	smmapis.InitMuxInvoke(mux4smm, "/smm.api/invoke")

	rmux4smm := aclmux.NewAclServerMux("http", mux4smm)
	httpServiceSMM := httpserver.NewHttpServer("httpPointSMM", rmux4smm)
	boot.AddService(httpServiceSMM)

	boot.Go(cfile)
}

func myInitor(gl *golua.GoLua, acclog *acclog.Service, mems *memserv.MemoryServ) {
	golua.InitCoreLibs(gl)
	vmmhttp.InitGoLuaWithHttpServ(gl)
	vmmhttp.InitGoLuaWithHttpClient(gl, acclog, "httpclient")
	vmmacclog.InitGoLua(gl)
	vmmjson.InitGoLua(gl)
	vmmsql.InitGoLua(gl)
	vmmclass.InitGoLua(gl)
	vmmesnp.InitGoLua(gl)
	vmmmemserv.InitGoLua(gl, mems)
}
