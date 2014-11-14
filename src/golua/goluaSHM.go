package golua

import (
	"bmautil/valutil"
	"boot"
	"bytes"
	"config"
	"fmt"
	"logger"
	"sync/atomic"
)

func (this *GoLua) GetGlobal(n string) (interface{}, bool) {
	this.globalMutex.RLock()
	defer this.globalMutex.RUnlock()
	v, ok := this.globals[n]
	return v, ok
}

func (this *GoLua) SetGlobal(n string, v interface{}) interface{} {
	this.globalMutex.Lock()
	defer this.globalMutex.Unlock()
	old := this.globals[n]
	this.globals[n] = v
	return old
}

func (this *GoLua) SetObjectFactory(n string, of GoObjectFactory) {
	this.ofMap[n] = of
}

func (this *GoLua) NewObject(n string) (interface{}, error) {
	if of, ok := this.ofMap[n]; ok {
		return of(n)
	}
	return nil, fmt.Errorf("invalid object type '%s'", n)
}

func (this *GoLua) GetConfig(n string) (interface{}, bool) {
	if len(n) > 0 && n[0] == '!' {
		gn := n[1:]
		switch gn {
		case "Debug":
			return boot.Debug, true
		case "DevMode":
			return boot.DevMode, true
		}
		return config.Global.GetConfig(gn)
	}

	this.configMutex.RLock()
	defer this.configMutex.RUnlock()
	v, ok := this.configs[n]
	return v, ok
}

func (this *GoLua) ParseConfig(str string) (string, error) {
	out := bytes.NewBuffer(make([]byte, 0))
	var c1 rune = 0
	word := bytes.NewBuffer(make([]byte, 0))

	this.configMutex.RLock()
	defer this.configMutex.RUnlock()

	for _, c := range []rune(str) {
		switch c1 {
		case 0:
			if c == '$' {
				c1 = c
			} else {
				out.WriteRune(c)
			}
		case '$':
			if c == '{' {
				c1 = '{'
			} else {
				out.WriteRune(c1)
				out.WriteRune(c)
				c1 = 0
			}
		case '{':
			if c == '}' {
				varname := word.String()
				word.Reset()

				var nv interface{}
				ok := false
				if len(varname) > 0 && varname[0] == '!' {
					nv, ok = config.Global.GetConfig(varname[1:])
				} else {
					nv, ok = this.configs[varname]
				}
				if !ok {
					return "", fmt.Errorf("invalid config(%s)", varname)
				}
				out.WriteString(valutil.ToString(nv, ""))
				c1 = 0
			} else {
				word.WriteRune(c)
			}
		}
	}

	if word.Len() > 0 {
		return "", fmt.Errorf("invalid parse format(%s)", str)
	}

	return out.String(), nil
}

func (this *GoLua) SetConfig(n string, v interface{}) bool {
	this.configMutex.Lock()
	defer this.configMutex.Unlock()
	if _, ok := this.configs[n]; ok {
		return false
	}
	this.configs[n] = v
	return true
}

func (this *GoLua) GetService(n string) (interface{}, bool) {
	this.serviceMutex.RLock()
	defer this.serviceMutex.RUnlock()
	v, ok := this.services[n]
	return v, ok
}

func (this *GoLua) SetService(n string, v interface{}) bool {
	this.serviceMutex.Lock()
	defer this.serviceMutex.Unlock()
	if _, ok := this.services[n]; ok {
		return false
	}
	this.services[n] = v
	return true
}

func (this *GoLua) AddService(title string, v interface{}) string {
	this.serviceMutex.Lock()
	defer this.serviceMutex.Unlock()
	for {
		sid := atomic.AddUint32(&this.sid, 1)
		id := fmt.Sprintf("__%d_%s", sid, title)
		if _, ok := this.services[id]; !ok {
			this.services[id] = v
			return id
		}
	}
}

func (this *GoLua) SingletonService(n string, creator func() (interface{}, error)) (interface{}, error) {
	this.serviceMutex.RLock()
	o, ok := this.services[n]
	this.serviceMutex.RUnlock()
	if ok {
		return o, nil
	}
	this.serviceMutex.Lock()
	defer this.serviceMutex.Unlock()
	o, ok = this.services[n]
	if ok {
		return o, nil
	}
	var err error
	o, err = creator()
	if err != nil {
		return nil, err
	}
	this.services[n] = o
	return o, nil
}

func (this *GoLua) RemoveService(id string) bool {
	this.serviceMutex.Lock()
	_, ok := this.services[id]
	if ok {
		delete(this.services, id)
	}
	this.serviceMutex.Unlock()
	return ok
}

func (this *GoLua) CloseService(id string) bool {
	this.serviceMutex.Lock()
	o, ok := this.services[id]
	if ok {
		delete(this.services, id)
	}
	this.serviceMutex.Unlock()
	closed := false
	if tc, ok := o.(SupportTryClose); ok {
		closed = tc.TryClose()
	} else {
		closed = doClose(o)
	}
	if closed {
		logger.Debug(tag, "%s close service '%s'", this, id)
	}
	return ok
}

type GoService struct {
	GL        *GoLua
	SID       string
	CloseFunc func()
	Data      interface{}
	Attr      interface{}
}

func (this *GoService) Close() {
	if this.CloseFunc != nil {
		this.CloseFunc()
	}
	if this.GL != nil {
		this.GL.RemoveService(this.SID)
	}
}

func (this *GoLua) CreateGoService(title string, data interface{}, closeFunc func()) *GoService {
	o := new(GoService)
	o.GL = this
	o.SID = this.AddService(title, o)
	o.CloseFunc = closeFunc
	o.Data = data
	return o
}
