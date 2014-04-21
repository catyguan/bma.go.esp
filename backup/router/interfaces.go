package router

import "esp/espnet"

type IRouter interface {
	GetChannel(addr espnet.Address) espnet.Channel
}
