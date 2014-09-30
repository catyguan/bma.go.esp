package glua

import (
	"bmautil/valutil"
	"context"
	"fmt"
	"logger"
	"lua51"
)

func m2v(l *lua51.State) (interface{}, bool) {
	v, _, ok := l.ToGValue(1)
	if !ok {
		return nil, false
	}
	if v == nil {
		return nil, false
	}
	m, ok2 := v.(map[string]interface{})
	if ok2 {
		key := l.ToString(2)
		rv, ok3 := m[key]
		return rv, ok3
	}
	a, ok3 := v.([]interface{})
	if ok3 {
		idx := l.ToInteger(2)
		if idx < len(a) && idx >= 0 {
			return a[idx], true
		}
	}
	return nil, false
}

func msv(l *lua51.State, v interface{}) int {
	gv, gid, ok := l.ToGValue(1)
	if !ok {
		return l.Error("not gvalue")
	}
	if gv == nil {
		return l.Error("nil gvalue")
	}
	m, ok2 := gv.(map[string]interface{})
	if ok2 {
		key := l.ToString(2)
		if v != nil {
			m[key] = v
		} else {
			delete(m, key)
		}
		return 0
	}
	a, ok3 := gv.([]interface{})
	if ok3 {
		idx := l.ToInteger(2)
		if idx < 0 {
			a = append(a, v)
			l.ReplaceGValue(gid, a)
			return 0
		}
		if idx < len(a) && idx >= 0 {
			a[idx] = v
			return 0
		}
		return l.Error(fmt.Sprintf("array out %d/%d", idx, len(a)))
	}
	return l.Error("not map or array")
}

func glua_getBool(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		rv := valutil.ToBool(v, false)
		l.PushBoolean(rv)
		return 1
	}
	l.PushNil()
	return 1
}

func glua_getInt(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		rv := valutil.ToInt(v, 0)
		l.PushInteger(rv)
		return 1
	}
	l.PushNil()
	return 1
}

func glua_getNumber(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		rv := valutil.ToFloat64(v, 0)
		l.PushNumber(rv)
		return 1
	}
	l.PushNil()
	return 1
}

func glua_getString(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		rv := valutil.ToString(v, "")
		l.PushString(rv)
		return 1
	}
	l.PushNil()
	return 1
}

func glua_getMap(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		if rv, ok2 := v.(map[string]interface{}); ok2 {
			l.PushGValue(rv)
		}
		return 1
	}
	l.PushNil()
	return 1
}

func glua_getArray(l *lua51.State) int {
	v, ok := m2v(l)
	if ok {
		if rv, ok2 := v.([]interface{}); ok2 {
			l.PushGValue(rv)
		}
		return 1
	}
	l.PushNil()
	return 1
}

func glua_setNil(l *lua51.State) int {
	return msv(l, nil)
}

func glua_setBool(l *lua51.State) int {
	v := l.ToBoolean(3)
	return msv(l, v)
}

func glua_setInt(l *lua51.State) int {
	v := l.ToInteger(3)
	return msv(l, v)
}

func glua_setNumber(l *lua51.State) int {
	v := l.ToNumber(3)
	return msv(l, v)
}

func glua_setString(l *lua51.State) int {
	v := l.ToString(3)
	return msv(l, v)
}

func glua_setMap(l *lua51.State) int {
	v, _, ok := l.ToGValue(3)
	if ok {
		m, ok2 := v.(map[string]interface{})
		if ok2 {
			return msv(l, m)
		}
	}
	return l.Error("val not map")
}

func glua_setArray(l *lua51.State) int {
	v, _, ok := l.ToGValue(3)
	if ok {
		a, ok2 := v.([]interface{})
		if ok2 {
			return msv(l, a)
		}
	}
	return l.Error("val not array")
}

func glua_newMap(l *lua51.State) int {
	m := make(map[string]interface{})
	l.PushGValue(m)
	return 1
}

func luav(l *lua51.State, idx int) interface{} {
	switch l.Type(idx) {
	case lua51.LUA_TBOOLEAN:
		return l.ToBoolean(idx)
	case lua51.LUA_TNUMBER:
		v1 := l.ToNumber(idx)
		v2 := l.ToInteger(idx)
		if v1 == float64(v2) {
			return v2
		}
		return v1
	case lua51.LUA_TSTRING:
		return l.ToString(idx)
	case lua51.LUA_TTABLE:
		return luam(l, idx)
	default:
		return nil
	}
}

func luam(l *lua51.State, idx int) map[string]interface{} {
	r := make(map[string]interface{})
	l.PushNil()
	for {
		if l.Next(idx) == 0 {
			break
		}
		k := l.ToString(-2)
		v := luav(l, l.GetTop())
		if v != nil {
			r[k] = v
		}
		l.Pop(1)
	}
	return r
}

func glua_toMap(l *lua51.State) int {
	if !l.IsTable(1) {
		return l.Error("not table")
	}
	m := luam(l, 1)
	l.PushGValue(m)
	return 1
}

func glua_newArray(l *lua51.State) int {
	a := make([]interface{}, 0)
	l.PushGValue(a)
	return 1
}

func glua_toArray(l *lua51.State) int {
	return 0
}

func (this *GLua) doPrint(l *lua51.State) int {
	if logger.EnableDebug(tag) {
		s := lua51.CheckString(this.l, 1)
		logger.Debug(tag, "'%s' print >> %s", this.name, s)
	}
	return 0
}

func (this *GLua) doTask(l *lua51.State) int {
	taskName := l.ToString(1)
	var req map[string]interface{}
	v, _, ok := l.ToGValue(2)
	if ok {
		req, ok = v.(map[string]interface{})
		if !ok {
			return l.Error("task request not map")
		}
	}
	ctx := this.context
	f := l.ToString(3)
	cb := func(taskName string, ctx context.Context, cu ContextUpdater, err error) {
		this.luaCallback4Task(taskName, f, ctx, cu, err)
	}
	err := this.StartTask(taskName, ctx, req, cb)
	if err != nil {
		return l.Error(err.Error())
	}
	return 0
}

func (this *GLua) initGoFunctions() {
	l := this.l
	l.Register("glua_print", this.doPrint)
	l.Register("glua_task", this.doTask)

	l.Register("glua_getBool", glua_getBool)
	l.Register("glua_getInt", glua_getInt)
	l.Register("glua_getNumber", glua_getNumber)
	l.Register("glua_getString", glua_getString)
	l.Register("glua_getMap", glua_getMap)
	l.Register("glua_getArray", glua_getArray)
	l.Register("glua_setNil", glua_setNil)
	l.Register("glua_setBool", glua_setBool)
	l.Register("glua_setInt", glua_setInt)
	l.Register("glua_setNumber", glua_setNumber)
	l.Register("glua_setString", glua_setString)
	l.Register("glua_setMap", glua_setMap)
	l.Register("glua_setArray", glua_setArray)
	l.Register("glua_newMap", glua_newMap)
	l.Register("glua_toMap", glua_toMap)
	l.Register("glua_newArray", glua_newArray)
	l.Register("glua_toArray", glua_toArray)

	l.Register("glua_print", this.doPrint)
}
