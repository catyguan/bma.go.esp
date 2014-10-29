package muxservice

import "net/http"

const (
	tag = "muxServer"
)

type HttpMuxBuilder func(mux *http.ServeMux)

type Service struct {
	name        string
	mux         *http.ServeMux
	muxBuilders []HttpMuxBuilder
}

func NewService(name string) *Service {
	this := new(Service)
	this.name = name
	this.muxBuilders = make([]HttpMuxBuilder, 0)
	return this
}

func (this *Service) Add(mb HttpMuxBuilder) {
	this.muxBuilders = append(this.muxBuilders, mb)
}

func (this *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := this.mux
	if m == nil {
		http.Error(w, "service invalid", http.StatusInternalServerError)
		return
	}
	m.ServeHTTP(w, req)
}
