package main

import (
	"esp/espnet/espsocket"
	"flag"
	"fmt"
	"logger"
	"net"
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

func createSocket(address string) espsocket.Socket {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		logger.Error(tag, "connect %s fail - %s", address, err)
		return nil
	}
	return espsocket.NewConnSocketN(conn, 0)
}
