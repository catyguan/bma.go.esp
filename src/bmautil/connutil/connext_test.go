package connutil

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func testDebuger(b []byte, read bool) {
	fmt.Printf("%v : %X\n", read, b)
}

func T2estConnUtil(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "www.163.com:80", 3*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	ce := NewConnExt(conn, testDebuger)
	fmt.Println("connected", ce)
	defer ce.Close()

	_, err1 := ce.WriteString("GET /\n\r\n\r")
	if err1 != nil {
		t.Error(err1)
		return
	}
	isb := ce.CheckBreak()
	fmt.Println("IsBreak", isb)

	in := bufio.NewReader(ce)
	line, err2 := in.ReadString('\n')
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Println(line)
	time.Sleep(100 * time.Millisecond)
}

func TestConnCheckBreak(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:1080", 3*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	ce := NewConnExt(conn, testDebuger)
	fmt.Println("connected", ce)
	defer ce.Close()

	for i := 0; i < 8; i++ {
		isb := ce.CheckBreak()
		fmt.Println("IsBreak", isb)
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
}
