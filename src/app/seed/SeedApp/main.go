package main

import (
	"bmautil/socket"
	"boot"
	"esp/cluster/n2n"
	"esp/cluster/nodeinfo"
	"esp/espnet/espservice"
)

const (
	tag = "seed"
)

func main() {
	cfile := "config/seed-config.json"

	ninfo := nodeinfo.NewService("nodeInfo")
	boot.Add(ninfo, "", false)

	n2n := n2n.NewService("n2nService", ninfo)
	boot.Add(n2n, "", false)

	mux := espservice.NewServiceMux(nil, nil)
	mux.AddServiceHandler("n2n", n2n.Serve)
	goService := espservice.NewGoService("goService", mux.Serve)

	lisPoint := socket.NewListenPoint("servicePoint", nil, goService.AcceptESP)
	boot.Add(lisPoint, "", false)

	boot.Go(cfile)
}
