package main

import (
	"bmautil/valutil"
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
		fmt.Println("\tadd,sadd,ladd")
		fmt.Println("\treload,shutdown")
		fmt.Println("\tpprof.heap pprof.thread pprof.block pprof.cpu pprof.gor")
		fmt.Println("sample: esp4c.exe 127.0.0.1:1080 add")
		return
	}

	raddr := flag.Arg(0)
	mode := strings.ToLower(flag.Arg(1))
	switch mode {
	case "add":
		doAdd(raddr)
	case "ladd":
		sec := 5
		if flag.NArg() > 2 {
			sec = valutil.ToInt(flag.Arg(2), sec)
		}
		doLAdd(raddr, sec)
	case "sadd":
		doSAdd(raddr)
	case "reload":
		doReload(raddr)
	case "shutdown":
		doShutdow(raddr)
	default:
		if strings.HasPrefix(mode, "pprof.") {
			doPProf(raddr, mode[6:])
		} else {
			logger.Error(tag, "unknow mode '%s'", mode)
		}
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
