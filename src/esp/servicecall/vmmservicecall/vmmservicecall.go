package vmmservicecall

import (
	"esp/servicecall"
	"golua"
)

const (
	tag = "vmmservicecall"
)

func InitGoLua(gl *golua.GoLua, s *servicecall.Service) {
	gl.SetObjectFactory("ServiceCall", ServiceCallGoLuaFactoryFunc(s))
}
