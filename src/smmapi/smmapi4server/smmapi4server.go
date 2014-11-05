package smmapi4server

import (
	"bmautil/valutil"
	"boot"
	"bytes"
	"fmt"
	"runtime"
	"smmapi"
)

type smmObject int

func s(v uint64) string {
	return valutil.SizeString(v, 1024, valutil.SizeM)
}

func (this smmObject) getContent() string {
	ver := runtime.Version()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("<pre>")
	buf.WriteString("Memory[")
	buf.WriteString("Alloc=")
	buf.WriteString(s(ms.Alloc))
	buf.WriteString(", ")
	buf.WriteString("TotalAlloc=")
	buf.WriteString(s(ms.TotalAlloc))
	buf.WriteString(", ")
	buf.WriteString("Sys=")
	buf.WriteString(s(ms.Sys))
	buf.WriteString(", ")
	buf.WriteString("Free=")
	buf.WriteString(s(ms.Frees))
	buf.WriteString("]\n")
	buf.WriteString("Version=")
	buf.WriteString(ver)
	buf.WriteString("\n")
	buf.WriteString("</pre>")

	return buf.String()
}

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Title = "Server"
	r.Content = this.getContent()
	r.Actions = make([]*smmapi.SMAction, 0)
	a1 := new(smmapi.SMAction)
	a1.Id = "boot.reload"
	a1.Title = "ReloadServer"
	a1.Type = smmapi.SMA_API
	r.Actions = append(r.Actions, a1)
	return r, nil
}

// Result, refreshInfo, error
func (this smmObject) ExecuteAction(aid, param string) (interface{}, bool, error) {
	switch aid {
	case "boot.reload":
		return boot.Restart(), false, nil
	}
	return nil, false, fmt.Errorf("unknow action(%s)", aid)
}

func init() {
	smmapi.Add("go.server", smmObject(0))
}
