package main

import (
	"esp/xmem/xmemservice"
	"logger"
)

type Tester struct {
	xmems *xmemservice.Service
}

func (this *Tester) Name() string {
	return "tester"
}

func (this *Tester) Run() bool {
	prof := new(xmemservice.MemGroupProfile)
	prof.Name = "test"
	prof.Coder = xmemservice.SimpleCoder(0)
	err := this.xmems.EnableMemGroup(prof)
	if err != nil {
		logger.Warn("test", "error - %s", err)
		return false
	}
	return true
}
