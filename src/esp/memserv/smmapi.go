package memserv

import (
	"bmautil/valutil"
	"fmt"
	"smmapi"
	"strings"
)

type smmObject struct {
	s *MemoryServ
}

func sizestr(v uint64) string {
	return valutil.SizeString(v, 1024, valutil.SizeAny)
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
		r := make([]map[string]interface{}, 0)
		this.s.lock.RLock()
		defer this.s.lock.RUnlock()
		for _, m := range this.s.mgs {
			o := make(map[string]interface{})
			s1, s2 := m.Size()
			o["Name"] = m.name
			o["Count"] = s1
			o["Size"] = sizestr(uint64(s2))
			r = append(r, o)
		}
		return r, nil
	case "memserv.delete":
		n := ""
		key := ""
		if param != nil {
			n = valutil.ToString(param["name"], "")
			key = valutil.ToString(param["key"], "")
		}
		mg := this.s.Get(n)
		if mg != nil {
			mg.DoSync(func(mgi *MemGoI) error {
				mgi.Remove(key)
				return nil
			})
		}
		return "OK", nil
	case "memserv.dump":
		n := ""
		str := ""
		s := 0
		c := 100
		if param != nil {
			n = valutil.ToString(param["name"], "")
			str = valutil.ToString(param["filter"], "")
			s = valutil.ToInt(param["start"], s)
			c = valutil.ToInt(param["count"], c)
		}

		r := make(map[string]interface{})
		mg := this.s.Get(n)
		if mg == nil {
			r["Message"] = fmt.Sprintf("Invalid MemGo('%s')", n)
			return r, nil
		}
		res := make([]map[string]interface{}, 0)
		r["Data"] = res

		sn := fmt.Sprintf("smmapi_%p", r)
		err0 := mg.BeginScan(sn)
		if err0 != nil {
			return nil, err0
		}
		defer mg.EndScan(sn)
		idx := 0
		for {
			ok, err1 := mg.Scan(sn, 100, func(k string, v interface{}) {
				defer func() {
					idx++
				}()
				if idx < s {
					return
				}
				if str != "" {
					if strings.Index(k, str) == -1 {
						return
					}
				}

				o := make(map[string]interface{})
				o["Key"] = k
				o["Value"] = fmt.Sprintf("%v", v)
				list := r["Data"].([]map[string]interface{})
				list = append(list, o)
				r["Data"] = list
			})
			if err1 != nil {
				return nil, err1
			}
			if ok {
				break
			}
			list := r["Data"].([]map[string]interface{})
			if len(list) >= c {
				break
			}
		}
		return r, nil
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}

func (this *MemoryServ) InitSMMAPI(n string) {
	smmapi.Add(n, &smmObject{this})
}
