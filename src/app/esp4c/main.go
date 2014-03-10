package main

import (
	"bmautil/socket"
	"esp/espnet/espchannel"
	"esp/espnet/espclient"
	"flag"
	"fmt"
	"logger"
	"strings"
	"time"
)

const (
	tag = "esp4c"
)

func main() {

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("esp4c.exe remoteAddress mode")
		fmt.Println("\tadd,madd")
		fmt.Println("\treload")
		fmt.Println("sample: esp4c.exe 127.0.0.1:1080 add")
		return
	}

	raddr := flag.Arg(0)
	mode := strings.ToLower(flag.Arg(1))
	switch mode {
	case "add":
		doAdd(raddr)
	case "madd":
		doMAdd(raddr)
	case "reload":
		doReload(raddr)
	default:
		logger.Error(tag, "unknow mode '%s'", mode)
	}
	time.Sleep(1 * time.Second)
}

func createClient(address string) *espclient.ChannelClient {
	c := espclient.NewChannelClient()
	cfg := new(socket.DialConfig)
	cfg.Address = address
	err := c.Dial(tag, cfg, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
	if err != nil {
		logger.Error(tag, "connect %s fail - %s", address, err)
		return nil
	}
	return c
}
