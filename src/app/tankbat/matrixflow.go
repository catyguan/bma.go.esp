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
	wox := ww / 2
	woy := wh / 2
	this.dmap = newDumpMap(ww, wh, MATRIX_UNIT/2, MATRIX_UNIT/2)

	w := NewWorld()
	this.world = w
	mxf := func(x int32) int32 {
		return (x-10)*MATRIX_UNIT/10 + MATRIX_UNIT - wox
	}
	myf := func(y int32) int32 {
		return (y-10)*MATRIX_UNIT/10 + MATRIX_UNIT - woy
	}
	// layout
	w.Add(newObject(mxf(0), myf(0), ww, MATRIX_UNIT, MCK_WALL))
	w.Add(newObject(mxf(0), myf((MATRIX_MAP_HEIGHT+2-1)*10), ww, MATRIX_UNIT, MCK_WALL))

	w.Add(newObject(mxf(0), myf(0), MATRIX_UNIT, wh-1*MATRIX_UNIT, MCK_WALL))
	w.Add(newObject(mxf((MATRIX_MAP_WIDTH+2-1)*10), myf(0), MATRIX_UNIT, wh-1*MATRIX_UNIT, MCK_WALL))

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
		w.Add(newObject(mxf(info.x), myf(info.y), info.w*MATRIX_UNIT/10, info.h*MATRIX_UNIT/10, MCK_WALL))
	}

	rocks := []*RECT{

		newRECT(65, 10, 5, 10),
		newRECT(80, 10, 5, 10),
		newRECT(65, 20, 20, 5),

		newRECT(20, 20, 10, 35),
		newRECT(40, 20, 10, 35),
		newRECT(100, 20, 10, 35),
		newRECT(120, 20, 10, 35),

		// 	newRECT(60, 30, 10, 30),
		// 	newRECT(80, 30, 10, 30),

		// 	// center begin
		// 	newRECT(60, 70, 10, 10),
		// 	newRECT(80, 70, 10, 10),
		// 	// center end
		// 	newRECT(60, 90, 10, 30),
		// 	newRECT(80, 90, 10, 30),

		// 	newRECT(20, 95, 10, 35),
		// 	newRECT(40, 95, 10, 35),
		// 	newRECT(100, 95, 10, 35),
		// 	newRECT(120, 95, 10, 35),

		// 	newRECT(65, 125, 20, 5),
		// 	newRECT(65, 130, 5, 10),
		// 	newRECT(80, 130, 5, 10),
	}
	for _, info := range rocks {
		w.Add(newObject4Rock(mxf(info.x), myf(info.y), info.w*MATRIX_UNIT/10, info.h*MATRIX_UNIT/10, MATRIX_ROCK_SIZE, MATRIX_ROCK_SIZE))
	}

	// w.AddBaseToMap((70-10)*4/10+2, (10-10)*4/10+2, 4, 4, 1)
	base := newObject(mxf(70), myf(10), 10*MATRIX_UNIT/10, 10*MATRIX_UNIT/10, MCK_BASE)
	base.Flag = 1
	w.Add(base)

	// test
	tank := newObject(mxf(10), myf(10), MATRIX_UNIT, MATRIX_UNIT, MCK_TANK)
	tank.Flag = 1
	w.Add(tank)
	if tank != nil {
		tankId = tank.Id
		tank.TurnDir(DIR_RIGHT)
		tank.StartMove(10)
	}
	w.BuildDumpMap(this.dmap)
	fmt.Print(this.dmap.View())
	return nil
}

func (this *Matrix) doInitPlayers(plist []*ServiceChannel) error {
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
	this.doSendAll("INFO", "begin")
	if this.timer != nil {
		this.timer.Stop()
	}
	this.ltime = 0
	this.timer = time.AfterFunc(time.Duration(MATRIX_DURATION_MS)*time.Millisecond, this.doOneTurn)
	return nil
}

func (this *Matrix) doFire(tankId int) error {
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
	bullet = bul
	return nil
}

func (this *Matrix) doHit(bu *Object, co *Object) {
	logger.Debug(mtag, "%s hited %s", bu, co)
	switch co.Kind {
	case MCK_BASE:
		this.doWin(bu.Flag)
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
			obj.StopMove()
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
	defer func() {
		this.timer.Reset(time.Duration(MATRIX_DURATION_MS) * time.Millisecond)
	}()
	// process
	w := this.world
	for _, obj := range this.world.MObj {
		this.doObjectMove(obj)
	}
	for _, obj := range this.world.MObj {
		if obj.removed {
			w.Remove(obj)
		}
	}

	if this.ltime%1000 == 0 {
		w.BuildDumpMap(this.dmap)
		logger.Debug(mtag, "\n%s", this.dmap.View())
		logger.Debug(mtag, "%d, %v, %s", this.ltime, bullet, dmessage)
	}

	if this.ltime > 1000 && this.ltime%2000 == 100 {
		this.doFire(tankId)
	}
}

var (
	bullet   *Object
	tankId   int
	dmessage string
)

func (this *Matrix) doDestroy(tank *Object) {
	logger.Info(mtag, "%s destroy", tank)
	tank.removed = true
}

func (this *Matrix) doEnd() {
	logger.Info(mtag, "%s end", this)
	go this.AskClose()
}

func (this *Matrix) doWin(flag int) {
	logger.Info(mtag, "%s end - winner is %d", this, flag)
	go this.AskClose()
}
