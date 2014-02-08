package main

import "boot"

func main() {
	cfile := "../test/cluster4c-config.json"

	// espnet
	boot.Define(boot.INIT, "my", myInit)

	boot.Go(cfile)
}

const (
	tag = "CLU4C"
)

func myInit() bool {

	// service = espnet.NewGoService("SERVICE", smux.Serve)
	// boot.QuickDefine(service, "", true)

	// cfg := new(espnet.DialConfig)
	// cfg.Address = "127.0.0.1:1082"
	// dp := espnet.NewDialPoint("N0", cfg, socketAcceptSer)
	// dp.RetryFailInfoDruation = 30 * time.Second

	// boot.QuickDefine(dp, "N0", true)

	return true
}
