package main

import (
	"bmautil/valutil"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

var (
	msg chan string
)

func isError(s string) bool {
	return strings.Contains(s, "ERROR:")
}

func waitResp(conn net.Conn) bool {
	to := time.Duration(10) * time.Second
	totm := time.NewTimer(to)
	defer totm.Stop()

	select {
	case s := <-msg:
		fmt.Printf("<<<< %s\n", s)
		if isError(s) {
			return false
		}
	case <-totm.C:
		fmt.Printf("#ERROR: wait timeout")
		return false
	}

	totm.Reset(to)
	du := time.Duration(100) * time.Millisecond
	rtm := time.NewTimer(du)
	defer rtm.Stop()
	for {
		select {
		case s := <-msg:
			fmt.Printf("<<<< %s\n", s)
			if isError(s) {
				return false
			}
			rtm.Reset(du)
		case <-totm.C:
			fmt.Printf("#ERROR: wait timeout")
			return false
		case <-rtm.C:
			return true
		}
	}
}

func main() {
	msg = make(chan string, 5)
	defer close(msg)

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Println("runshell[.exe] [shellServiceAddress|port] file")
		return
	}
	addr := flag.Arg(0)
	fn := flag.Arg(1)

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("#ERROR: file '%s' load fail => %s\n", fn, err)
		return
	}
	clist := strings.Split(string(content), "\n")

	if !strings.Contains(addr, ":") {
		port := valutil.ToInt(addr, -1)
		if port > 0 {
			addr = fmt.Sprintf("127.0.0.1:%d", port)
		}
	}

	fmt.Printf("connecting to %s\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("#ERROR %s\n", err)
		return
	}
	fmt.Printf("connected %s\n", addr)
	defer func() {
		conn.Close()
		time.Sleep(10 * time.Millisecond)
	}()

	go func() {
		in := bufio.NewReader(conn)
		for {
			line, err := in.ReadString('\n')
			if err != nil {
				return
			}
			msg <- strings.TrimSpace(line)
		}
	}()

	for _, line := range clist {
		str := strings.TrimSpace(line)
		if str == "" {
			continue
		}
		if strings.HasPrefix(str, "#") {
			//注释
			fmt.Printf("-- %s\n", str)
			continue
		}
		fmt.Printf(">> %s\n", str)
		conn.Write([]byte(str))
		conn.Write([]byte{'\n'})
		if !waitResp(conn) {
			fmt.Printf("----- last: %s -----", str)
			return
		}
	}
	fmt.Printf("----- run '%s' @ '%s' end -----\n", fn, addr)
}
