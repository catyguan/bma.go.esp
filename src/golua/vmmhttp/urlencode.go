package vmmhttp

import (
	"bmautil/valutil"
	"golua"
	"net/url"
)

type GOF_http_urlencode int

func (this GOF_http_urlencode) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	r := url.QueryEscape(vs)
	vm.API_push(r)
	return 1, nil
}

func (this GOF_http_urlencode) IsNative() bool {
	return true
}

func (this GOF_http_urlencode) String() string {
	return "GoFunc<http.urlencode>"
}

type GOF_http_urldecode int

func (this GOF_http_urldecode) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	r, err2 := url.QueryUnescape(vs)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(r)
	return 1, nil
}

func (this GOF_http_urldecode) IsNative() bool {
	return true
}

func (this GOF_http_urldecode) String() string {
	return "GoFunc<http.urldecode>"
}
