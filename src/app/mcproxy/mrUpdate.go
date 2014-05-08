package main

import "mcserver"

type mrUpdate struct {
	oks string
}

func (this *mrUpdate) HandleResponse(key string, r *mcserver.MemcacheResult) (more bool, done bool, err error) {
	res := new(remoteResult)
	res.key = key
	rdone := false
	if this.oks == "" || r.Response == this.oks {
		rdone = true
	}
	return false, rdone, nil
}

func (this *mrUpdate) CheckEnd(okc, failc, errc, total int) (end bool, done bool, iserr bool) {
	if failc > 0 {
		return true, false, errc > 0
	}
	if total == okc {
		return true, true, false
	}
	return false, false, false
}
