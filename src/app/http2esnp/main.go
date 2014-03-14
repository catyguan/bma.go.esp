package main

import (
	"boot"
	"httpserver"
	"net/http"
	"os"
)

const (
	tag = "http2esnp"
)

func main() {
	cfile := "config/http2esnp-config.json"

	service := new(Service)
	service.name = "service"
	boot.Add(service, "", true)

	var wd, _ = os.Getwd()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+"/public"))))
	mux.HandleFunc("/i/", service.InvokeESNP)

	httpService := httpserver.NewHttpServer("httpPoint", mux)
	boot.Add(httpService, "", true)

	boot.Go(cfile)
}
