package bombman

import (
	"bmautil/qexec"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"logger"
	"math/rand"
	"sort"
	"time"
)

var (
	globalMID int
)

const (
	ACTION_NONE = 0
	ACTION_MOVE = 1
	ACTION_BOMB = 2
)
const (
	DIR_NONE  = 0
	DIR_LEFT  = 4
	DIR_RIGHT = 6
	DIR_UP    = 8
	DIR_DOWN  = 2
)
const (
	MATRIX_MAX_SEQ    = 10 * 60
	MATRIX_DURATION   = 100
	MATRIX_MOVE_DELAY = 1
	MATRIX_BOMB_TIME  = (MATRIX_BOMB_RANGE + 1) * MATRIX_MOVE_DELAY
	MATRIX_BOMB_RANGE = 3
)

// Matrix
type Matrix struct {
	mid      int
	seqid    int
	executor qexec.QueueExecutor
	timer    *time.Timer
	players  playerList
	watchers map[string]ServiceChanel
	mapdata  Map
	bombList *list.List

	service *Service
}

func NewMatrix(bz int) *Matrix {
	this := new(Matrix)
	this.executor.InitQueueExecutor(mtag, bz, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.bombList = list.New()
	return this
}

func (this *Matrix) String() string {
	return fmt.Sprintf("Matrix[%d]", this.mid)
}

func (this *Matrix) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	case func(*Matrix) error:
		return true, rv(this)
	}
	return true, nil
}

func (this *Matrix) stopHandler() {
	logger.Info(mtag, "%s stop", this)
}

func (this *Matrix) Run(initor func(m *Matrix) error) bool {
	globalMID++
	this.mid = globalMID
	this.executor.Run()
	this.executor.DoNow("init", func() error {
		logger.Info(mtag, "%s run", this)
		rand.Seed(time.Now().UnixNano())
		this.seqid = 1
		err := initor(this)
		if err != nil {
			return err
		}
		this.timer = time.AfterFunc(time.Duration(MATRIX_DURATION)*time.Millisecond, this.timerPost)
		return nil
	})
	return true
}

func (this *Matrix) timerPost() {
	this.executor.DoNow("timer", func() error {
		err := this.processMatrix()
		if err != nil {
			return err
		}
		this.timer.Reset(time.Duration(MATRIX_DURATION) * time.Millisecond)
		return nil
	})
}

func (this *Matrix) AskClose() bool {
	this.executor.Stop()
	return true
}

func (this *Matrix) DoInit(players int, mapw, maph int) {
	logger.Info(mtag, "DoInit(players=%d, mapw=%d, maph=%d)", players, mapw, maph)
	this.mapdata.InitMap(mapw, maph)
	if players > 4 {
		players = 4
	}
	this.players = make(playerList, players)
	for i := 0; i < players; i++ {
		pid := i + 1
		this.players[i] = new(Player)
		this.players[i].id = pid
		this.players[i].name = fmt.Sprintf("Player%d", pid)
		var x, y int
		switch pid {
		case 1:
			x = 0
			y = 0
		case 2:
			x = mapw - 1
			y = maph - 1
		case 3:
			x = mapw - 1
			y = 0
		case 4:
			x = 0
			y = maph - 1
		}
		this.mapdata.PlayerIn(x, y, this.players[i])
	}
}

func (this *Matrix) DoAttachPlayer(pid int, sch *ServiceChanel) {
	idx := pid - 1
	if idx >= 0 && idx < len(this.players) {
		this.players[idx].channel = sch
	}
	this.doSendSnapshot(sch)
}

// func (this *Matrix) DoAddWatcher(w MatrixWatcher) {

// }

func (this *Matrix) DoView() string {
	return this.mapdata.View()
}

func (this *Matrix) doSendAll(cmd string, info interface{}) {
	bs, err := json.Marshal(info)
	if err != nil {
		logger.Error(mtag, "sendAll(%s, %v) fail - %s", cmd, info, err)
		return
	}
	s := fmt.Sprintf("%s %s\n", cmd, string(bs))
	logger.Debug(mtag, "sendAll ==> %s", s)
	for _, pl := range this.players {
		if pl.channel != nil {
			pl.channel.Send(s)
		}
	}
	for _, w := range this.watchers {
		w.Send(s)
	}
}

func (this *Matrix) doSend(sch *ServiceChanel, cmd string, info interface{}) {
	bs, err := json.Marshal(info)
	if err != nil {
		logger.Error(mtag, "doSend(%s, %v) fail - %s", cmd, info, err)
		return
	}
	s := fmt.Sprintf("%s %s\n", cmd, string(bs))
	logger.Debug(mtag, "doSend %s ==> %s", sch.channel, s)
	sch.Send(s)
}

