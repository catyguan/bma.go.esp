package vmmesnp

import (
	"bmautil/valutil"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"golua"
	"logger"
	"time"
)

func NewGOOSocket(sock *espsocket.Socket) golua.VMTable {
	return golua.NewGOO(sock, gooSocket(0))
}

type gooSocket int

func (gooSocket) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*espsocket.Socket); ok {
		switch key {
		case "Call":
			return golua.NewGOF("ESNPSocket.Call", func(vm *golua.VM, self interface{}) (int, error) {
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				msg, tms, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vtms := valutil.ToInt(tms, 0)
				if vtms <= 0 {
					return 0, fmt.Errorf("TimeoutMS invalid")
				}
				m := vm.API_toMap(msg)
				if m == nil {
					return 0, fmt.Errorf("Message invalid")
				}
				addr, ok := m["Address"]
				if !ok {
					return 0, fmt.Errorf("Message.Address invalid")
				}
				maddr := vm.API_toMap(addr)
				if maddr == nil {
					return 0, fmt.Errorf("Message.Address not table")
				}
				var vhs map[string]interface{}
				hs, ok3 := m["Header"]
				if ok3 {
					vhs = vm.API_toMap(hs)
				}
				var vdt map[string]interface{}
				dt, ok5 := m["Data"]
				if ok5 {
					vdt = vm.API_toMap(dt)
				}

				emsg := esnp.NewRequestMessage()
				emsg.GetAddress().BindMap(maddr)
				ehs := emsg.Headers()
				for k, v := range vhs {
					ehs.Set(k, golua.BaseData(v))
				}
				edt := emsg.Datas()
				for k, v := range vdt {
					edt.Set(k, golua.BaseData(v))
				}
				rmsg, err2 := obj.Call(emsg, time.Duration(vtms)*time.Millisecond)
				if err2 != nil {
					return 0, err2
				}
				r := make(map[string]interface{})
				rhs := rmsg.Headers()
				rhsm, err3 := rhs.ToMap()
				if err3 != nil {
					return 0, err3
				}
				r["Header"] = rhsm
				rdt := rmsg.Datas()
				rdtm, err4 := rdt.ToMap()
				if err4 != nil {
					return 0, err4
				}
				r["Data"] = rdtm
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Close":
			return golua.NewGOF("ESNPSocket.Close", func(vm *golua.VM, self interface{}) (int, error) {
				vm.API_popAll()
				logger.Debug(tag, "socket(%s) close", obj)
				obj.AskClose()
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooSocket) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooSocket) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooSocket) CanClose() bool {
	return false
}

func (gooSocket) Close(o interface{}) {
}
