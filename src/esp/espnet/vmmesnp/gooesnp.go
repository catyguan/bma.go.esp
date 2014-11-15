package vmmesnp

import (
	"bmautil/valutil"
	"esp/espnet/espsocket"
	"fmt"
	"golua"
)

func ESNPFactory(vm *golua.VM, n string) (interface{}, error) {
	return golua.NewGOO(0, gooESNP(0)), nil
}

func API_open(vm *golua.VM, host, token string, tms int) (*espsocket.Socket, error) {
	gos, errC := vm.GetGoLua().SingletonService("ESNP", createESNP)
	if errC != nil {
		return nil, errC
	}
	serv, ok := gos.(*esnpserv)
	if !ok {
		return nil, fmt.Errorf("invalid ESNP")
	}
	sock, err := serv.Open(host, token, tms)
	if err != nil {
		return nil, err
	}
	return sock, nil
}

type gooESNP int

func (gooESNP) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	gos, errC := vm.GetGoLua().SingletonService("ESNP", createESNP)
	if errC != nil {
		return nil, errC
	}
	if obj, ok := gos.(*esnpserv); ok {
		switch key {
		case "Create":
			return golua.NewGOF("ESNP.Create", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				host, tms, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vhost := valutil.ToString(host, "")
				if vhost == "" {
					return 0, fmt.Errorf("Host invalid")
				}
				vtms := valutil.ToInt(tms, 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				sock, err2 := obj.Create(vhost, vtms)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(NewGOOSocket(sock))
				return 1, nil
			}), nil
		case "Open":
			return golua.NewGOF("ESNP.Open", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(3)
				if err0 != nil {
					return 0, err0
				}
				host, token, tms, err1 := vm.API_pop3X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vhost := valutil.ToString(host, "")
				if vhost == "" {
					return 0, fmt.Errorf("Host invalid")
				}
				vtoken := valutil.ToString(token, "")
				vtms := valutil.ToInt(tms, 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				sock, err2 := obj.Open(vhost, vtoken, vtms)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(NewGOOSocket(sock))
				return 1, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooESNP) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooESNP) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooESNP) CanClose() bool {
	return false
}

func (gooESNP) Close(o interface{}) {
}