func (this *Matrix) doSendSnapshot(sch *ServiceChanel) {
	ss := new(pSnapshot)
	ss.sid = this.seqid
	ss.mapcell = this.mapdata.MakeSnapshot(this.seqid)
	info := make(map[string]interface{})
	ss.Build(info)
	go this.doSend(sch, "SNAPSHOT", info)
}

func (this *Matrix) playerFrom(sch *ServiceChanel) *Player {
	for _, p := range this.players {
		if p.channel != nil && p.channel.channel == sch.channel {
			return p
		}
	}
	return nil
}

func (this *Matrix) PostAction(sch *ServiceChanel, action int, dir int) {
	pl := this.playerFrom(sch)
	if pl == nil {
		sch.BeError(logger.Error(mtag, "invalid %s player action", sch.channel))
		return
	}
	if pl.died {
		sch.BeError(errors.New("you are die"))
		return
	}
	pl.actionTime = time.Now().UnixNano()
	pl.action = action
	pl.actionDir = dir
	sch.Replay("OK\n")
}

func (this *Matrix) processMatrix() error {
	this.seqid++
	if this.seqid >= MATRIX_MAX_SEQ {
		logger.Info(mtag, "%s end", this)
		this.doSendAll("END", "")
		this.service.OnMatrixEnd(this)
		return nil
	}
	// logger.Debug(mtag, "%s process %d", this.seqid)
	sort.Sort(this.players)
	liveCount := 0
	for _, pl := range this.players {
		if pl.died {
			continue
		}
		liveCount++
		switch pl.action {
		case ACTION_MOVE:
			old := pl.MapPos
			if this.mapdata.PlayerMove(pl, pl.actionDir, this.seqid) {
				if logger.EnableDebug(mtag) {
					logger.Debug(mtag, "%s move %s ~> %s", pl, old.PosString(), pl.MapPos.PosString())
				}
			}
		case ACTION_BOMB:
			bomb := new(bombInfo)
			bomb.seqid = this.seqid
			if bpos := this.mapdata.PlayerBomb(pl, pl.actionDir, bomb); bpos != nil {
				bomb.MapPos = *bpos
				bomb.who = pl.id
				this.bombList.PushBack(bomb)
				if logger.EnableDebug(mtag) {
					logger.Debug(mtag, "%s bomb %s ~> %s", pl, pl.MapPos.PosString(), bpos.PosString())
				}
			}
		}
		pl.action = ACTION_NONE
		pl.actionDir = DIR_NONE
	}

	if liveCount < 2 {
		logger.Info(mtag, "%s end", this)
		this.doSendAll("END", "")
		this.service.OnMatrixEnd(this)
		return nil
	}

	for e := this.bombList.Front(); e != nil; e = e.Next() {
		binfo := e.Value.(*bombInfo)
		if binfo.fired {
			continue
		}
		if this.seqid-binfo.seqid >= MATRIX_BOMB_TIME {
			// bomb!!
			binfo.fired = true
			this.doBombFire(binfo)
		}
	}

	es := this.mapdata.MakeEvent(this.seqid)
	for _, pmc := range es {
		info := make(map[string]interface{})
		info["sid"] = this.seqid
		pmc.Build(info)
		this.doSendAll("INFO", info)
	}
	return nil
}

var (
	dirs []int = []int{DIR_UP, DIR_RIGHT, DIR_DOWN, DIR_LEFT}
)

func (this *Matrix) doBombFire(binfo *bombInfo) {
	if logger.EnableDebug(mtag) {
		logger.Debug(mtag, "doBombFire(%s)", binfo.PosString())
	}
	cell := this.mapdata.Cell(binfo.x, binfo.y)
	if cell != nil {
		cell.seqid = this.seqid
		cell.event = EVENT_BOMB
		cell.kind = MCK_NONE
	}

	for _, dir := range dirs {
		p1 := binfo.MapPos
		for i := 0; i < MATRIX_BOMB_RANGE; i++ {
			p2 := p1.Dir(dir)
			c2 := this.mapdata.Cell(p2.x, p2.y)
			if c2 == nil {
				break
			}
			cdo := true
			switch c2.kind {
			case MCK_WALL:
				cdo = false
			case MCK_ROCK:
				c2.seqid = this.seqid
				c2.event = EVENT_BOMB
				cdo = false
			case MCK_BOMB:
				if c2.bomb != nil && !c2.bomb.fired {
					c2.bomb.fired = true
					this.doBombFire(c2.bomb)
				}
			case MCK_NONE:
				c2.seqid = this.seqid
				c2.event = EVENT_BOMB
			case MCK_PLAYER:
				c2.seqid = this.seqid
				c2.event = EVENT_BOMB
				if c2.player != nil {
					c2.player.died = true
					this.mapdata.PlayerDie(c2.player)
				}
			}
			if !cdo {
				break
			}
			p1 = p2
		}
	}
}
