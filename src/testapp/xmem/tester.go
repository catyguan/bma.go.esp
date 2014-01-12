package main

import (
	"esp/xmem"
	"logger"
)

type Tester struct {
	xmems *xmem.Service
}

func (this *Tester) Name() string {
	return "tester"
}

func (this *Tester) Run() bool {
	prof := new(xmem.MemGroupProfile)
	prof.Name = "test"
	prof.Coder = xmem.SimpleCoder(0)
	err := this.xmems.EnableMemGroup(prof)
	if err != nil {
		logger.Warn("test", "error - %s", err)
		return false
	}
	return true
}
