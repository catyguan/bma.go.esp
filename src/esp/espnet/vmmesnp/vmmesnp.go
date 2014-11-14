package vmmesnp

import "golua"

func InitGoLua(gl *golua.GoLua) {
	gl.SetObjectFactory("ESNP", ESNPFactory)
}
