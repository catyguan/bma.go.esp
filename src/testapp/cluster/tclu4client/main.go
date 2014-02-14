package main

import (
	"bmautil/socket"
	"boot"
	"config"
	"esp/espnet"
	"fmt"
	"logger"
	"time"
)

func main() {
	cfile := "../test/cluster4c-config.json"

	// espnet
	app := new(App)
	boot.QuickDefine(app, "", true)

	boot.Go(cfile)
}

const (
	tag = "CLU4C"
)

func socketInitor(s *socket.Socket) error {
	s.Trace = 128
	return nil
}

type appConfig struct {
	Remote    []string
	TimeoutMS int
}

func (this *appConfig) Valid() error {
	if len(this.Remote) == 0 {
		return fmt.Errorf("remote empty")
	}
	if this.TimeoutMS < 0 {
		this.TimeoutMS = 5000
	}
	return nil
}

type App struct {
	appConfig appConfig

	remotes []*espnet.DialPool
	pch     *espnet.PChannel
}

func (this *App) Name() string {
	return "application"
}

func (this *App) Start() bool {
	var cfg appConfig
	if !config.GetBeanConfig(this.Name(), &cfg) {
		logger.Error(tag, "no application config")
		return false
	}
	err := cfg.Valid()
	if err != nil {
		logger.Error(tag, "config invalid - %s", err)
		return false
	}
	this.remotes = make([]*espnet.DialPool, 0)
	for i, addr := range cfg.Remote {
		cfg := new(espnet.DialPoolConfig)
		cfg.Dial.Address = addr
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := espnet.NewDialPool(fmt.Sprintf("Remote%d", i), cfg, socketInitor)
		this.remotes = append(this.remotes, pool)
		boot.RuntimeStartRun(pool)
	}
	pch := espnet.NewPChannel("ch4cluc")
	for _, pool := range this.remotes {
		cf := pool.NewChannelFactory("espnet", time.Duration(cfg.TimeoutMS)*time.Millisecond)
		pch.Add(cf)
	}
	pch.Run()
	this.pch = pch

	this.appConfig = cfg
	return true
}

func (this *App) Run() bool {
	return true
}

func (this *App) Stop() bool {
	if this.pch != nil {
		this.pch.Stop()
	}
	for _, pool := range this.remotes {
		boot.RuntimeStopCloseClean(pool, true)
	}
	return true
}
