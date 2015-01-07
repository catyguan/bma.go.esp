package main

import (
	"boot"
	"esp/servicecall"
	"esp/servicecall/servicecallcfg"
	"fmt"
	"logger"
	"strings"
	"time"
)

const (
	tag = "esp4n"
)

func main() {
	cfile := "config/esp4n-config.json"

	scs := servicecallcfg.NewService("serviceCall")
	servicecallcfg.DefaultInit()
	boot.AddService(scs)

	bw := boot.NewBootWrap("main")
	bw.SetRun(func(ctx *boot.BootContext) bool {
		defer time.AfterFunc(1*time.Second, func() {
			boot.Shutdown()
		})

		if len(boot.Args) < 1 {
			fmt.Println("esp4n.exe mode")
			fmt.Println("\nhello")
			fmt.Println("sample: esp4n.exe hello")
			return true
		}

		mode := strings.ToLower(boot.Args[0])
		switch mode {
		case "hello":
			doHello()
		default:
			logger.Error(tag, "unknow mode '%s'", mode)
		}
		return true
	})
	boot.AddService(bw)

	boot.Go(cfile)
}

func doHello() {
	params := make(map[string]interface{})
	params["word"] = "kitty"
	rv, err1 := servicecall.CallTimeout("serviceCall", "hello", params, 1*time.Second)
	if err1 != nil {
		logger.Warn(tag, "call 'hello' fail - %s", err1)
		return
	}
	logger.Info(tag, "result = %v", rv)
}
