package goyacc

import (
	"fmt"
	"io/ioutil"
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

func loadFile(f string) (string, error) {
	bs, err0 := ioutil.ReadFile(f)
	if err0 != nil {
		return "", err0
	}
	return string(bs), nil
}

func TestP1(t *testing.T) {
	safeCall()
	yyDebug = 2
	f := "test1.gom"
	// f := "test_go_syn.lua"
	content, err0 := loadFile(f)
	if err0 != nil {
		t.Error(err0)
		return
	}
	p := NewParser("test", content)
	node, err := p.Parse()
	if err != nil {
		fmt.Println(content)
		t.Error(err)
		return
	}
	fmt.Println("------------NODE---------------")
	fmt.Println(DumpNode("", node))
}
