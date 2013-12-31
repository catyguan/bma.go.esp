package tbus

import (
	"esp/espnet/cfprototype"
)

func CreateChannelFactoryPrototype(kind string) cfprototype.ChannelFactoryPrototype {
	switch kind {
	case "dial":
		r := new(cfprototype.DialPoolPrototype)
		return r
	case "loadbalance":
		r := new(cfprototype.LoadBalancePrototype)
		return r
	}
	return nil
}

func ListChannelFactoryPrototype() []string {
	return []string{"dial", "loadbalance"}
}
