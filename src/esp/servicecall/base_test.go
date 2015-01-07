package servicecall

import (
	"fmt"
	"testing"
	"time"
)

func TestBase(t *testing.T) {
	safeCall()
	InitBaseFactory()

	cfg := make(map[string]interface{})
	cfg["Type"] = "esnp"
	cfg["Address"] = "127.0.0.1:1080"
	cfg["TimeoutMS"] = 500

	if true {
		sok, scid, err := SetServiceCall("serviceCall", cfg, nil, 0)
		fmt.Println(sok, scid, err)
	}
	if false {
		sok, scid, err := SetServiceCall(NAME_GATE_SERVICE, cfg, nil, 0)
		fmt.Println(sok, scid, err)
	}
	if false {
		SetLookup(func(serviceName string, deadline time.Time) (map[string]interface{}, error) {
			return cfg, nil
		})
	}
	if false {
		BindLookupService()
	}

	defer func() {
		time.Sleep(100 * time.Millisecond)
		RemoveAll()
	}()
	fmt.Println("test start")
	v, err := Call("serviceCall", "hello", map[string]interface{}{"word": "world"}, time.Now().Add(1*time.Second))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("result", v)
}
