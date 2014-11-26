package espnetss

import (
	"fmt"
	"logger"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estMakeSplit(t *testing.T) {
	host := "127.0.0.1:1080"
	user := "test"
	lt := "base"
	cert := "123456"
	key := Make(host, user, lt, cert)
	fmt.Println("Make =>", key)
	s1, s2, s3, s4 := Split(key)
	fmt.Printf("Split =>H:%s, U:%s, T:%s, C:%s\n", s1, s2, s3, s4)
}

func TestSocketSource(t *testing.T) {
	safeCall()

	cfg := new(Config)
	cfg.Host = "172.19.16.97:80"
	cfg.User = "test"
	cfg.PoolSize = 2
	cfg.PreConns = 1
	ss := NewSocketSource(cfg)
	defer func() {
		ss.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	ss.Start()
	time.Sleep(100 * time.Millisecond)
	sock, err := ss.Open(100)
	if err != nil {
		t.Error(err)
		return
	}
	logger.Debug("test", "open & close")
	// sock.AskClose()
	ss.Return(sock)

	time.Sleep(100 * time.Millisecond)
	logger.Debug("test", "end")
}

func T2estConnFail(t *testing.T) {
	cfg := new(Config)
	cfg.Host = "127.0.0.1:1080"
	cfg.User = "test"
	cfg.PoolSize = 1
	cfg.PreConns = 1
	ss := NewSocketSource(cfg)
	defer func() {
		ss.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	ss.Start()

	time.Sleep(5 * time.Second)
}
