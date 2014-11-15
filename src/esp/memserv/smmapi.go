package memserv

import (
	"bmautil/valutil"
	"fmt"
	"smmapi"
)

type smmObject struct {
	s *MemoryServ
}

func sizestr(v uint64) string {
	return valutil.SizeString(v, 1024, valutil.SizeM)
}

func (this *smmObject) getContent() string {
	var sz1 int
	var sz2 uint64
	this.s.lock.RLock()
	defer this.s.lock.RUnlock()
	c := len(this.s.mgs)
	for _, m := range this.s.mgs {
		s1, s2 := m.Size()
		sz1 += s1
		sz2 += uint64(s2)
	}
	return fmt.Sprintf("MemGo=%d, Items=%d, MemSize=%s", c, sz1, sizestr(sz2))
}

func (this *smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Title = "MemServ"
	r.Content = this.getContent()

	r.Actions = make([]*smmapi.SMAction, 0)
	if true {
		a := new(smmapi.SMAction)
		a.Id = "memserv.detail"
		a.Title = "Detail"
		a.Type = smmapi.SMA_HTTPUI
		a.UIN = "go.memserv/smm.ui:detail.gl.lua"
		r.Actions = append(r.Actions, a)
	}

	return r, nil
}

func (this *smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	switch aid {
	case "memserv.list":
		return 0, nil
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}

func (this *MemoryServ) InitSMMAPI(n string) {
	smmapi.Add(n, &smmObject{this})
}
