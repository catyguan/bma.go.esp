package http4fileloader

import (
	"fileloader"
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

func TestHttpLoader(t *testing.T) {
	safeCall()

	// url := "http://my.oschina.net/u/698121/blog/$F"
	// fn := "156245"
	url := "http://127.0.0.1:1085/query?m=test&f=$F&v=&c=09fd752adf1cf436a2fb132247af2f1f"
	fn := "hello.lua"

	cfg := make(map[string]interface{})
	cfg["Type"] = "http"
	cfg["URL"] = url

	fl, err := fileloader.CommonFileLoaderFactory.Create(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(bs))
}
