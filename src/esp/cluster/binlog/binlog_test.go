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
	err := cfg.Valid()
	if err != nil {
		panic(err)
	}
	return cfg
}

func TestWrite(t *testing.T) {
	cfg := makeConfig()
	cfg.FileMaxSize = 130

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	w, _ := bl.NewWriter()

	lver, _ := w.GerVersion()
	fmt.Println(w.Write(lver+1, []byte("hello")))
	fmt.Println(w.Write(lver+2, []byte("hello2")))
	lver, _ = w.GerVersion()
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

	r, _ := bl.NewReader(clusterbase.OpVer(9))
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

	bl.Stop()
	bl.WaitStop()
}

/*
func TestMix(t *testing.T) {
	cfg := new(BinlogConfig)
	cfg.FileName = "test.blog"

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	go func() {
		w, _ := bl.NewWriter()
		for {
			time.Sleep(100 * time.Millisecond)
			if !w.Write([]byte(time.Now().String())) {
				fmt.Println("WRITE EXIT")
				return
			}
		}
	}()

	lis := func(seq BinlogVer, bs []byte, closed bool) {
		if closed {
			fmt.Println("LIS CLOSED")
		} else {
			fmt.Println("PUSH", seq, string(bs))
		}
	}
	r, _ := bl.NewReader()

	seq := BinlogVer(1390639648191602509)
	if true {
		r.SeekAndListen(seq, lis)
	} else {
		r.Seek(seq)
		for {
			seq, bs, err := r.Read()
			if bs == nil {
				break
			}
			fmt.Println("READ", seq, string(bs), err)
		}
		r.SetListener(lis)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("END")
	bl.Stop()
	bl.WaitStop()
}
*/
