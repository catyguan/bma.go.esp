package memserv4session

import (
	"acl"
	"bmautil/valutil"
	"esp/memserv"
	"strings"
	"time"
)

const (
	KEY_SESSION = "session"
)

var (
	userKeys []string
)

func init() {
	userKeys = []string{"USER_ACCOUNT", "USER_DOMAIN", "USER_GROUPS"}
}

func QueryUser(s *memserv.MemoryServ, sid string, timeoutMS int) (*acl.User, error) {
	datas, err := MGetSession(s, sid, userKeys, timeoutMS)
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
		return nil, nil
	}
	a := valutil.ToString(datas["USER_ACCOUNT"], "")
	if a == "" {
		a = "anonymous"
	}
	d := valutil.ToString(datas["USER_DOMAIN"], "")
	g := valutil.ToString(datas["USER_GROUPS"], "")
	var gs []string
	if g != "" {
		gs = strings.Split(g, ",")
	}
	user := acl.NewUser(a, d, gs)
	return user, nil
}

func GetSession(s *memserv.MemoryServ, sid string, key string, timeoutMS int) (interface{}, error) {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return nil, err
	}
	var rv interface{}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		tm := time.Now()
		ok, v, err0 := mgi.Get(sid, &tm)
		if err0 != nil {
			return err0
		}
		if ok {
			mgi.Touch(sid, timeoutMS)
			if m, ok2 := v.(map[string]interface{}); ok2 {
				if v2, ok3 := m[key]; ok3 {
					rv = v2
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rv, err
}

func MGetSession(s *memserv.MemoryServ, sid string, keys []string, timeoutMS int) (map[string]interface{}, error) {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return nil, err
	}
	rv := make(map[string]interface{})
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		tm := time.Now()
		ok, v, err0 := mgi.Get(sid, &tm)
		if err0 != nil {
			return err0
		}
		if ok {
			mgi.Touch(sid, timeoutMS)
			if m, ok2 := v.(map[string]interface{}); ok2 {
				for _, k := range keys {
					if v2, ok3 := m[k]; ok3 {
						rv[k] = v2
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rv, err
}

func SetUser(s *memserv.MemoryServ, sid string, user *acl.User, timeoutMS int) error {
	if user == nil {
		return MDeleteSession(s, sid, userKeys, timeoutMS)
	} else {
		mv := make(map[string]interface{})
		mv["USER_ACCOUNT"] = user.Account
		mv["USER_DOMAIN"] = user.Domain
		if user.Groups != nil {
			mv["USER_GROUPs"] = strings.Join(user.Groups, ",")
		}
		return MSetSession(s, sid, mv, timeoutMS)
	}
}

func SetSession(s *memserv.MemoryServ, sid string, key string, val interface{}, timeoutMS int) error {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return err
	}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		ok, v, err0 := mgi.Get(sid, nil)
		if err0 != nil {
			return err0
		}
		var m map[string]interface{}
		if ok {
			m, ok = v.(map[string]interface{})
		}
		if m == nil {
			m = make(map[string]interface{})
		}
		m[key] = val
		return mgi.Set(sid, m, timeoutMS)
	})
	return err
}

func MSetSession(s *memserv.MemoryServ, sid string, mv map[string]interface{}, timeoutMS int) error {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return err
	}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		ok, v, err0 := mgi.Get(sid, nil)
		if err0 != nil {
			return err0
		}
		var m map[string]interface{}
		if ok {
			m, ok = v.(map[string]interface{})
		}
		if m == nil {
			m = make(map[string]interface{})
		}
		for key, val := range mv {
			m[key] = val
		}
		return mgi.Set(sid, m, timeoutMS)
	})
	return err
}

func DeleteSession(s *memserv.MemoryServ, sid string, key string, timeoutMS int) error {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return err
	}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		ok, v, err0 := mgi.Get(sid, nil)
		if err0 != nil {
			return err0
		}
		if ok {
			if m, ok2 := v.(map[string]interface{}); ok2 {
				delete(m, key)
				return mgi.Set(sid, m, timeoutMS)
			}
		}
		return nil
	})
	return err
}

func MDeleteSession(s *memserv.MemoryServ, sid string, keys []string, timeoutMS int) error {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return err
	}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		ok, v, err0 := mgi.Get(sid, nil)
		if err0 != nil {
			return err0
		}
		if ok {
			if m, ok2 := v.(map[string]interface{}); ok2 {
				for _, key := range keys {
					delete(m, key)
				}
				return mgi.Set(sid, m, timeoutMS)
			}
		}
		return nil
	})
	return err
}

func CloseSession(s *memserv.MemoryServ, sid string) error {
	mg, _, err := s.GetOrCreate(KEY_SESSION, nil)
	if err != nil {
		return err
	}
	err = mg.DoSync(func(mgi *memserv.MemGoI) error {
		return mgi.Remove(sid)
	})
	return err
}
