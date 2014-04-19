package tankbat

import (
	"fmt"
	"logger"
	"math/rand"
	"time"
)

func (this *Matrix) doInitMatrix() error {
	rand.Seed(time.Now().UnixNano())
	this.seqid = 1

	ww := MATRIX_WORLD_WIDTH + MATRIX_UNIT*2
	wh := MATRIX_WORLD_HEIGHT + MATRIX_UNIT*2
	this.dmap = newDumpMap(ww, wh, MATRIX_UNIT/2, MATRIX_UNIT/2)
	this.builder = newMapBuilder(ww, wh, MATRIX_UNIT/10, MATRIX_UNIT/10)

	w := NewWorld()
	this.world = w

	// layout
	w.Add(this.builder.newObject(0, 0, ww, MATRIX_UNIT, MCK_WALL))
	w.Add(this.builder.newObject(0, (MATRIX_MAP_HEIGHT+2-1)*10, ww, MATRIX_UNIT, MCK_WALL))

	w.Add(this.builder.newObject(0, 0, MATRIX_UNIT, wh-1*MATRIX_UNIT, MCK_WALL))
	w.Add(this.builder.newObject((MATRIX_MAP_WIDTH+2-1)*10, 0, MATRIX_UNIT, wh-1*MATRIX_UNIT, MCK_WALL))

	walls := []*RECT{
		newRECT(70, 70, 10, 10),

		newRECT(10, 75, 10, 5),
		newRECT(30, 70, 10, 5),
		newRECT(50, 75, 10, 5),

		newRECT(90, 70, 10, 5),
		newRECT(110, 75, 10, 5),
		newRECT(130, 70, 10, 5),
	}
	for _, info := range walls {
		w.Add(this.builder.newObject(info.x, info.y, info.w*MATRIX_UNIT/10, info.h*MATRIX_UNIT/10, MCK_WALL))
	}

	rocks := []*RECT{

		newRECT(65, 10, 5, 10),
		newRECT(80, 10, 5, 10),
		newRECT(65, 20, 20, 5),

		newRECT(20, 20, 10, 35),
		newRECT(40, 20, 10, 35),
		newRECT(100, 20, 10, 35),
		newRECT(120, 20, 10, 35),
	}
	for _, info := range rocks {
		w.Add(this.builder.newObject4Rock(info.x, info.y, info.w*MATRIX_UNIT/10, info.h*MATRIX_UNIT/10, MATRIX_ROCK_SIZE, MATRIX_ROCK_SIZE))
	}

	// w.AddBaseToMap((70-10)*4/10+2, (10-10)*4/10+2, 4, 4, 1)
	base := this.builder.newObject(70, 10, 10*MATRIX_UNIT/10, 10*MATRIX_UNIT/10, MCK_BASE)
	base.Flag = 1
	w.Add(base)

	// test
	w.BuildDumpMap(this.dmap)
	fmt.Print(this.dmap.View())
	return nil
}

func (this *Matrix) doAddTank(pl *Player, event bool) {
	var x, y int32
	var dir DIR
	switch pl.teamId {
	case 1:
		dir = DIR_DOWN
		if pl.teamNum == 0 {
			x = 10
			y = 10
		} else {
			x = 10 + MATRIX_MAP_WIDTH*10
			y = 10
		}
	case 2:
		dir = DIR_UP
		if pl.teamNum == 0 {
			x = 10
			y = MATRIX_MAP_HEIGHT * 10
		} else {
			x = 10 + MATRIX_MAP_WIDTH*10
			y = MATRIX_MAP_HEIGHT * 10
		}
	}
	tank := this.builder.newObject(x, y, MATRIX_UNIT, MATRIX_UNIT, MCK_TANK)
	tank.Flag = pl.teamId
	tank.TurnDir(dir)
	this.world.Add(tank)
	pl.tankId = tank.Id

	if event {
		m := Event.New(tank)
		this.doSendEvent(m)
	}
}

func (this *Matrix) doJoinPlayers(plist []*ServiceChannel, teamId int) error {
	if plist != nil {
		for idx, sch := range plist {
			pl := new(Player)
			pl.sch = sch
			pl.teamId = teamId
			pl.teamNum = idx
			this.players[pl.Id()] = pl
			this.doAddTank(pl, false)
		}
	}
	return nil
}

