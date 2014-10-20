package vmmjson

import (
	"bmautil/valutil"
	"encoding/json"
	"golua"
)

const tag = "vmmjson"

func Module() *golua.VMModule {
	m := golua.NewVMModule("json")
	m.Init("encode", GOF_json_encode(0))
	m.Init("decode", GOF_json_decode(0))
	return m
}

// json.encode(data)
type GOF_json_encode int

func (this GOF_json_encode) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	data, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	data = golua.GoData(data)
	buf, err2 := json.Marshal(data)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(string(buf))
	return 1, nil
}

func (this GOF_json_encode) IsNative() bool {
	return true
}

func (this GOF_json_encode) String() string {
	return "GoFunc<json.encode>"
}

// json.decoe(str, godata)
type GOF_json_decode int

func (this GOF_json_decode) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	str, isgd, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	var rv interface{}
	vstr := valutil.ToString(str, "")
	if vstr != "" {
		rv = make(map[string]interface{})
		err2 := json.Unmarshal([]byte(vstr), &rv)
		if err2 != nil {
			return 0, err2
		}
		if !valutil.ToBool(isgd, false) {
			rv = golua.GoluaData(rv)
		}
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_json_decode) IsNative() bool {
	return true
}

func (this GOF_json_decode) String() string {
	return "GoFunc<json.decode>"
}
