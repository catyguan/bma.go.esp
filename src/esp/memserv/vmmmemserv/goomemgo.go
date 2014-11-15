package vmmmemserv

import (
	"esp/memserv"
	"golua"
)

func NewGOOMemGo(mg *memserv.MemGo) golua.VMTable {
	return golua.NewGOO(mg, gooMemGo(0))
}

type gooMemGo int

func (gooMemGo) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*ObjectMemServ); !ok {
		switch key {
		case "Create":
			return golua.NewGOF("ESNP.Create", func(vm *golua.VM, self interface{}) (int, error) {
				if obj != nil {

				}
				return 0, nil
			}), nil
		case "Get":
			return golua.NewGOF("MemServ.Get", func(vm *golua.VM, self interface{}) (int, error) {
				return 0, nil
			}), nil
		}
	}
	return nil, nil
}

func (gooMemGo) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	return nil
}

func (gooMemGo) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	return r
}

func (gooMemGo) CanClose() bool {
	return false
}

func (gooMemGo) Close(o interface{}) {
}
