package cfprototype

import (
	"esp/espnet"
	"fmt"
	"net/url"
	"uprop"
)

const (
	tag = "CFP"
)

type Provider4ChannelFactoryPrototype func() ChannelFactoryPrototype

type ChannelFactoryStorage interface {
	GetStorageVersion() uint64
	GetChannelFactory(n string) (espnet.ChannelFactory, error)
}

type ChannelFactoryPrototype interface {
	GetProperties() []*uprop.UProperty

	Valid() error

	ToMap() map[string]interface{}

	FromMap(data map[string]interface{}) error

	ToURI() (*url.URL, error)

	FromURI(u *url.URL) error

	CreateChannelFactory(storage ChannelFactoryStorage, name string, start bool) (espnet.ChannelFactory, error)
}

var (
	regProviders map[string]Provider4ChannelFactoryPrototype
)

func init() {
	regProviders = make(map[string]Provider4ChannelFactoryPrototype)
	regProviders["dial"] = func() ChannelFactoryPrototype {
		return new(DialPoolPrototype)
	}
	regProviders["loadbalance"] = func() ChannelFactoryPrototype {
		return new(LoadBalancePrototype)
	}
}

func RegProvider4ChannelFactoryPrototype(n string, p Provider4ChannelFactoryPrototype) {
	regProviders[n] = p
}

func CreateChannelFactoryPrototype(kind string) ChannelFactoryPrototype {
	p, _ := regProviders[kind]
	if p != nil {
		return p()
	}
	return nil
}

func BuildChannelFactoryPrototype(s string) (ChannelFactoryPrototype, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	p, _ := regProviders[u.Scheme]
	if p == nil {
		return nil, fmt.Errorf("unknow scheme '%s'", u.Scheme)
	}
	cfp := p()
	err = cfp.FromURI(u)
	if err != nil {
		return nil, err
	}
	err = cfp.Valid()
	if err != nil {
		return nil, err
	}
	return cfp, nil
}

func ListChannelFactoryPrototype() []string {
	r := []string{}
	for k, _ := range regProviders {
		r = append(r, k)
	}
	return r
}
