package smmapi

import (
	"fmt"
	"sort"
	"sync"
)

type service struct {
	sync.RWMutex
	objects map[string]SMMObject
}

var (
	gs service
)

func Add(id string, obj SMMObject) {
	gs.Lock()
	defer gs.Unlock()
	if gs.objects == nil {
		gs.objects = make(map[string]SMMObject)
	}
	gs.objects[id] = obj
}

func Remove(id string) {
	gs.Lock()
	defer gs.Unlock()
	if gs.objects != nil {
		delete(gs.objects, id)
	}
}

func Get(id string) SMMObject {
	gs.RLock()
	defer gs.RUnlock()
	if gs.objects == nil {
		return nil
	}
	if o, ok := gs.objects[id]; ok {
		return o
	}
	return nil
}

func List() map[string]SMMObject {
	r := make(map[string]SMMObject)
	gs.RLock()
	defer gs.RUnlock()
	for k, o := range gs.objects {
		r[k] = o
	}
	return r
}

func Invoke(id string, aid string, param string) (interface{}, bool, error) {
	if id == "" {
		switch aid {
		case "list":
			m := List()
			r := make(smlist, 0)
			for k, o := range m {
				info, err := o.GetInfo()
				if err != nil {
					info = new(SMInfo)
					info.Title = "Unknow"
					info.Content = fmt.Sprintf("<pre>error : %s</pre>", err)
				}
				info.Id = k
				r = append(r, info)
			}
			sort.Sort(r)
			return r, false, nil
		case "one":
			o := Get(param)
			if o == nil {
				return nil, false, fmt.Errorf("invalid SMMObject(%s)", param)
			}
			info, err := o.GetInfo()
			if info != nil {
				info.Id = param
			}
			return info, false, err
		}
		return nil, false, fmt.Errorf("unknow command(%s)(%s)", aid, param)
	} else {
		obj := Get(id)
		if obj == nil {
			return nil, false, fmt.Errorf("invalid SMMObject(%s)", id)
		}
		return obj.ExecuteAction(aid, param)
	}
}
