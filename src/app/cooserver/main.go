package main

import (
	"bmautil/socket"
	"boot"
	"esp/aclserv"
	"esp/cluster/espnode"
	"esp/espnet/aclesnpmux"
	"esp/espnet/espservice"
	"smmapi/espmux4smmapi"
	_ "smmapi/smmapi4config"
	_ "smmapi/smmapi4server"

	_ "fileloader/http4fileloader"
	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-sqlite3"
)

const (
	tag = "xmemServer"
)

func main() {
	cfile := "config/cooserver-config.json"

	acls := aclserv.NewService("acl")
	boot.AddService(acls)

	// smm.api
	if true {
		smmapis := espmux4smmapi.NewService("smmapiServ")
		boot.AddService(smmapis)

		mux := espservice.NewServiceMux(nil, nil)
		smmapis.InitMuxInvoke(mux, "smm.api", "invoke")

		rmux := aclesnpmux.NewAclServerMux("esnp", mux.DoServe)

		goservice := espservice.NewGoService("serviceSMM", rmux.DoServe)

		lisPoint := socket.NewListenPoint("esnpPointSMM", nil, goservice.AcceptESP)
		boot.AddService(lisPoint)
	}

	// coo serv
	if true {
		// service := coo4espgroup.NewService("GroupCOO")
		// boot.AddService(service)
		nodeserv := espnode.NewService("espnode")
		boot.AddService(nodeserv)

		mux := espservice.NewServiceMux(nil, nil)
		// mux.AddServiceHandler("group.coo", service.DoServiceHandle)
		nodeserv.InitServiceMux(mux)

		rmux := aclesnpmux.NewAclServerMux("espnode", mux.DoServe)

		goservice := espservice.NewGoService("service", nodeserv.BindAuth(rmux.DoServe))

		lisPoint := socket.NewListenPoint("esnpPoint", nil, goservice.AcceptESP)
		boot.AddService(lisPoint)
	}

	boot.Go(cfile)
}
