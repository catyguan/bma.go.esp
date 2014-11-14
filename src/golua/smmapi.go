package golua

import (
	"boot"
	"fmt"
	"smmapi"
)

func (this *GoLua) InitSMMApi() {
	o := new(smmObject)
	o.gl = this
	smmapi.Add(fmt.Sprintf("golua.%s", this.name), o)
}

func (this *GoLua) CloseSMMApi() {
	smmapi.Remove(fmt.Sprintf("golua.%s", this.name))
}

////////////////////////////
type smmObject struct {
	gl *GoLua
}

func (this smmObject) getContent() string {
	return fmt.Sprintf("%s", this.gl)
}

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	var extinfo *smmapi.SMInfo
	if this.gl.ExtSMMApi != nil {
		var err error
		extinfo, err = this.gl.ExtSMMApi.GetInfo()
		if err != nil {
			return nil, err
		}
	}

	r := new(smmapi.SMInfo)
	r.Title = "GoLua App"
	r.Content = this.getContent()
	if extinfo != nil && extinfo.Content != "" {
		r.Content = r.Content + ";" + extinfo.Content
	}

	r.Actions = make([]*smmapi.SMAction, 0)
	if extinfo != nil {
		for _, a := range extinfo.Actions {
			r.Actions = append(r.Actions, a)
		}
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "golua.debugger"
		a.Title = "Debugger"
		a.Type = smmapi.SMA_HTTPUI
		a.UIN = "golua/smm.ui:debugger.gl.lua"
		r.Actions = append(r.Actions, a)
	}

	return r, nil
}

// Result, refreshInfo, error
func (this smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	switch aid {
	case "boot.reload":
		msg := "OK"
		if !boot.Restart() {
			msg = "Fail"
		}
		return msg, nil
	}
	if this.gl.ExtSMMApi != nil {
		return this.gl.ExtSMMApi.ExecuteAction(aid, param)
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}
