package fileloader

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estFFLoader(t *testing.T) {
	// tt := "app1/default/../../hello.txt"
	// fmt.Println(filepath.Clean("/" + tt))
	// if true {
	// 	return
	// }
	safeCall()

	cfg := make(map[string]interface{})
	cfg["Type"] = "file"
	dirs := make([]interface{}, 0)
	dirs = append(dirs, "./")
	cfg["Dirs"] = dirs

	fl, err := DoCreate(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	// fn := "interfaces.go"
	fn := "testdir:../../hello.txt"
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Println(string(bs))
}

func TestCFLoader(t *testing.T) {

	safeCall()

	cfg0 := make(map[string]interface{})
	cfg0["Type"] = "file"
	cfg0["Dirs"] = "./testdir/$!M$F"
	ff, err0 := DoCreate(cfg0)
	if err0 != nil {
		t.Error(err0)
		return
	}
	DefineFileLoader("fl1", ff)

	cfg := make(map[string]interface{})
	cfg["Type"] = "c"
	cfg["FL"] = "fl1"

	fl, err := DoCreate(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	fn := "test:hello.txt"
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Println(string(bs))
	fmt.Println("END")
}
