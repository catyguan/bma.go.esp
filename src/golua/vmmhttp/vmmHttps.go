package vmmhttp

import (
	"bmautil/valutil"
	"esp/acclog"
	"fmt"
	"golua"
	"net"
	"net/http"
)

const tag = "vmmhttp"

func HttpServModule() *golua.VMModule {
	m := golua.NewVMModule("httpserv")
	m.Init("query", GOF_httpserv_query(0))
	m.Init("formValue", GOF_httpserv_formValue(0))
	m.Init("form", newQueryv("form"))
	m.Init("host", newQueryv("host"))
	m.Init("path", newQueryv("path"))
	m.Init("requestURI", newQueryv("requestURI"))
	m.Init("post", newQueryv("post"))
	m.Init("header", newQueryv("header"))
	m.Init("remoteAddr", newQueryv("remoteAddr"))
	m.Init("writeHeader", GOF_httpserv_writeHeader(0))
	m.Init("setHeader", GOF_httpserv_setHeader(0))
	m.Init("write", GOF_httpserv_write(0))
	return m
}

func doQuery(vm *golua.VM, n string) (interface{}, error) {
	ctx := vm.API_getContext()
	if ctx == nil {
		return nil, fmt.Errorf("doQuery(%s) fail - context nil", n)
	}
	req, _ := RequestFromContext(ctx)
	if req == nil {
		return nil, fmt.Errorf("doQuery(%s) fail - request nil", n)
	}
	switch n {
	case "host":
		return req.Host, nil
	case "path":
		return req.URL.Path, nil
	case "requestURI":
		return req.RequestURI, nil
	case "remoteAddr":
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		return ip, nil
	case "form":
		o := golua.NewVMTable(nil)
		err := req.ParseForm()
		if err != nil {
			return nil, err
		}
		for k, _ := range req.Form {
			v := req.FormValue(k)
			o.Rawset(k, v)
		}
		return o, nil
	case "post":
		o := golua.NewVMTable(nil)
		err := req.ParseForm()
		if err != nil {
			return nil, err
		}
		for k, _ := range req.PostForm {
			v := req.PostFormValue(k)
			o.Rawset(k, v)
		}
		return o, nil
	case "header":
		o := golua.NewVMTable(nil)
		for k, _ := range req.Header {
			v := req.Header.Get(k)
			o.Rawset(k, v)
		}
		return o, nil
	}
	return nil, fmt.Errorf("doQuery(%s) unknow action", n)
}

// httpserv.query(n:string[, defval])
type GOF_httpserv_query int

func (this GOF_httpserv_query) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	n, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")
	rv, err2 := doQuery(vm, vn)
	if err2 != nil {
		return 0, err2
	}
	if rv == nil {
		rv = dv
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_httpserv_query) IsNative() bool {
	return true
}

func (this GOF_httpserv_query) String() string {
	return "GoFunc<httpserv.query>"
}

// httpserv.formValue(n:string[, defval])
type GOF_httpserv_formValue int

func (this GOF_httpserv_formValue) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	n, dv, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vn := valutil.ToString(n, "")

	ctx := vm.API_getContext()
	if ctx == nil {
		return 0, fmt.Errorf("formValue(%v) fail - context nil", n)
	}
	req, _ := RequestFromContext(ctx)
	if req == nil {
		return 0, fmt.Errorf("formValue(%v) fail - request nil", n)
	}

	err2 := req.ParseForm()
	if err2 != nil {
		return 0, err2
	}

	rv := req.FormValue(vn)
	if rv == "" {
		rv = valutil.ToString(dv, "")
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_httpserv_formValue) IsNative() bool {
	return true
}

func (this GOF_httpserv_formValue) String() string {
	return "GoFunc<httpserv.formValue>"
}

// httpserv.queryv()
type GOF_httpserv_queryv struct {
	n string
}

func newQueryv(n string) *GOF_httpserv_queryv {
	return &GOF_httpserv_queryv{n}
}

func (this *GOF_httpserv_queryv) Exec(vm *golua.VM) (int, error) {
	vm.API_popAll()
	rv, err := doQuery(vm, this.n)
	if err != nil {
		return 0, err
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_httpserv_queryv) IsNative() bool {
	return true
}

func (this GOF_httpserv_queryv) String() string {
	return fmt.Sprintf("GoFunc<httpserv.%s>", this.n)
}

func respWriter(vm *golua.VM) (http.ResponseWriter, error) {
	ctx := vm.API_getContext()
	if ctx == nil {
		return nil, fmt.Errorf("respWriter fail - context nil")
	}
	w, _ := ResponseFromContext(ctx)
	if w == nil {
		return nil, fmt.Errorf("respWriter fail - responseWriter nil")
	}
	return w, nil
}

// httpserv.writeHeader(status:int)
type GOF_httpserv_writeHeader int

func (this GOF_httpserv_writeHeader) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	st, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vst := valutil.ToInt(st, http.StatusOK)

	w, err2 := respWriter(vm)
	if err2 != nil {
		return 0, err2
	}

	w.WriteHeader(vst)

	ctx := vm.API_getContext()
	if ctx != nil {
		adt, ok := acclog.AcclogDataFromContext(ctx)
		if ok {
			adt["status"] = vst
		}
	}

	return 0, nil
}

func (this GOF_httpserv_writeHeader) IsNative() bool {
	return true
}

func (this GOF_httpserv_writeHeader) String() string {
	return "GoFunc<httpserv.writeHeader>"
}

// httpserv.setHeader(n:string, val)
// httpserv.setHeader(hs:map)
type GOF_httpserv_setHeader int

func (this GOF_httpserv_setHeader) Exec(vm *golua.VM) (int, error) {
	w, errW := respWriter(vm)
	if errW != nil {
		vm.API_popAll()
		return 0, errW
	}

	top := vm.API_gettop()
	switch top {
	case 1:
		hs, err1 := vm.API_pop1X(-1, true)
		if err1 != nil {
			return 0, err1
		}
		vhs := vm.API_table(hs)
		if vhs == nil {
			return 0, fmt.Errorf("headers not a table(%T)", hs)
		}
		h := w.Header()
		for k, v := range vhs.ToMap() {
			vv := valutil.ToString(v, "")
			h.Set(k, vv)
		}
	case 2:
		n, v, err1 := vm.API_pop2X(-1, true)
		if err1 != nil {
			return 0, err1
		}
		vn := valutil.ToString(n, "")
		if vn == "" {
			return 0, fmt.Errorf("header name invalid(%v)", n)
		}
		vv := valutil.ToString(v, "")
		w.Header().Set(vn, vv)
	default:
		return 0, fmt.Errorf("invalid setHeader param count(%d)", top)
	}
	return 0, nil
}

func (this GOF_httpserv_setHeader) IsNative() bool {
	return true
}

func (this GOF_httpserv_setHeader) String() string {
	return "GoFunc<httpserv.setHeader>"
}

// httpserv.write(status:int)
type GOF_httpserv_write int

func (this GOF_httpserv_write) Exec(vm *golua.VM) (int, error) {
	err0 := vm.API_checkstack(1)
	if err0 != nil {
		return 0, err0
	}
	v, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}

	w, err2 := respWriter(vm)
	if err2 != nil {
		return 0, err2
	}

	switch rv := v.(type) {
	case []byte:
		w.Write(rv)
	default:
		vstr := valutil.ToString(rv, "")
		w.Write([]byte(vstr))
	}
	return 0, nil
}

func (this GOF_httpserv_write) IsNative() bool {
	return true
}

func (this GOF_httpserv_write) String() string {
	return "GoFunc<httpserv.write>"
}
