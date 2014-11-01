package vmmhttp

import (
	"bmautil/httputil"
	"bmautil/valutil"
	"bytes"
	"esp/acclog"
	"fmt"
	"golua"
	"io"
	"mime"
	"net"
	"net/http"
	"path/filepath"
)

const tag = "vmmhttp"

func HttpServModule() *golua.VMModule {
	m := golua.NewVMModule("httpserv")
	m.Init("query", GOF_httpserv_query(0))
	m.Init("formValue", GOF_httpserv_formValue(0))
	m.Init("formFile", GOF_httpserv_formFile(0))
	m.Init("form", newQueryv("form"))
	m.Init("host", newQueryv("host"))
	m.Init("path", newQueryv("path"))
	m.Init("requestURI", newQueryv("requestURI"))
	m.Init("post", newQueryv("post"))
	m.Init("header", newQueryv("header"))
	m.Init("remoteAddr", newQueryv("remoteAddr"))
	m.Init("writeHeader", GOF_httpserv_writeHeader(0))
	m.Init("setHeader", GOF_httpserv_setHeader(0))
	m.Init("setContentType", GOF_httpserv_setContentType(0))
	m.Init("setContentFile", GOF_httpserv_setContentFile(0))
	m.Init("write", GOF_httpserv_write(0))
	m.Init("render", GOF_httpserv_render(0))
	m.Init("writeFile", GOF_httpserv_writeFile(0))
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
		// err := req.ParseForm()
		// if err != nil {
		// 	return nil, err
		// }
		for k, _ := range req.Form {
			v := req.FormValue(k)
			o.Rawset(k, v)
		}
		return o, nil
	case "post":
		o := golua.NewVMTable(nil)
		if httputil.IsMultipartForm(req) {
			fmt.Println(req.MultipartForm)
			for k, vs := range req.MultipartForm.Value {
				v := ""
				if len(vs) > 0 {
					v = vs[0]
				}
				o.Rawset(k, v)
			}
		} else {
			for k, _ := range req.PostForm {
				v := req.PostFormValue(k)
				o.Rawset(k, v)
			}
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

func (this GOF_httpserv_query) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

func (this GOF_httpserv_formValue) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

	// err2 := req.ParseForm()
	// if err2 != nil {
	// 	return 0, err2
	// }

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

func (this *GOF_httpserv_queryv) Exec(vm *golua.VM, self interface{}) (int, error) {
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

func (this GOF_httpserv_writeHeader) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

func (this GOF_httpserv_setHeader) Exec(vm *golua.VM, self interface{}) (int, error) {
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

// httpserv.write(content:string)
type GOF_httpserv_write int

func (this GOF_httpserv_write) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

	bs := golua.ToBytes(v)
	if bs != nil {
		w.Write(bs)
	} else {
		vstr := valutil.ToString(v, "")
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

// httpserv.setContentType(ct:string)
type GOF_httpserv_setContentType int

func (this GOF_httpserv_setContentType) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

	vstr := valutil.ToString(v, "")
	w.Header().Set("Content-Type", vstr)
	return 0, nil
}

func (this GOF_httpserv_setContentType) IsNative() bool {
	return true
}

func (this GOF_httpserv_setContentType) String() string {
	return "GoFunc<httpserv.setContentType>"
}

// httpserv.setContentFile(file:string)
type GOF_httpserv_setContentFile int

func (this GOF_httpserv_setContentFile) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
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

	vstr := valutil.ToString(v, "")
	ctype := mime.TypeByExtension(filepath.Ext(vstr))
	w.Header().Set("Content-Type", ctype)
	return 0, nil
}

func (this GOF_httpserv_setContentFile) IsNative() bool {
	return true
}

func (this GOF_httpserv_setContentFile) String() string {
	return "GoFunc<httpserv.setContentFile>"
}

// httpserv.formFile(n:string) : bytes:Object
type GOF_httpserv_formFile int

func (this GOF_httpserv_formFile) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	n, err1 := vm.API_pop1X(-1, true)
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

	if !httputil.IsMultipartForm(req) {
		return 0, nil
	}
	file, fh, err2 := req.FormFile(vn)
	if err2 == http.ErrMissingFile {
		return 0, nil
	}
	if err2 != nil {
		return 0, err2
	}
	fn := fh.Filename
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	sz, err3 := io.Copy(buf, file)
	if err3 != nil {
		return 0, err3
	}
	bo := golua.CreateGoBytes(buf.Bytes())

	vm.API_push(bo)
	vm.API_push(fn)
	vm.API_push(sz)

	return 3, nil
}

func (this GOF_httpserv_formFile) IsNative() bool {
	return true
}

func (this GOF_httpserv_formFile) String() string {
	return "GoFunc<httpserv.formFile>"
}

// httpserv.render(scriptName:string , viewData:map)
type GOF_httpserv_render int

func (this GOF_httpserv_render) Exec(vm *golua.VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(1)
	if err1 != nil {
		return 0, err1
	}
	sn, vo, err2 := vm.API_pop2X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vsn := valutil.ToString(sn, "")
	if vsn == "" {
		return 0, fmt.Errorf("invalid ScriptName(%v)", sn)
	}
	vtb := vm.API_toMap(vo)
	if vtb == nil {
		return 0, fmt.Errorf("invalid ViewData(%v)", vo)
	}

	gl := vm.GetGoLua()
	cc, err3 := gl.ChunkLoad(vm, vsn, true, RenderScriptPreprocess)
	if err3 != nil {
		return 0, err3
	}
	vm.API_push(cc)
	rc, err4 := vm.Call(0, 0, vtb)
	if err4 != nil {
		return rc, err4
	}
	return rc, nil
}

func (this GOF_httpserv_render) IsNative() bool {
	return true
}

func (this GOF_httpserv_render) String() string {
	return "GoFunc<httpserv.render>"
}

// httpserv.writeFile(fileName:string)
type GOF_httpserv_writeFile int

func (this GOF_httpserv_writeFile) Exec(vm *golua.VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(1)
	if err1 != nil {
		return 0, err1
	}
	n, err2 := vm.API_pop1X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vn := valutil.ToString(n, "")
	if vn == "" {
		return 0, fmt.Errorf("invalid FileName(%v)", n)
	}

	gl := vm.GetGoLua()
	bs, err3 := gl.FileLoad(vm, vn)
	if err3 != nil {
		return 0, err3
	}

	w, err4 := respWriter(vm)
	if err4 != nil {
		return 0, err4
	}

	if bs != nil {
		w.Write(bs)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
	return 0, nil
}

func (this GOF_httpserv_writeFile) IsNative() bool {
	return true
}

func (this GOF_httpserv_writeFile) String() string {
	return "GoFunc<httpserv.writeFile>"
}
