package tankbat

import "fmt"

type mapBuilder struct {
	offsetX, offsetY int32
	cw, ch           int32
}

func newMapBuilder(w, h int32, cw, ch int32) *mapBuilder {
	this := new(mapBuilder)
	fmt.Println(w, h, w/2, h/2)
	this.offsetX = w / 2
	this.offsetY = h / 2
	this.cw = cw
	this.ch = ch
	return this
}

func (this *mapBuilder) Cell(cx, cy int32, w, h int32, mdir int) (int32, int32, int32, int32) {
	rx := cx*this.cw - this.offsetX
	ry := cy*this.ch - this.offsetY
	fmt.Println(cx, cy, "=>", rx, ry, "|", this.offsetX, this.offsetY)
	return rx, ry, w, h
}

func (this *mapBuilder) newObject(x, y int32, w, h int32, k int) *Object {
	rx, ry, rw, rh := this.Cell(x, y, w, h, 0)
	return newObject(rx, ry, rw, rh, k)
}

func (this *mapBuilder) newObject4Rock(x, y int32, w, h int32, cw, ch int32) *Object {
	rx, ry, rw, rh := this.Cell(x, y, w, h, 0)
	return newObject4Rock(rx, ry, rw, rh, cw, ch)
}
