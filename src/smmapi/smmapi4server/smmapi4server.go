package smmapi4server

import (
	"bmautil/valutil"
	"boot"
	"fmt"
	"runtime"
	"smmapi"
)

type smmObject int

func s(v uint64) string {
	return valutil.SizeString(v, 1024, valutil.SizeM)
}

func (this smmObject) getProfiles() map[string]interface{} {
	r := make(map[string]interface{})
	r["Version"] = runtime.Version()
	r["StartupTime"] = boot.StartTime.String()
	r["InitTime"] = boot.LoadTime.String()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	r["Memory_Alloc"] = s(ms.Alloc)
	r["Memory_HeapIdle"] = s(ms.HeapIdle)
	r["Memory_Sys"] = s(ms.Sys)

	return r
}

func (this smmObject) getContent() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return fmt.Sprintf("Alloc=%s, HeapIdle=%s, Sys=%s", s(ms.Alloc), s(ms.HeapIdle), s(ms.Sys))
}

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Title = "Server"
	r.Content = this.getContent()

	r.Actions = make([]*smmapi.SMAction, 0)
	if true {
		a := new(smmapi.SMAction)
		a.Id = "boot.detail"
		a.Title = "Detail"
		a.Type = smmapi.SMA_HTTPUI
		a.UIN = "go.server/smm.ui:detail.gl.lua"
		r.Actions = append(r.Actions, a)
	}
	if true {
		a := new(smmapi.SMAction)
		a.Id = "boot.reload"
		a.Title = "Reload"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}

	return r, nil
}

func (this smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	switch aid {
	case "boot.reload":
		msg := "OK"
		if !boot.Restart() {
			msg = "Fail"
		}
		return msg, nil
	case "boot.profiles":
		p := this.getProfiles()
		return p, nil
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}

func init() {
	smmapi.Add("go.server", smmObject(0))
}
