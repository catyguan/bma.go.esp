package tankbat

import (
	"bytes"
	"fmt"
	"math"
)

type dumpMap struct {
	cells [][]int

	offsetX, offsetY int32
	mw, mh           int32
	cw, ch           int32
}

func newDumpMap(w, h int32, cw, ch int32) *dumpMap {
	this := new(dumpMap)

	this.offsetX = w / 2
	this.offsetY = h / 2
	this.mw = w / cw
	this.mh = h / ch
	this.cw = cw
	this.ch = ch

	this.cells = make([][]int, this.mw)
	for x := int32(0); x < this.mw; x++ {
		this.cells[x] = make([]int, this.mh)
	}
	return this
}

func (this *dumpMap) Clear() {
	for y := int32(0); y < this.mh; y++ {
		for x := int32(0); x < this.mw; x++ {
			this.cells[x][y] = MCK_NONE
		}
	}
}

func (this *dumpMap) View() string {
	buf := bytes.NewBuffer([]byte{})
	for y := int32(0); y < this.mh; y++ {
		buf.WriteString(fmt.Sprintf("%2d: ", y))
		for x := int32(0); x < this.mw; x++ {
			switch this.cell(x, y) {
			case MCK_NONE:
				buf.WriteByte(' ')
			case MCK_WALL:
				buf.WriteByte('#')
			case MCK_ROCK:
				buf.WriteByte('M')
			case MCK_BASE:
				buf.WriteByte('*')
			case MCK_TANK:
				buf.WriteByte('T')
			case MCK_BULLET:
				buf.WriteByte('O')
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (this *dumpMap) W2MX(v int32) int32 {
	return (v + this.offsetX) / this.cw
}
func (this *dumpMap) W2MY(v int32) int32 {
	return (v + this.offsetY) / this.ch
}
func (this *dumpMap) W2MW(v int32) int32 {
	return int32(math.Ceil(float64(v) / float64(this.cw)))
}
func (this *dumpMap) W2MH(v int32) int32 {
	return int32(math.Ceil(float64(v) / float64(this.ch)))
}

func (this *dumpMap) cell(x, y int32) int {
	if y >= int32(0) && x >= int32(0) && x < this.mw && y < this.mh {
		return this.cells[x][y]
	}
	return MCK_NONE
}

func (this *dumpMap) Set(x, y, w, h int32, k int) {
	mx := this.W2MX(x)
	my := this.W2MY(y)
	mw := this.W2MW(w)
	mh := this.W2MH(h)
	if mw < 1 {
		mw = 1
	}
	if mh < 1 {
		mh = 1
	}
	// if k == MCK_BULLET {
	// 	fmt.Println("set", x, y, w, h, mx, my, mw, mh)
	// }
	for ix := int32(0); ix < mw; ix++ {
		for iy := int32(0); iy < mh; iy++ {
			rx := mx + ix
			ry := my + iy
			if ry >= int32(0) && rx >= int32(0) && rx < this.mw && ry < this.mh {
				kv := this.cells[rx][ry]
				if k > kv {
					this.cells[rx][ry] = k
				}
			}
		}
	}
}
