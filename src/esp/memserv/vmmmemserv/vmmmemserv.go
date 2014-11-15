package vmmmemserv

import (
	"esp/memserv"
	"golua"
)

func InitGoLua(gl *golua.GoLua, s *memserv.MemoryServ) {
	mf := new(MemServFactory)
	mf.s = s
	gl.SetObjectFactory("MemServ", mf.FactoryFunc)
}
