package espnetdriver

import (
	"esp/espnet"
	"fmt"
	"net/url"
	"uprop"
)

const (
	tag = "CFP"
)

type ChannelDriver func() ChannelSource

type ChannelFactoryStorage interface {
	GetStorageVersion() uint64
	GetChannelFactory(n string) (espnet.ChannelFactory, error)
}

type ChannelSource interface {
	GetProperties() []*uprop.UProperty

	Valid() error

	ToMap() map[string]interface{}

	FromMap(data map[string]interface{}) error

	ToURI() (*url.URL, error)

	FromURI(u *url.URL) error

	CreateChannelFactory(storage ChannelFactoryStorage, name string, start bool) (espnet.ChannelFactory, error)
}

var (
	regDrivers map[string]ChannelDriver
)

func init() {
	regDrivers = make(map[string]ChannelDriver)
	regDrivers["dial"] = func() ChannelSource {
		return new(DialPoolSource)
	}
	regDrivers["loadbalance"] = func() ChannelSource {
		return new(LoadBalancePrototype)
	}
}

func RegChannelDriver(n string, p ChannelDriver) {
	regDrivers[n] = p
}

func ListChannelDriver() []string {
	r := []string{}
	for k, _ := range regDrivers {
		r = append(r, k)
	}
	return r
}

// Global ChannelSource Util
func CreateChannelSource(kind string) ChannelSource {
	p, _ := regDrivers[kind]
	if p != nil {
		return p()
	}
	return nil
}

func BuildChannelSource(s string) (ChannelSource, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	p, _ := regDrivers[u.Scheme]
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

// Global Channel Util
