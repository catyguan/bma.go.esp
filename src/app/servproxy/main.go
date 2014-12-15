package main

import (
	"boot"
	"esp/aclserv"
	"esp/goluaserv"
	"esp/servproxy"
	"fileloader"
	"golua"
	"golua/vmmclass"
	"golua/vmmjson"
	"golua/vmmsql"
	"httpserver"
	"httpserver/aclmux"
	"net/http"
	"smmapi/httpmux4smmapi"
)

const (
	tag = "servproxyApp"
)

func main() {
	cfile := "config/servproxy-config.json"

	acls := aclserv.NewService("acl")
	boot.AddService(acls)

	fl := fileloader.NewService("fileloader")
	boot.AddService(fl)

	goluaServ := goluaserv.NewService("goluaServ", func(gl *golua.GoLua) {
		myInitor(gl)
	})
	boot.AddService(goluaServ)

	service := servproxy.NewService("servproxy", goluaServ)
	boot.AddService(service)

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
) {
	golua.InitCoreLibs(gl)
	vmmjson.InitGoLua(gl)
	vmmsql.InitGoLua(gl)
	vmmclass.InitGoLua(gl)
}
