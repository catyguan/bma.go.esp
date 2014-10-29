package muxservice

import (
	"boot"
	"logger"
	"net/http"
)

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
}
func (this *Service) CheckConfig(ctx *boot.BootContext) bool {
	ccr := boot.NewConfigCheckResult(boot.CCR_CHANGE, nil)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	logger.Debug(tag, "'%s' build http mux", this.name)
	bl := this.muxBuilders
	mux := http.NewServeMux()
	for _, b := range bl {
		b(mux)
	}
	this.mux = mux
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
