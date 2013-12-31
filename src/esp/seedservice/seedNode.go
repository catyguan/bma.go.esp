package seedservice

import (
	"esp/espnet"
	"logger"
)

type SeedNode struct {
	espnet.ExecuteNode
}

func NewSeedNode(name string) *SeedNode {
	this := new(SeedNode)
	this.Init(name, this.executeFunc, this.stopHandler)
	return this
}

func (this *SeedNode) executeFunc(ctx *espnet.EventContext) error {
	logger.Info(tag, "request %s, %v", ctx.Channel, ctx.Event)
	return nil
}

func (this *SeedNode) stopHandler() {

}
