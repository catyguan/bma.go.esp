package tankbat

const (
	MCK_GROUP  = -1
	MCK_NONE   = 0
	MCK_WALL   = 1
	MCK_ROCK   = 2
	MCK_BASE   = 3
	MCK_TANK   = 4
	MCK_BULLET = 5
)

const (
	EVENT_NONE    = 0
	EVENT_CHANGE  = 1
	EVENT_NEW     = 2
	EVENT_DELETE  = 3
	EVENT_DESTROY = 4
)

type fsEvent int

var (
	Event fsEvent
)

func (O fsEvent) New(obj *Object) map[string]interface{} {
	m := make(map[string]interface{})
	m["ekind"] = EVENT_NEW
	m["x"] = obj.X1()
	m["y"] = obj.Y1()
	m["w"] = obj.Size.w
	m["h"] = obj.Size.h
	m["oid"] = obj.Id
	m["okind"] = obj.SKind()
	if obj.Dir != DIR_NONE {
		m["dir"] = obj.Dir
	}
	if obj.Speed > 0 {
		m["speed"] = obj.Speed * MATRIX_SEC_COUNT
	}
	return m
}

func (O fsEvent) ChangeSpeed(obj *Object) map[string]interface{} {
	m := make(map[string]interface{})
	m["ekind"] = EVENT_CHANGE
	m["x"] = obj.X1()
	m["y"] = obj.Y1()
	m["oid"] = obj.Id
	m["okind"] = obj.Kind
	m["speed"] = obj.Speed * MATRIX_SEC_COUNT
	return m
}

func (O fsEvent) ChangeDir(obj *Object) map[string]interface{} {
	m := make(map[string]interface{})
	m["ekind"] = EVENT_CHANGE
	m["x"] = obj.X1()
	m["y"] = obj.Y1()
	m["oid"] = obj.Id
	m["okind"] = obj.Kind
	m["dir"] = obj.Dir
	return m
}

func (O fsEvent) Remove(obj *Object) map[string]interface{} {
	m := make(map[string]interface{})
	if obj.Kind != MCK_ROCK {
		m["ekind"] = EVENT_DELETE
	} else {
		m["ekind"] = EVENT_DESTROY
		m["w"] = obj.Size.w
		m["h"] = obj.Size.h
	}
	m["x"] = obj.X1()
	m["y"] = obj.Y1()
	m["oid"] = obj.SID()
	m["okind"] = obj.SKind()
	return m
}
