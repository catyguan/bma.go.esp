package main

import (
	"boot"
	"httpserver"
	"net/http"
	"os"
)

const (
	tag = "flServer"
)

func main() {
	cfile := "config/flServer-config.json"

	service := new(Service)
	service.name = "service"
	boot.AddService(service)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))
	mux.HandleFunc("/create", service.InvokeCreate)
	mux.HandleFunc("/", service.InvokeFL)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.AddService(httpService)

	boot.Go(cfile)
}
