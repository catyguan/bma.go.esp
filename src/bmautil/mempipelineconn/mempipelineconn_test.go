package mempipelineconn

import (
	"bmautil/netutil"
	"bufio"
	"fmt"
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

func TestMemPipeline(t *testing.T) {
	safeCall()
	defer func() {
		time.Sleep(100 * time.Millisecond)
	}()
	lis, err0 := netutil.Listen(NETWORK, "abc")
	if err0 != nil {
		t.Error(err0)
		return
	}
	defer lis.Close()
	go func() {
		for {
			c, err := lis.Accept()
			if c == nil {
				fmt.Println("listen end")
				return
			}
			if err != nil {
				fmt.Println("listen fail - ", err)
				return
			}
			fmt.Println("accepted - ", c)
			go func() {
				defer c.Close()
				b := make([]byte, 8)
				for {
					n, err := c.Read(b)
					if n > 0 {
						fmt.Println("read -- ", b[:n])
						c.Write(b[:n])
					}
					if err != nil {
						fmt.Println("read fail - ", err)
						return
					}
				}
			}()
		}
	}()

	conn, err := netutil.DialTimeout(NETWORK, "abc", 3*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	ce := conn
	fmt.Println("connected", ce)
	defer ce.Close()

	_, err1 := ce.Write([]byte("GET /\n\r\n\r"))
	if err1 != nil {
		t.Error(err1)
		return
	}
	fmt.Println("writed", ce)

	in := bufio.NewReader(ce)
	line, err2 := in.ReadString('\n')
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Println(line)

}
