package gom

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"golua"
	"os"
	"path/filepath"
)

func InitGoLua(gl *golua.GoLua) {
	Module().Bind(gl)
}

func Module() *golua.VMModule {
	m := golua.NewVMModule("gomf")
	m.Init("render", GOF_gomf_render(0))
	return m
}

// httpserv.render(scriptName:string , viewData:map, fname:string)
type GOF_gomf_render int

func (this GOF_gomf_render) Exec(vm *golua.VM, self interface{}) (int, error) {
	err1 := vm.API_checkStack(1)
	if err1 != nil {
		return 0, err1
	}
	sn, vo, fn, err2 := vm.API_pop3X(-1, true)
	if err2 != nil {
		return 0, err2
	}
	vsn := valutil.ToString(sn, "")
	if vsn == "" {
		return 0, fmt.Errorf("invalid ScriptName(%v)", sn)
	}
	vfn := valutil.ToString(fn, "")
	if vfn == "" {
		return 0, fmt.Errorf("invalid FileName(%v)", fn)
	}
	vtb := vm.API_toMap(vo)
	if vtb == nil {
		return 0, fmt.Errorf("invalid ViewData(%v)", vo)
	}

	dir := filepath.Dir(vfn)
	errC := os.MkdirAll(dir, os.ModePerm)
	if errC != nil {
		return 0, errC
	}

	f, errF := os.Create(vfn)
	if errF != nil {
		return 0, errF
	}
	defer f.Close()

	vtb["_FILE"] = golua.CreateGoFile(f)

	gl := vm.GetGoLua()
	cc, err3 := gl.ChunkLoad(vm, vsn, true, golua.CreateRenderScriptPreprocess(func(buf *bytes.Buffer) error {
		buf.WriteString("local out = _FILE.Write\n")
		return nil
	}, nil))
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

func (this GOF_gomf_render) IsNative() bool {
	return true
}

func (this GOF_gomf_render) String() string {
	return "GoFunc<gomf.render>"
}
