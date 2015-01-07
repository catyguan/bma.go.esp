package servicecall

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estHttpFactory(t *testing.T) {
	safeCall()
	InitBaseFactory()

	cfg := make(map[string]interface{})
	cfg["Type"] = "http"
	cfg["URL"] = "http://127.0.0.1:1080/sample/servicecall.gl"

	sc, err := DoCreate("test", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	ps := make(map[string]interface{})
	ps["world"] = "Kitty"
	rv, err2 := sc.Call("", "say", ps, time.Time{})
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("Answer = %v\n", rv)
}

func T2estESNPNet(t *testing.T) {
	safeCall()
	InitBaseFactory()

	defer func() {
		time.Sleep(100 * time.Millisecond)
	}()

	cfg := make(map[string]interface{})
	cfg["Type"] = "esnp"
	cfg["Address"] = "127.0.0.1:1080"
	cfg["TimeoutMS"] = 500

	sc, err := DoCreate("serviceCall", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	rv, err2 := sc.Call("", "hello", map[string]interface{}{"word": "world"}, time.Time{})
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("%v\n", rv)
}
