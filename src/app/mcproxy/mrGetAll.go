package main

import "mcserver"

type mrGetResult struct {
	key    string
	result *mcserver.MemcacheResult
}

type mrGetAll struct {
	results []*mrGetResult
}

func (this *mrGetAll) Init() {
	this.results = make([]*mrGetResult, 0)
}

func (this *mrGetAll) HandleResponse(key string, r *mcserver.MemcacheResult) (more bool, done bool, err error) {
	o := new(mrGetResult)
	o.key = key
	o.result = r
	this.results = append(this.results, o)
	if r.Response == "END" {
		return false, true, nil
	}
	if iserr, _ := r.ToError(); iserr {
		return false, true, nil
	}
	return true, false, nil
}

func (this *mrGetAll) CheckEnd(okc, failc, errc, total int) (end bool, done bool, iserr bool) {
	if total == okc+failc {
		return true, true, false
	}
	return false, false, false
}
