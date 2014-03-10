package main

import (
	"boot"
	"esp/cluster/nodegroup"
)

func CreateBootObject(node *nodegroup.NodeGroup) boot.BootObject {
	r := new(boot.BootWrap)
	r.SetStart(func(ctx *boot.BootContext) bool {
		return node.Start()
	})
	r.SetGraceStop(func(ctx *boot.BootContext) bool {
		return node.Stop()
	})
	r.SetStop(func() bool {
		return node.Stop()
	})
	r.SetCleanup(func() bool {
		node.WaitStop()
		return true
	})
	return r
}
