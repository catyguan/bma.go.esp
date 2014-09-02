package acclog

import (
	"boot"
	"fmt"
	"logger"
	"os"
	"strings"
)

type fileConfig struct {
	Name      string
	Path      string
	File      string
	QueueSize int
	MaxLines  int
	MaxSize   int
	NoDaily   bool
}

func (this *fileConfig) FilePrex() string {
	return this.Path + this.File
}

func (this *fileConfig) Valid() error {
	if this.Name == "" {
		return fmt.Errorf("accesslog Name empty")
	}
	if this.Path == "" {
		pw, _ := os.Getwd()
		this.Path = pw
	}
	if !(strings.HasSuffix(this.Path, "/") || strings.HasSuffix(this.Path, "\\")) {
		this.Path += "/"
	}
	if this.File == "" {
		if this.Name != "*" {
			this.File = this.Name
		} else {
			this.File = "acc"
		}
	}
	if this.QueueSize <= 0 {
		this.QueueSize = 128
	}
	if this.MaxSize < 0 {
		this.MaxSize = 0
	}

	return nil
}

func (this *fileConfig) Compare(old *fileConfig) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.Path != old.Path {
		return boot.CCR_NEED_START
	}
	if this.File != old.File {
		return boot.CCR_NEED_START
	}
	if this.QueueSize != old.QueueSize {
		return boot.CCR_NEED_START
	}
	if this.MaxLines != old.MaxLines {
		return boot.CCR_NEED_START
	}
	if this.MaxSize != old.MaxSize {
		return boot.CCR_NEED_START
	}
	if this.NoDaily != old.NoDaily {
		return boot.CCR_NEED_START
	}
	return boot.CCR_NONE
}

type configInfo struct {
	Nodes []*fileConfig
}

func (this *configInfo) Valid() error {
	m := make(map[string]bool)
	for _, fcfg := range this.Nodes {
		err := fcfg.Valid()
		if err != nil {
			return err
		}
		if ok := m[fcfg.Name]; ok {
			return fmt.Errorf("'%s' exists", fcfg.Name)
		}
		m[fcfg.Name] = true
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if len(this.Nodes) != len(old.Nodes) {
		return boot.CCR_NEED_START
	}
	m := make(map[string]*fileConfig)
	for _, fcfg := range this.Nodes {
		m[fcfg.Name] = fcfg
	}
	r := boot.CCR_NONE
	for _, of := range old.Nodes {
		if nf, ok := m[of.Name]; ok {
			v := nf.Compare(of)
			if v == boot.CCR_NEED_START {
				return v
			}
			if v != boot.CCR_NONE {
				r = v
			}
		} else {
			return boot.CCR_NEED_START
		}
	}
	return r
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
}
func (this *Service) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(configInfo)
	if !co.GetBeanConfig(this.name, cfg) {
		logger.Error(tag, "'%s' miss config", this.name)
		return false
	}
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := ccr.Config.(*configInfo)
	this.config = cfg
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	cfg := this.config
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, fc := range cfg.Nodes {
		if _, ok := this.nodes[fc.Name]; ok {
			continue
		}
		this._createNode(fc)
	}

	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}

	cfg := ccr.Config.(*configInfo)

	this.lock.Lock()
	defer this.lock.Unlock()

	for k, node := range this.nodes {
		del := false
		for _, fc := range cfg.Nodes {
			if fc.Name == k {
				v := fc.Compare(node.cfg)
				if v == boot.CCR_NEED_START {
					del = true
				} else {
					node.cfg = fc
				}
				break
			}
		}
		if del {
			this._closeNode(k, node)
		}
	}
	return true
}

func (this *Service) Stop() bool {
	return true
}

func (this *Service) Close() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, node := range this.nodes {
		this._closeNode(k, node)
	}
	return true
}

func (this *Service) Cleanup() bool {
	return true
}
