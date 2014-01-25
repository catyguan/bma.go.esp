package binlog

import (
	"fmt"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	cfg := new(BinlogConfig)
	cfg.FileName = "test.blog"

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	w, _ := bl.NewWriter()

	w.Write([]byte("hello"))
	w.Write([]byte("hello2"))

	time.Sleep(1 * time.Second)
	bl.Stop()
	bl.WaitStop()
}

func TestRead(t *testing.T) {
	cfg := new(BinlogConfig)
	cfg.Readonly = true
	cfg.FileName = "test.blog"

	bl := NewBinLog("test", 32, cfg)
	bl.Run()

	r, _ := bl.NewReader()

	fmt.Println(r.Seek(1389522970898990710))
	for {
		seq, bs, err := r.Read()
		if bs == nil {
			break
		}
		fmt.Println("READ", seq, string(bs), err)
	}
	time.Sleep(1 * time.Second)

	bl.Stop()
	bl.WaitStop()
}

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