func (this *Matrix) doReadyGo() error {

	cdown := MATRIX_BEGIN_COUNT
	var cdf func()
	cdf = func() {
		if this.executor.IsClosing() {
			return
		}
		this.doSendAll("INFO", fmt.Sprintf("%d seconds to begin", cdown))
		cdown--
		if cdown == 0 {
			this.doGo()
		} else {
			if this.timer != nil {
				this.timer.Stop()
			}
			this.timer = time.AfterFunc(1*time.Second, cdf)
		}
	}
	cdf()
	return nil
}

func (this *Matrix) doGo() error {
	this.doSendAll("BEGIN", "")
	if this.timer != nil {
		this.timer.Stop()
	}
	// send all teamInfo
	for _, pl := range this.players {
		ts := "A"
		if pl.teamId == 2 {
			ts = "B"
		}
		str := fmt.Sprintf("%s %d %s", ts, pl.tankId, pl.sch.name)
		this.doSendAll("TEAM", str)
	}
	// send snapshot
	this.doSendSnapshot()
	this.isBegin = true

	this.ltime = 0
	this.timer = time.AfterFunc(time.Duration(MATRIX_DURATION_MS)*time.Millisecond, this.doOneTurn)
	return nil
}

func (this *Matrix) doSendSnapshot() {
	w := this.world
	for _, obj := range w.SObj {
		m := Event.New(obj)
		this.doSendEvent(m)
	}
	for _, obj := range w.GObj {
		if obj.Kind == MCK_GROUP {
			m := Event.New(obj)
			this.doSendEvent(m)
		}
	}
	for _, obj := range w.MObj {
		m := Event.New(obj)
		this.doSendEvent(m)
	}
}

func (this *Matrix) doTankMove(tankId int, sp int) error {
	w := this.world
	obj := w.Objects[tankId]
	if obj == nil {
		return fmt.Errorf("tank(%d) not found", tankId)
	}
	if obj.Kind != MCK_TANK {
		return fmt.Errorf("object(%d) not a tank", tankId)
	}
	if sp != obj.Speed {
		obj.StartMove(sp)
		m := Event.ChangeSpeed(obj)
		this.doSendEvent(m)
	}
	return nil
}

func (this *Matrix) doTankTurn(tankId int, dir DIR) error {
	w := this.world
	obj := w.Objects[tankId]
	if obj == nil {
		return fmt.Errorf("tank(%d) not found", tankId)
	}
	if obj.Kind != MCK_TANK {
		return fmt.Errorf("object(%d) not a tank", tankId)
	}
	if dir != obj.Dir {
		obj.TurnDir(dir)
		m := Event.ChangeDir(obj)
		this.doSendEvent(m)
	}
	return nil
}

func (this *Matrix) doTankFire(tankId int) error {
	w := this.world
	obj := w.Objects[tankId]
	if obj == nil {
		return fmt.Errorf("tank(%d) not found", tankId)
	}
	if obj.Kind != MCK_TANK {
		return fmt.Errorf("object(%d) not a tank", tankId)
	}
	dir := obj.Dir
	if dir == DIR_NONE {
		dir = DIR_DOWN
	}
	var pos POS
	pos.x, pos.y = obj.Center()
	pos.x, pos.y = pos.Move(dir, MATRIX_BULLET_FIRE_RANGE)
	bul := newObject(pos.x-MATRIX_BULLET_SIZE/2, pos.y-MATRIX_BULLET_SIZE/2, MATRIX_BULLET_SIZE, MATRIX_BULLET_SIZE, MCK_BULLET)
	bul.Flag = obj.Flag
	w.Add(bul)
	bul.TurnDir(dir)
	bul.StartMove(MATRIX_BULLET_SPEED)
	ev := Event.New(bul)
	this.doSendEvent(ev)
	return nil
}

func (this *Matrix) doTankKillMe(tankId int) error {
	return nil
}

