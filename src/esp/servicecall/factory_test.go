package servicecall

import (
	"esp/espnet/mempipeline"
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

func TestHttpFactory(t *testing.T) {
	safeCall()

	cfg := make(map[string]interface{})
	cfg["Type"] = "http"
	cfg["URL"] = "http://127.0.0.1:1080/sample/servicecall.gl"

	fac := HttpServiceCallerFactory(0)
	sc, err := fac.Create("test", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	rv, err2 := sc.Call("say", []interface{}{"kitty"}, 0)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("Answer = %v\n", rv)
}

func T2estESNPMemp(t *testing.T) {
	safeCall()

	s := mempipeline.NewService()
	// sock := s.Open("test", "b")

	cfg := make(map[string]interface{})
	cfg["Type"] = "esnp.mem"
	cfg["Name"] = "test:a"
	cfg["TimeoutMS"] = 500

	fac := new(ESNPMemPipelineServiceCallerFactory)
	fac.S = s
	defer fac.S.Close()

	sc, err := fac.Create("test", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	rv, err2 := sc.Call("hello", []interface{}{"world"}, 0)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("%v\n", rv)
}
