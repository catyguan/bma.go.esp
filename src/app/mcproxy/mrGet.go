package main

import (
	"bytes"
	"errors"
	"mcserver"
	"strings"
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

func mergeResults(results []*mcserver.MemcacheResult) []byte {
	data := bytes.NewBuffer([]byte{})
	for _, r := range results {
		data.WriteString(r.Response)
		if len(r.Params) > 0 {
			data.WriteString(" ")
			data.WriteString(strings.Join(r.Params, " "))
		}
		data.WriteString("\r\n")
		if r.Data != nil {
			data.Write(r.Data)
			data.WriteString("\r\n")
		}
	}
	return data.Bytes()
}
