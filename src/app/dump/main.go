package main

import (
	"filelog"
	"flag"
	"fmt"
	"net"
	"time"
)

func proxy(id int, first bool, dstr bool, c1 net.Conn, c2 net.Conn, flog *filelog.FileLog) {
	defer c1.Close()
	defer c2.Close()

	ff := "<<"
	if first {
		ff = ">>"
	}

	for {
		b := make([]byte, 1024*1024)
		n, err := c1.Read(b)
		if err != nil && n == 0 {
			fmt.Println(c1.RemoteAddr().String(), "read fail", err)
			return
		}
		data := b[:n]

		var str string
		if dstr {
			str = string(data)
		} else {
			str = fmt.Sprintf("%X", data)
		}
		msg := fmt.Sprintf("[%d] %s %s %s:\n%s", id, c1.RemoteAddr().String(), ff, c2.RemoteAddr().String(), str)
		flog.Println(msg)
		c2.Write(data)
	}
}

func main() {

	dstr := false
	flag.BoolVar(&dstr, "s", false, "dump string")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("dump.exe localAddress remoteAddress [file]")
		return
	}
	laddr := flag.Arg(0)
	raddr := flag.Arg(1)
	fn := "dump"
	if flag.NArg() > 2 {
		fn = flag.Arg(2)
	}

	ln, err := net.Listen("tcp", laddr)
	fmt.Println("start at", ln.Addr().String())
	if err != nil {
		// handle error
		fmt.Println("listen fail", err)
		return
	}
	defer ln.Close()

	timeFormatString := "20060102_150405"
	filename := fmt.Sprintf("%s_%s.log", fn, time.Now().Format(timeFormatString))
	flog := filelog.NewFileLog(filename, 1024)
	err3 := flog.Open()
	if err3 != nil {
		fmt.Println("open file log fail", err3)
		return
	}

	cid := 1
	for {
		conn1, err := ln.Accept()
		if err != nil {
			return
		}

		conn2, err2 := net.Dial("tcp", raddr)
		if err2 != nil {
			fmt.Println("dial fail", err2)
			conn1.Close()
			return
		}
		flog.EnablePrint = true
		go proxy(cid, true, dstr, conn1, conn2, flog)
		go proxy(cid, false, dstr, conn2, conn1, flog)
		cid++
	}
}
