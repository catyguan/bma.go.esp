package gom

import (
	"boot"
	"fileloader"
	"testing"
)

func TestService(t *testing.T) {
	fl := fileloader.NewService("fileloader")
	boot.AddService(fl)

	service := NewService("gomServ", nil)
	boot.AddService(service)

	f := func() {
		fn := "none"
		// fn := "goyacc/test1.gom"
		sc := "dump"
		sc = "mysql:mysqltest.lua"
		ps := []string{}
		err := service.RunCommands(fn, sc, ps)
		if err != nil {
			t.Error(err)
		}
	}
	boot.TestGo("service_test.json", 3, []func(){f})
}
