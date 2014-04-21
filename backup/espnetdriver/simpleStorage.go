package espnetdriver

import (
	"boot"
	"errors"
	"esp/espnet"
)

type SimpleChannelFactoryStorage struct {
	version uint64
	Factory map[string]espnet.ChannelFactory
}

func (this *SimpleChannelFactoryStorage) GetStorageVersion() uint64 {
	if this.version == 0 {
		this.version = 1
	}
	return this.version
}

func (this *SimpleChannelFactoryStorage) GetChannelFactory(n string) (espnet.ChannelFactory, error) {
	if this.Factory != nil {
		r, ok := this.Factory[n]
		if ok {
			return r, nil
		}
	}
	return nil, errors.New("not exists")
}

func (this *SimpleChannelFactoryStorage) SetChannelFactory(n string, cf espnet.ChannelFactory) error {
	if this.Factory == nil {
		this.Factory = make(map[string]espnet.ChannelFactory)
	}
	this.Factory[n] = cf
	this.version = this.version + 1
	return nil
}

func (this *SimpleChannelFactoryStorage) Close(wait bool) {
	for _, o := range this.Factory {
		boot.RuntimeStopCloseClean(o, wait)
	}
	this.Factory = nil
}
