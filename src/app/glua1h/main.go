package main

import (
	"boot"
	"esp/acclog"
	"esp/glua"
	"esp/glua/httpmux4glua"
	"httpserver"
	"net/http"
	"os"
)

const (
	tag = "glua1h"
)

func main() {
	cfile := "config/glua4h-config.json"

	acclog := acclog.NewService("acclog")
	boot.AddService(acclog)

	service := glua.NewService("gluaService")
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))

	mux4gl := httpmux4glua.NewService("gluaMux", service)
	mux4gl.AccLog = acclog
	mux4gl.AccName = "http"
	mux4gl.InitMux(mux, "/i/")
	boot.AddService(mux4gl)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}
