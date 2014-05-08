package main

import (
	"errors"
	"mcserver"
)

type mrGet struct {
	results []*mcserver.MemcacheResult
}

func (this *mrGet) Init() {
	this.results = make([]*mcserver.MemcacheResult, 0)
}

func (this *mrGet) HandleResponse(key string, r *mcserver.MemcacheResult) (more bool, done bool, err error) {
	isErr, errMsg := r.ToError()
	if isErr {
		return false, false, errors.New(errMsg)
	}
	this.results = append(this.results, r)
	if r.Response == "END" {
		return false, true, nil
	}
	return true, false, nil
}

func (this *mrGet) CheckEnd(okc, failc, errc, total int) (end bool, done bool, iserr bool) {
	if okc > 0 {
		return true, true, false
	}
	return false, false, false
}
