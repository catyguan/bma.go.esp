package cfprototype

import (
	"esp/espnet"
	"uprop"
)

const (
	tag = "CFP"
)

type ChannelFactoryStorage interface {
	GetStorageVersion() uint64
	GetChannelFactory(n string) (espnet.ChannelFactory, error)
}

type ChannelFactoryPrototype interface {
	GetProperties() []*uprop.UProperty

	Valid() error

	ToMap() map[string]interface{}

	FromMap(data map[string]interface{}) error

	CreateChannelFactory(storage ChannelFactoryStorage, name string, start bool) (espnet.ChannelFactory, error)
}
