package main

import (
	"boot"
	"esp/cluster/n2n"
)

const (
	tag = "seed"
)

func main() {
	cfile := "config/seed-config.json"

	// mux := espservice.NewServiceMux(nil, nil)
	// mux.AddHandler("test", "add", H4Add)
	// mux.AddHandler("sys", "reload", H4Reload)

	n2nService := n2n.NewService("n2n")
	boot.Add(n2nService, "", false)

	// n2nPoint := socket.NewListenPoint("n2nPoint", nil, goservice.AcceptESP)
	// boot.Add(lisPoint, "", true)

	boot.Go(cfile)
}
