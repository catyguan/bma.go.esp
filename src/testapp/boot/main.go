package main

import (
	"boot"
	"fmt"
	"os"
)

type obk struct {
}

func (this *obk) Init() bool {
	fmt.Println("my init")
	return true
}

func (this *obk) Stop() bool {
	fmt.Println("on my stop")
	return true
}

func (this *obk) Init2() bool {
	fmt.Println("my init2")
	return true
}

func main() {
	var wd, _ = os.Getwd()
	cfile := wd + "/../test/esp-config.json"
	// cfile = ""

	ob := obk{}
	boot.Define(boot.INIT, "", ob.Init)
	boot.DefineOrder(boot.INIT, "i2", ob.Init2, -10)
	boot.DefineAfter("i2", boot.INIT, "i3", func() bool {
		fmt.Println("my init3")
		return true
	})
	boot.Define(boot.RUN, "", func() bool {
		boot.Shutdown()
		// boot.Shutdown()
		// boot.Shutdown()
		return true
	})
	boot.Define(boot.STOP, "", ob.Stop)
	boot.Define(boot.STOP, "", func() bool {
		fmt.Println("stop 2")
		return true
	})

	boot.Go(cfile)
}
