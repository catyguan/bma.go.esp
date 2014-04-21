package router

import (
	"bmautil/valutil"
	"errors"
	"esp/espnet"
	"time"
)

type ChannelFactoryProvider func(name string, cfg map[string]interface{}) (espnet.ChannelFactory, error)

var channelFactoryProviders map[string]ChannelFactoryProvider = make(map[string]ChannelFactoryProvider)

func RegChannelFactoryProvider(typ string, p ChannelFactoryProvider) {
	channelFactoryProviders[typ] = p
}

func NewChannelFactory(typ string, name string, cfg map[string]interface{}) (espnet.ChannelFactory, error) {
	p, ok := channelFactoryProviders[typ]
	if !ok {
		return nil, errors.New("unknow provider '" + typ + "'")
	}
	return p(name, cfg)
}

func init() {
	RegChannelFactoryProvider("DialPool", providerDialPool)
	RegChannelFactoryProvider("Dial", providerDial)
}

// DialPool
type myDialPool struct {
	service *espnet.DialPool
	factory espnet.ChannelFactory
}

func (this *myDialPool) Start() bool {
	return this.service.Start()
}

func (this *myDialPool) Run() bool {
	return this.service.Run()
}

func (this *myDialPool) Close() bool {
	return this.service.Close()
}

func (this *myDialPool) NewChannel() (espnet.Channel, error) {
	return this.factory.NewChannel()
}

func providerDialPool(name string, cfg map[string]interface{}) (espnet.ChannelFactory, error) {
	var dc espnet.DialConfig
	if !valutil.ToBean(cfg, &dc) {
		return nil, errors.New("miss DialConfig")
	}
	if err := dc.Valid(name); err != nil {
		return nil, err
	}
	max := valutil.ToInt(cfg["MaxSize"], 0)
	initSize := valutil.ToInt(cfg["InitSize"], 0)
	tm := valutil.ToInt(cfg["Timeout"], 0)
	chcoder := valutil.ToString(cfg["Coder"], "")

	pool := espnet.NewDialPool(name, &dc, max, nil)
	pool.InitSize = initSize
	f := pool.NewChannelFactory(chcoder, time.Duration(tm)*time.Millisecond)

	r := new(myDialPool)
	r.service = pool
	r.factory = f

	return r, nil
}

// Dial
func providerDial(name string, cfg map[string]interface{}) (espnet.ChannelFactory, error) {
	var dc espnet.DialConfig
	if !valutil.ToBean(cfg, &dc) {
		return nil, errors.New("miss DialConfig")
	}
	if err := dc.Valid(name); err != nil {
		return nil, err
	}
	chcoder := valutil.ToString(cfg["Coder"], "")

	r := espnet.NewSimpleDialPool(name, &dc, nil, chcoder)
	return r, nil
}
