package tankbat

import "fmt"

type DIR byte

const (
	DIR_NONE  = DIR(0)
	DIR_LEFT  = DIR(4)
	DIR_RIGHT = DIR(8)
	DIR_UP    = DIR(1)
	DIR_DOWN  = DIR(2)
)

var (
	DIRS []DIR = []DIR{DIR_UP, DIR_RIGHT, DIR_DOWN, DIR_LEFT}
)

type POS struct {
	x, y int32
}

func (this *POS) POSKey() int64 {
	return int64(this.y)<<32 + int64(this.x)
}

func (this *POS) POSString() string {
	return fmt.Sprintf("%d. %d.", this.x, this.y)
}

func (this *POS) Move(dir DIR, step int32) (x, y int32) {
	switch dir {
	case DIR_UP:
		return this.x, this.y - step
	case DIR_DOWN:
		return this.x, this.y + step
	case DIR_LEFT:
		return this.x - step, this.y
	case DIR_RIGHT:
		return this.x + step, this.y
	default:
		return this.x, this.y
	}
}

type SIZE struct {
	w, h int32
}

func (this *SIZE) SIZEString() string {
	return fmt.Sprintf("%d, %d", this.w, this.h)
}

type RECT struct {
	POS
	SIZE
}

func newRECT(x, y, w, h int32) *RECT {
	return &RECT{POS{x, y}, SIZE{w, h}}
}

type Body struct {
	Pos  POS
	Size SIZE
}

func (this *Body) InitBody(pos POS, sz SIZE) {
	this.Pos = pos
	this.Size = sz
}

func (this *Body) X1() int32 {
	return this.Pos.x
}
func (this *Body) X2() int32 {
	return this.Pos.x + this.Size.w
}
func (this *Body) Y1() int32 {
	return this.Pos.y
}
func (this *Body) Y2() int32 {
	return this.Pos.y + this.Size.h
}
func (this *Body) IsCollide(x1, y1, x2, y2 int32) bool {
	if x1 >= this.X2() || y1 >= this.Y2() || x2 <= this.X1() || y2 <= this.Y1() {
		return false
	} else {
		return true
	}
}
func (this *Body) Center() (int32, int32) {
	return this.Pos.x + this.Size.w/2, this.Pos.y + this.Size.h/2
}
