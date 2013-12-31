package bombman

import (
	"bytes"
	"fmt"
	"math/rand"
)

type MapPos struct {
	x, y int
}

func (this *MapPos) PosKey() int {
	return this.y*10000 + this.x
}

func (this *MapPos) PosString() string {
	return fmt.Sprintf("%d, %d", this.x, this.y)
}

func (this *MapPos) Dir(dir int) MapPos {
	var pos MapPos
	pos = *this
	switch dir {
	case DIR_DOWN:
		pos.y++
	case DIR_UP:
		pos.y--
	case DIR_LEFT:
		pos.x--
	case DIR_RIGHT:
		pos.x++
	}
	return pos
}

type MapCell struct {
	kind   byte
	event  int
	seqid  int
	player *Player
	bomb   *bombInfo
}

func (this *MapCell) View() string {
	if this.player != nil {
		return this.player.View()
	}
	switch this.kind {
	case MCK_WALL:
		return "#"
	case MCK_ROCK:
		return "M"
	case MCK_BOMB:
		return "@"
	default:
		return " "
	}
}

func (this *MapCell) IsEmpty() bool {
	if this.kind != MCK_NONE {
		return false
	}
	return true
}

type Map struct {
	width, height int
	cell          [][]MapCell
}

func (this *Map) InitMap(w, h int) {
	this.width = w
	this.height = h
	this.cell = make([][]MapCell, w)
	for i := 0; i < w; i++ {
		this.cell[i] = make([]MapCell, h)
	}
	spos := func(v, mv int) bool {
		if v <= 1 {
			return true
		}
		if v >= mv-2 {
			return true
		}
		return false
	}
	free := make([]*MapCell, 0)
	for y := 0; y < this.height; y++ {
		for x := 0; x < this.width; x++ {
			if x%2 == 1 && y%2 == 1 {
				this.cell[x][y].kind = MCK_WALL
			} else {
				if !(spos(x, this.width) && spos(y, this.height)) {
					free = append(free, &this.cell[x][y])
				}
			}
		}
	}

	l := len(free)
	for i := 0; i < len(free)/3; i++ {
		pos := rand.Intn(l)
		cell := free[pos]
		cell.kind = MCK_ROCK
		free[pos], free[l-1] = free[l-1], free[pos]
		l--
	}

}

func (this *Map) Cell(x, y int) *MapCell {
	if this.validPos(x, y) {
		return &this.cell[x][y]
	}
	return nil
}

func (this *Map) View() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	for x := 0; x < this.width+2; x++ {
		buf.WriteString("#")
	}
	buf.WriteString("\n")
	for y := 0; y < this.height; y++ {
		buf.WriteString("#")
		for x := 0; x < this.width; x++ {
			buf.WriteString(this.cell[x][y].View())
		}
		buf.WriteString("#\n")
	}
	for x := 0; x < this.width+2; x++ {
		buf.WriteString("#")
	}
	buf.WriteString("\n")
	return buf.String()
}

func (this *Map) PlayerIn(x, y int, pl *Player) {
	cell := &this.cell[x][y]
	cell.player = pl
	cell.kind = MCK_PLAYER
	pl.x = x
	pl.y = y
}

func (this *Map) validPos(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	return x < this.width && y < this.height
}

func (this *Map) PlayerMove(pl *Player, dir int, sid int) bool {
	pos := pl.MapPos.Dir(dir)
	if !this.validPos(pos.x, pos.y) {
		return false
	}
	cell := &this.cell[pos.x][pos.y]
	if cell.kind == MCK_NONE {
		old := &this.cell[pl.x][pl.y]
		old.kind = MCK_NONE
		old.player = nil
		old.seqid = sid
		cell.kind = MCK_PLAYER
		cell.player = pl
		cell.seqid = sid
		pl.MapPos = pos
		return true
	}
	return false
}

func (this *Map) PlayerBomb(pl *Player, dir int, b *bombInfo) *MapPos {
	pos := pl.MapPos.Dir(dir)
	if this.BombPlace(&pos, b) {
		return &pos
	}
	return nil
}

func (this *Map) BombPlace(pos *MapPos, b *bombInfo) bool {
	if !this.validPos(pos.x, pos.y) {
		return false
	}
	cell := &this.cell[pos.x][pos.y]
	if cell.kind == MCK_NONE {
		cell.kind = MCK_BOMB
		cell.bomb = b
		cell.seqid = b.seqid
		return true
	}
	return false
}

func (this *Map) PlayerDie(pl *Player) {
	cell := this.Cell(pl.x, pl.y)
	if cell != nil && cell.player == pl {
		cell.kind = MCK_NONE
		cell.player = nil
	}
}

func (this *Map) ToPMapCell(x, y int, seqid int) *pMapCell {
	cell := &this.cell[x][y]
	mc := new(pMapCell)
	mc.kind = int(cell.kind)
	mc.x = x
	mc.y = y
	if cell.seqid == seqid {
		mc.event = cell.event
	}
	return mc
}

func (this *Map) MakeSnapshot(seqid int) []*pMapCell {
	r := make([]*pMapCell, 0)
	for y := 0; y < this.height; y++ {
		for x := 0; x < this.width; x++ {
			cell := &this.cell[x][y]
			if cell.IsEmpty() {
				continue
			}
			r = append(r, this.ToPMapCell(x, y, 0))
		}
	}
	return r
}

func (this *Map) MakeEvent(seqid int) []*pMapCell {
	r := make([]*pMapCell, 0)
	for y := 0; y < this.height; y++ {
		for x := 0; x < this.width; x++ {
			cell := &this.cell[x][y]
			if cell.seqid == seqid {
				if cell.event == EVENT_BOMB && cell.kind == MCK_ROCK {
					cell.kind = MCK_NONE
				}
				r = append(r, this.ToPMapCell(x, y, seqid))
			}
			cell.event = 0
		}
	}
	return r
}
