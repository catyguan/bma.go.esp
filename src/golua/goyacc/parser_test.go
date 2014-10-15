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
	bs, err0 := ioutil.ReadFile("../samplecodes/" + f)
	if err0 != nil {
		return "", err0
	}
	return string(bs), nil
}

func TestExecOp(t *testing.T) {
	safeCall()

	var v1, v2 interface{}

	op := OP_EQ
	v1 = 1
	v2 = nil
	ok, r, err := ExecOp2(op, v1, v2)
	fmt.Println("result", ok, r, err)
}

func T2estP1(t *testing.T) {
	safeCall()
	yyDebug = 0
	f := "test1.lua"
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
