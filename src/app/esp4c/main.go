package main

import (
	"bmautil/socket"
	"esp/espnet/esnp"
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
		fmt.Println("esp4c.exe remoteAddress mode[add]")
		return
	}

	raddr := flag.Arg(0)
	mode := strings.ToLower(flag.Arg(1))
	switch mode {
	case "add":
		doAdd(raddr)
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

func doAdd(address string) {
	c := createClient(address)
	if c == nil {
		return
	}
	defer c.Close()

	msg := esnp.NewMessage()
	msg.GetAddress().Set(esnp.ADDRESS_OP, "add")
	ds := msg.Datas()
	ds.Set("a", 1)
	ds.Set("b", 2)
	rmsg, err := c.Call(msg, nil)
	if err != nil {
		logger.Warn(tag, "call 'add' fail - %s", err)
		return
	}
	ds2 := rmsg.Datas()
	res, err2 := ds2.GetInt("c", 0)
	if err2 != nil {
		logger.Warn(tag, "result fail - %s", err2)
		return
	}
	logger.Info(tag, "result = %d", res)
}
