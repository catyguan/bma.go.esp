package http4goluaserv

import (
	"acl"
	"bmautil/valutil"
	"esp/memserv/memserv4httpsession"
	"fmt"
	"golua"
	"golua/vmmhttp"
)

type SessionObject struct {
	s *memserv4httpsession.Service
}

func (this *SessionObject) FactoryFunc(vm *golua.VM, n string) (interface{}, error) {
	return golua.NewGOO(this, gooSession(0)), nil
}

type gooSession int

func (gooSession) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*SessionObject); ok {
		switch key {
		case "Close":
			return golua.NewGOF("Session.Close", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				err2 := obj.s.CloseSession(req)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "GetUser":
			return golua.NewGOF("Session.GetUser", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				user, err2 := obj.s.GetUser(req)
				if err2 != nil {
					return 0, err2
				}
				if user == nil {
					vm.API_push(nil)
				} else {
					r := make(map[string]interface{})
					r["Account"] = user.Account
					r["Domain"] = user.Domain
					if user.Groups != nil {
						gs := make([]interface{}, len(user.Groups))
						for i, g := range user.Groups {
							gs[i] = g
						}
						r["Groups"] = gs
					}
					vm.API_push(r)
				}
				return 1, nil
			}), nil
		case "SetUser":
			return golua.NewGOF("Session.SetUser", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				w := vmmhttp.ResponseFromVM(vm)
				if w == nil {
					return 0, fmt.Errorf("http response invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				user, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				if user == nil {
					err2 := obj.s.SetUser(w, req, nil)
					if err2 != nil {
						return 0, err2
					}
					return 0, nil
				}
				m := vm.API_toMap(user)
				if m == nil {
					return 0, fmt.Errorf("invalid User table")
				}
				o := new(acl.User)
				o.Account = valutil.ToString(m["Account"], "")
				o.Domain = valutil.ToString(m["Domain"], "")
				gs := m["Groups"]
				if gs != nil {
					if arr, ok := gs.([]interface{}); ok {
						o.Groups = make([]string, len(arr))
						for i, g := range arr {
							o.Groups[i] = valutil.ToString(g, "")
						}
					}
				}
				err2 := obj.s.SetUser(w, req, o)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Get":
			return golua.NewGOF("Session.Get", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				r, err2 := obj.s.GetSession(req, vn)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(r)
				return 1, nil
			}), nil
		case "MGet":
			return golua.NewGOF("Session.Get", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				arr := vm.API_toSlice(v)
				if arr == nil {
					return 0, fmt.Errorf("keys invalid")
				}
				ks := make([]string, len(arr))
				for i, kv := range arr {
					ks[i] = valutil.ToString(kv, "")
				}
				r, err2 := obj.s.MGetSession(req, ks)
				if err2 != nil {
					return 0, err2
				}
				vm.API_push(r)
				return 1, nil
			}), nil
		case "Set":
			return golua.NewGOF("Session.Set", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				w := vmmhttp.ResponseFromVM(vm)
				if w == nil {
					return 0, fmt.Errorf("http response invalid")
				}
				err0 := vm.API_checkStack(2)
				if err0 != nil {
					return 0, err0
				}
				n, v, err1 := vm.API_pop2X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				err2 := obj.s.SetSession(w, req, vn, v)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "MSet":
			return golua.NewGOF("Session.MSet", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				w := vmmhttp.ResponseFromVM(vm)
				if w == nil {
					return 0, fmt.Errorf("http response invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				m, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				mo := vm.API_toMap(m)
				if mo == nil {
					return 0, fmt.Errorf("MSet table invalid")
				}
				err2 := obj.s.MSetSession(w, req, mo)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "Delete":
			return golua.NewGOF("Session.Delete", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				n, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				vn := valutil.ToString(n, "")
				if vn == "" {
					return 0, fmt.Errorf("Name invalid")
				}
				err2 := obj.s.DeleteSession(req, vn)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		case "MDelete":
			return golua.NewGOF("Session.MDelete", func(vm *golua.VM, self interface{}) (int, error) {
				req := vmmhttp.RequestFromVM(vm)
				if req == nil {
					return 0, fmt.Errorf("http request invalid")
				}
				err0 := vm.API_checkStack(1)
				if err0 != nil {
					return 0, err0
				}
				v, err1 := vm.API_pop1X(-1, true)
				if err1 != nil {
					return 0, err1
				}
				arr := vm.API_toSlice(v)
				if arr == nil {
					return 0, fmt.Errorf("keys invalid")
				}
				ks := make([]string, len(arr))
				for i, kv := range arr {
					ks[i] = valutil.ToString(kv, "")
				}
				err2 := obj.s.MDeleteSession(req, ks)
				if err2 != nil {
					return 0, err2
				}
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooSession) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooSession) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooSession) CanClose() bool {
	return false
}

func (gooSession) Close(o interface{}) {
}
