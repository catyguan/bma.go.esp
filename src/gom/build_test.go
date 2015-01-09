package gom

import (
	"bytes"
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

func TestB1(t *testing.T) {
	safeCall()
	fname := "goyacc/test1.gom"
	gm, err1 := Compile(fname)
	if err1 != nil {
		t.Error(err1)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	gm.Dump(buf, "")
	fmt.Println(buf.String())
}
