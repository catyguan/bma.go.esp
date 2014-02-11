package main

import (
	"flag"
	"fmt"
	"net"
)

func readAll(id int, c1 net.Conn) {
	defer func() {
		c1.Close()
		fmt.Println("OnClosed", id, c1.RemoteAddr())
	}()

	fmt.Println("OnAccept", id, c1.RemoteAddr())
	for {
		b := make([]byte, 1024*1024)
		n, err := c1.Read(b)
		if err != nil && n == 0 {
			fmt.Println("OnError", id, c1.RemoteAddr(), "read fail", err)
			return
		}
		fmt.Println("OnRead", id, c1.RemoteAddr(), ">>", n)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("bhole.exe localAddress")
		return
	}
	laddr := flag.Arg(0)

	ln, err := net.Listen("tcp", laddr)
	fmt.Println("start at", ln.Addr().String())
	if err != nil {
		// handle error
		fmt.Println("listen fail", err)
		return
	}
	defer ln.Close()

	cid := 1
	for {
		conn1, err := ln.Accept()
		if err != nil {
			return
		}
		go readAll(cid, conn1)
		cid++
	}
}
