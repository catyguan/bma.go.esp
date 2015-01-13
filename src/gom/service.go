package gom

import (
	"bytes"
	"context"
	"golua"
)
import (
	"fileloader"
	"fmt"
)

type Service struct {
	name    string
	config  *configInfo
	floader fileloader.FileLoader
	gl      *golua.GoLua
	gli     golua.GoLuaInitor
}

func NewService(name string, gli golua.GoLuaInitor) *Service {
	this := new(Service)
	this.name = name
	this.gli = gli
	return this
}

func (this *Service) RunCommands(fname string, script string, params []string) error {
	var content string
	if fname != "none" {
		bs, err0 := this.floader.Load(fname)
		if err0 != nil {
			return err0
		}
		if bs == nil {
			return fmt.Errorf("dom file '%s' not exists", fname)
		}
		content = string(bs)
	}
	gm := NewGOM()
	if content != "" {
		err1 := gm.Compile(content, fname)
		if err1 != nil {
			return err1
		}
	}
	if script == "" || script == "dump" {
		buf := bytes.NewBuffer([]byte{})
		gm.Dump(buf, "")
		fmt.Println("----------------- DUMP ----------------")
		fmt.Println(buf.String())
		fmt.Println("----------------- DUMP END-------------")
		return nil
	}

	ri := golua.NewRequestInfo()
	ri.Script = script

	ctx := context.Background()
	ctx, _ = context.CreateExecId(ctx)
	ctx = golua.CreateRequest(ctx, ri)

	locals := make(map[string]interface{})
	locals["gom"] = gm
	locals["PARAMS"] = params
	_, errE := this.gl.DoExecute(ctx, locals)

	return errE
}
