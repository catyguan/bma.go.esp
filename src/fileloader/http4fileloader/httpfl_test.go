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

	cfg := make(map[string]interface{})
	cfg["Type"] = "http"
	cfg["URL"] = "http://my.oschina.net/u/698121/blog/"

	fl, err := fileloader.CommonFileLoaderFactory.Create(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	fn := "156245"
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(bs))
}
