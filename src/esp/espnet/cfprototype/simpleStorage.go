package cfprototype

import (
	"errors"
	"esp/espnet"
)

type SimpleChannelFactoryStorage struct {
	Factory map[string]espnet.ChannelFactory
}

func (this *SimpleChannelFactoryStorage) GetStorageVersion() uint64 {
	return 1
}

func (this *SimpleChannelFactoryStorage) GetChannelFactory(n string) (espnet.ChannelFactory, error) {
	r, ok := this.Factory[n]
	if ok {
		return r, nil
	}
	return nil, errors.New("not exists")
}
