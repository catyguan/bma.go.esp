package tankbat

import "fmt"

type World struct {
	Objects  map[int]*Object
	SObj     map[int]*Object
	GObj     map[int]*Object
	MObj     map[int]*Object
	objectId int
}

func NewWorld() *World {
	this := new(World)
	this.Objects = make(map[int]*Object)
	this.SObj = make(map[int]*Object)
	this.GObj = make(map[int]*Object)
	this.MObj = make(map[int]*Object)
	return this
}

func (this *World) Add(o *Object) bool {
	if o.Id == 0 {
		o.Id = this.NextId()
	}
	_, ok := this.Objects[o.Id]
	if ok {
		return false
	}
	this.Objects[o.Id] = o
	switch o.Kind {
	case MCK_TANK, MCK_BULLET:
		this.MObj[o.Id] = o
	case MCK_GROUP:
		this.GObj[o.Id] = o
		if o.child != nil {
			for _, c := range o.child {
				this.Add(c)
			}
		}
	case MCK_WALL, MCK_BASE:
		this.SObj[o.Id] = o
	}
	return true
}

func (this *World) Remove(o *Object) {
	delete(this.Objects, o.Id)
	delete(this.SObj, o.Id)
	delete(this.GObj, o.Id)
	delete(this.MObj, o.Id)
}

func (this *World) NextId() int {
	this.objectId++
	return this.objectId
}

func (this *World) Collide(res map[*Object]bool, x1, y1, x2, y2 int32, exc *Object, rockOnly bool) {
	if !rockOnly {
		for _, o := range this.SObj {
			if o.IsCollide(x1, y1, x2, y2) {
				res[o] = true
			}
		}
	}
	for _, o := range this.GObj {
		o.Collide(res, x1, y1, x2, y2)
		if o.removed {
			this.Remove(o)
		}
	}
	if !rockOnly {
		for _, o := range this.MObj {
			if o != exc && o.IsCollide(x1, y1, x2, y2) {
				res[o] = true
			}
		}
	}
}

func (this *World) BuildDumpMap(dmap *dumpMap) {
	dmap.Clear()
	for _, o := range this.Objects {
		if o.Kind == MCK_GROUP {
			continue
		}
		dmap.Set(o.Pos.x, o.Pos.y, o.Size.w, o.Size.h, o.Kind)
	}
}

type Object struct {
	Id   int
	Kind int
	Name string
	Flag int

	Body
	Dir   DIR
	Speed int

	child []*Object

	removed bool
}

func newObject(x1, y1, w, h int32, k int) *Object {
	r := new(Object)
	r.Kind = k
	switch k {
	case MCK_BASE:
		r.Name = "BASE"
	case MCK_BULLET:
		r.Name = "BULLET"
	case MCK_GROUP:
		r.Name = "GROUP"
	case MCK_ROCK:
		r.Name = "ROCK"
	case MCK_TANK:
		r.Name = "TANK"
	case MCK_WALL:
		r.Name = "WALL"
	}
	r.InitBody(POS{x1, y1}, SIZE{w, h})
	return r
}

func newObject4Rock(x1, y1, w, h int32, cw, ch int32) *Object {
	r := newObject(x1, y1, w, h, MCK_GROUP)
	for y := y1; y < y1+h; y += ch {
		rh := ch
		if y+rh > h {
			rh = h - y
		}
		for x := x1; x < x1+w; x += cw {
			rw := cw
			if x+rw > w {
				rw = w - x
			}
			r.AddChild(newObject(x, y, rw, rh, MCK_ROCK))
		}
	}
	return r
}

func (this *Object) AddChild(o *Object) {
	if this.Kind != MCK_GROUP {
		panic("invalid kind")
	}
	if this.child == nil {
		this.child = make([]*Object, 0)
	}
	this.child = append(this.child, o)
}

func (this *Object) String() string {
	return fmt.Sprintf("%s[%s:%s]", this.Name, this.Body.Pos.POSString(), this.Body.Size.SIZEString())
}

func (this *Object) Moving() bool {
	return this.Dir != DIR_NONE && this.Speed > 0
}

func (this *Object) TurnDir(dir DIR) {
	this.Dir = dir
}

func (this *Object) StopMove() {
	this.Speed = 0
}

func (this *Object) StartMove(sp int) {
	this.Speed = sp
}

func (this *Object) Collide(res map[*Object]bool, x1, y1, x2, y2 int32) {
	if this.removed {
		return
	}
	if !this.IsCollide(x1, y1, x2, y2) {
		return
	}
	if this.Kind == MCK_GROUP {
		valid := false
		if this.child != nil {
			for _, o := range this.child {
				if o.removed {
					continue
				}
				valid = true
				o.Collide(res, x1, y1, x2, y2)
			}
		}
		if !valid {
			this.removed = true
		}
	} else {
		res[this] = true
	}
}
