package main

import (
	"boot"
	"esp/glua"
	"esp/glua/httpmux4glua"
	"httpserver"
	"net/http"
)

const (
	tag = "glua1h"
)

func main() {
	cfile := "config/glua1h-config.json"

	service := glua.NewService("gluaService")
	boot.Add(service, "", true)

	mux := http.NewServeMux()
	httpmux4glua.InitMux(mux, "/", service)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.Add(httpService, "", true)

	boot.Go(cfile)
}