func (this *Matrix) doHit(bu *Object, co *Object) {
	logger.Debug(mtag, "%s hited %s", bu, co)
	switch co.Kind {
	case MCK_BASE:
		if this.winner == 0 {
			this.winner = bu.Flag
		} else {
			this.winner = -1
		}
		if this.events[co.Id] == nil {
			this.events[co.Id] = Event.Remove(co)
		}
	case MCK_WALL:
		return
	case MCK_BULLET:
		co.removed = true
	case MCK_ROCK:
		// bomb it
		px, py := bu.Center()
		var x1, y1, x2, y2 int32
		switch bu.Dir {
		case DIR_LEFT, DIR_RIGHT:
			x1 = px - MATRIX_BULLET_SIZE/2
			x2 = x1 + MATRIX_BULLET_SIZE
			y1 = py - MATRIX_UNIT/2
			y2 = y1 + MATRIX_UNIT
		case DIR_UP, DIR_DOWN:
			x1 = px - MATRIX_UNIT/2
			x2 = x1 + MATRIX_UNIT
			y1 = py - MATRIX_BULLET_SIZE/2
			y2 = y1 + MATRIX_BULLET_SIZE
		}
		m2 := make(map[*Object]bool)
		this.world.Collide(m2, x1, y1, x2, y2, nil, true)
		for ro, _ := range m2 {
			logger.Debug(mtag, "bomb %s", ro)
			ro.removed = true
			ev := Event.Remove(ro)
			this.events[ro.Id] = ev
		}
	case MCK_TANK:
		this.doDestroy(co)
	}
}

func (this *Matrix) doCrash(bu *Object, co *Object) {
	logger.Debug(mtag, "%s crashed %s", bu, co)
	switch co.Kind {
	case MCK_BULLET:
		co.removed = true
		this.doDestroy(bu)
	}
}

func (this *Matrix) doObjectMove(obj *Object) {
	if !obj.Moving() {
		return
	}
	w := this.world

	x, y := obj.Pos.Move(obj.Dir, int32(obj.Speed))
	// fmt.Println("try", x, y)

	// move it
	defer func() {
		for k, _ := range this.colmap {
			delete(this.colmap, k)
		}
	}()
	w.Collide(this.colmap, x, y, x+obj.Size.w, y+obj.Size.h, obj, false)
	if len(this.colmap) > 0 {
		// do coll
		switch obj.Kind {
		case MCK_BULLET:
			obj.removed = true
			obj.Pos.x = x
			obj.Pos.y = y
			for co, _ := range this.colmap {
				this.doHit(obj, co)
			}
		case MCK_TANK:
			done := obj.StopMove()
			if done && this.events[obj.Id] == nil {
				this.events[obj.Id] = Event.ChangeSpeed(obj)
			}
			for co, _ := range this.colmap {
				this.doCrash(obj, co)
			}
		}
		return
	}
	obj.Pos.x = x
	obj.Pos.y = y
	dmessage = fmt.Sprintf("%d - %v, %v", this.ltime, obj.Speed, obj.Dir)
}

func (this *Matrix) doOneTurn() {
	if this.executor.IsClosing() {
		return
	}
	this.ltime++
	if this.ltime >= MATRIX_LTIME_MAX {
		this.doEnd()
		return
	}
	this.timer.Reset(time.Duration(MATRIX_DURATION_MS) * time.Millisecond)

	// process
	w := this.world
	for _, obj := range this.world.MObj {
		this.doObjectMove(obj)
	}
	for _, obj := range this.world.MObj {
		if obj.removed {
			w.Remove(obj)
			ev := Event.Remove(obj)
			this.events[obj.Id] = ev
		}
	}
	for k, ev := range this.events {
		delete(this.events, k)
		this.doSendEvent(ev)
	}

	if this.winner != 0 {
		if this.winner == -1 {
			this.doEnd()
		} else {
			this.doWin(this.winner)
		}
	}

	if this.ltime%1000 == 0 {
		w.BuildDumpMap(this.dmap)
		logger.Debug(mtag, "\n%s", this.dmap.View())
		logger.Debug(mtag, "%d, %s", this.ltime, dmessage)
	}

}

var (
	dmessage string
)

func (this *Matrix) doDestroy(tank *Object) {
	logger.Info(mtag, "%s destroy", tank)
	tank.removed = true
}

func (this *Matrix) doEnd() {
	logger.Info(mtag, "%s end", this)
	this.doSendAll("END", "")
	go this.AskClose()
}

func (this *Matrix) doWin(flag int) {
	logger.Info(mtag, "%s end - winner is %d", this, flag)
	ts := "A"
	if flag == 2 {
		ts = "B"
	}
	this.doSendAll("END", ts)
	go this.AskClose()
}
