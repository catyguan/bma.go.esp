package http4goluaserv

import (
	"esp/memserv/memserv4httpsession"
	"golua"
)

func InitGoLua(gl *golua.GoLua, s *memserv4httpsession.Service) {
	so := new(SessionObject)
	so.s = s
	gl.SetObjectFactory("HttpSession", so.FactoryFunc)
}
