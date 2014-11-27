package main

import (
	"boot"
	"esp/acclog"
	"esp/aclserv"
	"esp/espnet/espnetss"
	"esp/espnet/mempipeline"
	"esp/espnet/vmmesnp"
	"esp/goluaserv"
	"esp/http4goluaserv"
	"esp/memserv"
	"esp/memserv/memserv4httpsession"
	"esp/memserv/vmmmemserv"
	"esp/servicecall"
	"esp/servicecall/vmmservicecall"
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
	_ "fileloader/sql4fileloader"
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

	mems := memserv.NewMemoryServ("memServ")
	boot.AddService(mems)
	mems.InitSMMAPI("go.memserv")

	sess := memserv4httpsession.NewService("sessionServ", mems)
	boot.AddService(sess)

	memp := mempipeline.NewService()
	boot.AddService(memp.CreateBootService("mempipeline"))

	ssp := espnetss.NewService()
	boot.AddService(ssp.CreateBootService("espnetss"))

	scs := servicecall.NewService("serviceCall", nil)
	scs.AddServiceCallerFactory("http", servicecall.HttpServiceCallerFactory(0))
	if true {
		fac := new(servicecall.ESNPMemPipelineServiceCallerFactory)
		fac.S = memp
		scs.AddServiceCallerFactory("memp", fac)
	}
	if true {
		fac := new(servicecall.ESNPNetServiceCallerFactory)
		fac.S = ssp
		scs.AddServiceCallerFactory("esnp", fac)
	}
	boot.AddService(scs)

	service := goluaserv.NewService("goluaServ", func(gl *golua.GoLua) {
		myInitor(gl, acclog, mems, sess, scs)
	})
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux4app := http.NewServeMux()
	mux4app.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	http4gl := http4goluaserv.NewService("goluaHttp", service)
	http4gl.SetupAcclog(acclog, "httpserv")
	http4gl.InitMux(mux4app, "/")
	boot.AddService(http4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux4app)
	boot.AddService(httpService)

	mux4smm := http.NewServeMux()
	smmapis := httpmux4smmapi.NewService("smmapiServ")
	boot.AddService(smmapis)
	smmapis.InitMuxInvoke(mux4smm, "/smm.api/invoke")

	rmux4smm := aclmux.NewAclServerMux("http", mux4smm, nil)
	httpServiceSMM := httpserver.NewHttpServer("httpPointSMM", rmux4smm)
	boot.AddService(httpServiceSMM)

	boot.Go(cfile)
}

func myInitor(
	gl *golua.GoLua,
	acclog *acclog.Service,
	mems *memserv.MemoryServ,
	sess *memserv4httpsession.Service,
	scs *servicecall.Service,
) {
	golua.InitCoreLibs(gl)
	vmmhttp.InitGoLuaWithHttpServ(gl)
	vmmhttp.InitGoLuaWithHttpClient(gl, acclog, "httpclient")
	vmmacclog.InitGoLua(gl)
	vmmjson.InitGoLua(gl)
	vmmsql.InitGoLua(gl)
	vmmclass.InitGoLua(gl)
	vmmesnp.InitGoLua(gl)
	vmmmemserv.InitGoLua(gl, mems)
	vmmservicecall.InitGoLua(gl, scs)
	http4goluaserv.InitGoLua(gl, sess)
}
