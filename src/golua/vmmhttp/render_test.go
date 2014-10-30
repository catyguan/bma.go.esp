package vmmhttp

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRenderPrepreprocess(t *testing.T) {
	fn := "test.vlua"
	str, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Error(err)
		return
	}
	r, err2 := RenderScriptPreprocess(string(str))
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Println("--- start --------------------------------")
	fmt.Println(r)
	fmt.Println("--- end ----------------------------------")
}
