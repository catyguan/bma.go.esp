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

func manageInvoke(cmd string, param map[string]interface{}) (interface{}, error) {
	switch cmd {
	case "list":
		m := List()
		r := make(smlist, 0)
		for k, o := range m {
			info, err := o.GetInfo()
			if err != nil {
				info = new(SMInfo)
				info.Title = "Unknow"
				info.Content = fmt.Sprintf("error : %s", err)
			}
			info.Id = k
			r = append(r, info)
		}
		sort.Sort(r)
		return r, nil
	}
	return nil, fmt.Errorf("unknow command(%s)(%s)", cmd, param)
}

func Invoke(id string, aid string, param map[string]interface{}) (interface{}, error) {
	if id == "" {
		return manageInvoke(aid, param)
	} else {
		obj := Get(id)
		if obj == nil {
			return nil, fmt.Errorf("invalid SMMObject(%s)", id)
		}
		return obj.ExecuteAction(aid, param)
	}
}
