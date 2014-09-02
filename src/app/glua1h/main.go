package main

import (
	"boot"
	"esp/acclog"
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

	acclog := acclog.NewService("acclog")
	boot.Add(acclog, "", false)

	service := glua.NewService("gluaService")
	boot.Add(service, "", false)

	mux := http.NewServeMux()
	mux4gl := httpmux4glua.NewService("gluaMux", service)
	mux4gl.AccLog = acclog
	mux4gl.AccName = "http"
	mux4gl.InitMux(mux, "/")
	boot.Add(mux4gl, "", false)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.Add(httpService, "", false)

	boot.Go(cfile)
}
