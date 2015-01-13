package main

import (
	"boot"
	"fileloader"
	"fmt"
	"golua"
	"golua/vmmclass"
	"golua/vmmjson"
	"gom"
	"logger"
	"time"
)

const (
	tag = "gom"
)

func main() {
	cfile := "config/gom-config.json"

	fl := fileloader.NewService("fileloader")
	boot.AddService(fl)

	service := gom.NewService("gomServ", func(gl *golua.GoLua) {
		myInitor(gl)
	})
	boot.AddService(service)

	bw := boot.NewBootWrap("main")
	bw.SetRun(func(ctx *boot.BootContext) bool {
		go doRun(service)
		return true
	})
	boot.AddService(bw)

	boot.Go(cfile)
}

func myInitor(
	gl *golua.GoLua,
) {
	golua.InitCoreLibs(gl)
	vmmjson.InitGoLua(gl)
	vmmclass.InitGoLua(gl)
}

func doRun(s *gom.Service) {
	defer func() {
		time.Sleep(100 * time.Millisecond)
		boot.Shutdown()
	}()

	if len(boot.Args) < 1 {
		fmt.Println(">> gom gomFile [gomScript [param...]]")
		return
	}
	gfname := boot.Args[0]
	gscript := ""
	var ps []string
	if len(boot.Args) > 1 {
		gscript = boot.Args[1]
		ps = boot.Args[2:]
	}
	fmt.Printf(">> run %s, %s, %s\n", gfname, gscript, ps)
	err := s.RunCommands(gfname, gscript, ps)
	if err != nil {
		logger.Error(tag, "RunCommands - %s,%s fail\n - %s", gfname, gscript, err)
	}
}
