package binlog

import (
	"esp/cluster/clusterbase"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeConfig() *BinlogConfig {
	cfg := new(BinlogConfig)
	wd, _ := os.Getwd()
	cfg.LogDir = filepath.Join(wd, "testdir")
	cfg.FileFormatter = "test_%04d.blog"
	cfg.FileMaxSize = 1024

	err := cfg.Valid()
	if err != nil {
		panic(err)
	}
	return cfg
}

func TestWrite(t *testing.T) {
	cfg := makeConfig()

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	w, _ := bl.NewWriter()

	lver, _ := w.GetVersion()
	fmt.Println(w.Write(lver+1, []byte("hello")))
	fmt.Println(w.Write(lver+2, []byte("hello2")))
	lver, _ = w.GetVersion()
	fmt.Println("last version", lver)
	time.Sleep(1 * time.Second)
	bl.Stop()
	bl.WaitStop()
}

func TestRead(t *testing.T) {
	cfg := makeConfig()
	cfg.Readonly = true

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	r, _ := bl.NewReader()
	err0 := r.Seek(8)
	if err0 == nil {
		for {
			ok, seq, bs, err := r.Read()
			if !ok {
				break
			}
			fmt.Println("READ", seq, string(bs), err)
			if err != nil {
				break
			}
		}
		time.Sleep(1 * time.Second)
	} else {
		fmt.Println("seek fail", err0)
	}

	bl.Stop()
	bl.WaitStop()
}

func l4test(ver clusterbase.OpVer, bs []byte, closed bool) {
	if closed {
		fmt.Println("LISTENER CLOSED")
	} else {
		fmt.Println("PUSH", ver, string(bs))
	}
}

func TestFollow(t *testing.T) {
	cfg := makeConfig()
	cfg.Readonly = true

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	r, _ := bl.NewReader()
	r.Follow(4, l4test)

	bl.Stop()
	bl.WaitStop()
}

func TestMix(t *testing.T) {
	cfg := makeConfig()

	bl := NewBinLog("test", 32, cfg)
	bl.Run()
	w, _ := bl.NewWriter()

	go func() {
		wver, _ := w.GetVersion()
		for {
			wver++
			time.Sleep(100 * time.Millisecond)
			ok, _ := w.Write(wver, []byte(time.Now().String()))
			if !ok {
				fmt.Println("WRITE EXIT")
				return
			}
		}
	}()

	r, _ := bl.NewReader()
	r.Follow(0, l4test)

	time.Sleep(1 * time.Second)
	fmt.Println("END")
	lver, _ := w.GetVersion()
	fmt.Println("last version", lver)
	bl.Stop()
	bl.WaitStop()
}
