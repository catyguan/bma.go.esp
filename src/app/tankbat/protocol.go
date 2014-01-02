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
	EVENT_NONE = 0
	EVENT_BOMB = 1
)

type pMapCell struct {
	x, y  int
	kind  int
	event int
}

func (this *pMapCell) Build(m map[string]interface{}) {
	m["x"] = this.x
	m["y"] = this.y
	m["kind"] = this.kind
	if this.event != 0 {
		m["event"] = this.event
	}
}

type pSnapshot struct {
	sid     int
	mapcell []*pMapCell
}

func (this *pSnapshot) Build(m map[string]interface{}) {
	m["sid"] = this.sid
	if this.mapcell != nil {
		a := make([]interface{}, len(this.mapcell))
		m["map"] = a
		for i := 0; i < len(this.mapcell); i++ {
			cm := make(map[string]interface{})
			this.mapcell[i].Build(cm)
			a[i] = cm
		}
	}
}
