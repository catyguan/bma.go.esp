package httpserver

import (
	"boot"
	"fmt"
	"logger"
	"net/http"
	"strings"
)

const (
	tag2 = "httpFileServer"
)

type HttpFileServer struct {
	name   string
	config *HttpFileServerConfigInfo
}

func NewHttpFileServer(name string) *HttpFileServer {
	this := new(HttpFileServer)
	this.name = name
	return this
}

func (this *HttpFileServer) BuildMux(mux *http.ServeMux) {
	if this.config != nil {
		for _, l := range this.config.Location {
			p := l.Path
			d := l.Dir
			logger.Debug(tag2, "[%s] => [%s]", p, d)
			mux.Handle(p, http.StripPrefix(p, http.FileServer(http.Dir(d))))
		}
	}
}

type HttpFileServerConfigItem struct {
	Path string
	Dir  string
}

type HttpFileServerConfigInfo struct {
	Location []HttpFileServerConfigItem
}

func (this *HttpFileServerConfigInfo) Valid() error {
	chk := make(map[string]bool)
	for _, item := range this.Location {
		if !strings.HasSuffix(item.Path, "/") {
			item.Path = item.Path + "/"
		}
		if _, ok := chk[item.Path]; ok {
			return fmt.Errorf("path[%s] duplicate", item.Path)
		}
		if item.Dir == "" {
			return fmt.Errorf("path[%s] location empty", item.Path)
		}
	}
	return nil
}

func (this *HttpFileServerConfigInfo) Compare(old *HttpFileServerConfigInfo) int {
	return boot.CCR_CHANGE
}

func (this *HttpFileServer) Name() string {
	return this.name
}

func (this *HttpFileServer) Prepare() {
}
func (this *HttpFileServer) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(HttpFileServerConfigInfo)
	if !co.GetBeanConfig(this.name, cfg) {
		logger.Error(tag2, "'%s' miss config", this.name)
		return false
	}
	if err := cfg.Valid(); err != nil {
		logger.Error(tag2, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *HttpFileServer) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*HttpFileServerConfigInfo)
	this.config = cfg
	return true
}

func (this *HttpFileServer) Start(ctx *boot.BootContext) bool {
	return true
}

func (this *HttpFileServer) Run(ctx *boot.BootContext) bool {
	return true
}

func (this *HttpFileServer) GraceStop(ctx *boot.BootContext) bool {
	return true
}

func (this *HttpFileServer) Stop() bool {
	return true
}

func (this *HttpFileServer) Close() bool {
	return true
}

func (this *HttpFileServer) Cleanup() bool {
	return true
}
