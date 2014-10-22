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
	safeCall()

	cfg := make(map[string]interface{})
	cfg["Type"] = "file"
	dirs := make([]interface{}, 0)
	dirs = append(dirs, "./")
	cfg["Dirs"] = dirs

	fl, err := CommonFileLoaderFactory.Create(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	// fn := "interfaces.go"
	fn := "testdir:hello.txt"
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(bs))
}

func TestMFLoader(t *testing.T) {
	safeCall()

	ff := new(FileFileLoader)
	ff.Dirs = []string{"./testdir/$F"}
	SetModuleFileLoader("*", ff)

	cfg := make(map[string]interface{})
	cfg["Type"] = "m"

	fl, err := CommonFileLoaderFactory.Create(cfg)
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
}
