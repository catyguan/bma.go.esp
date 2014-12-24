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

	cfg := make(map[string]interface{})
	cfg["Type"] = "http"
	cfg["URL"] = "http://127.0.0.1:1080/sample/servicecall.gl"

	fac := HttpServiceCallerFactory(0)
	sc, err := fac.Create("test", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	ps := make(map[string]interface{})
	ps["world"] = "Kitty"
	rv, err2 := sc.Call("say", ps, 0)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("Answer = %v\n", rv)
}

func TestESNPNet(t *testing.T) {
	safeCall()

	defer func() {
		time.Sleep(100 * time.Millisecond)
	}()

	cfg := make(map[string]interface{})
	cfg["Type"] = "esnp"
	cfg["Address"] = "127.0.0.1:1080"
	cfg["TimeoutMS"] = 500

	fac := new(ESNPServiceCallerFactory)

	sc, err := fac.Create("serviceCall", cfg)
	if err != nil {
		t.Error(err)
		return
	}
	rv, err2 := sc.Call("hello", map[string]interface{}{"word": "world"}, 0)
	if err2 != nil {
		t.Error(err2)
		return
	}
	fmt.Printf("%v\n", rv)
}
