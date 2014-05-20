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
	boot.Add(service, "", false)

	mux := http.NewServeMux()
	mux4gl := httpmux4glua.NewService("gluaMux", service)
	mux4gl.InitMux(mux, "/")
	boot.AddService(mux4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}
