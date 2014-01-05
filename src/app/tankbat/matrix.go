package tankbat

import (
	"bmautil/qexec"
	"encoding/json"
	"fmt"
	"logger"
	"time"
)

var (
	globalMID int
)

const (
	MATRIX_BASE_UNIT    = int32(10000)
	MATRIX_UNIT         = 4 * MATRIX_BASE_UNIT
	MATRIX_MAP_WIDTH    = 13
	MATRIX_WORLD_WIDTH  = MATRIX_MAP_WIDTH * MATRIX_UNIT
	MATRIX_MAP_HEIGHT   = 11
	MATRIX_WORLD_HEIGHT = MATRIX_MAP_HEIGHT * MATRIX_UNIT
	MATRIX_DURATION_MS  = 1
	MATRIX_SEC_COUNT    = 1000 / MATRIX_DURATION_MS
	MATRIX_DURATION_MAX = 60 * 1000
	MATRIX_LTIME_MAX    = MATRIX_DURATION_MAX / MATRIX_DURATION_MS
	MATRIX_BEGIN_COUNT  = 5

	MATRIX_ROCK_SIZE         = MATRIX_BASE_UNIT
	MATRIX_TANK_SIZE         = MATRIX_UNIT
	MATRIX_TANK_SPEED        = 10
	MATRIX_BULLET_SIZE       = 1 * MATRIX_BASE_UNIT
	MATRIX_BULLET_SPEED      = 20
	MATRIX_BULLET_FIRE_RANGE = MATRIX_TANK_SIZE + MATRIX_BASE_UNIT/2
)

// Matrix
type Matrix struct {
	mid      int
	seqid    int
	eventId  int
	executor qexec.QueueExecutor
	timer    *time.Timer
	players  map[uint32]*Player
	watchers map[string]*ServiceChannel
	world    *World
	ltime    int

	service *Service

	isBegin bool
	colmap  map[*Object]bool
	events  map[int]interface{}
	winner  int
	dmap    *dumpMap
	builder *mapBuilder
}

func NewMatrix(s *Service, bz int) *Matrix {
	this := new(Matrix)
	globalMID++
	this.mid = globalMID
	this.service = s
	this.executor.InitQueueExecutor(mtag, bz, this.requestHandler)
	this.executor.StopHandler = this.stopHandler

	this.players = make(map[uint32]*Player)
	this.watchers = make(map[string]*ServiceChannel)
	this.colmap = make(map[*Object]bool)
	this.events = make(map[int]interface{})

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
	if this.timer != nil {
		this.timer.Stop()
	}
	s := this.service
	this.service = nil
	if s != nil && !s.executor.IsClosing() {
		s.OnMatrixEnd(this)
	}
}

func (this *Matrix) Run(mapId int, teamA []*ServiceChannel, teamB []*ServiceChannel) bool {
	globalMID++
	this.mid = globalMID
	this.executor.Run()
	this.executor.DoNow("init", func() error {
		logger.Info(mtag, "%s run", this)
		if err := this.doInitMatrix(); err != nil {
			return err
		}
		this.doJoinPlayers(teamA, 1)
		this.doJoinPlayers(teamB, 2)
		if err := this.doReadyGo(); err != nil {
			return err
		}
		return nil
	})
	return true
}

func (this *Matrix) doFindTank(sch *ServiceChannel) int {
	for _, pl := range this.players {
		if pl.Id() == sch.Id() {
			return pl.tankId
		}
	}
	return 0
}

func (this *Matrix) CommandMove(sch *ServiceChannel) {
	err := this.executor.DoNow("cmdMove", func() error {
		if !this.isBegin {
			sch.ReplyOK()
			return nil
		}
		tid := this.doFindTank(sch)
		if tid == 0 {
			sch.BeError(fmt.Errorf("not found player tank"))
		} else {
			err := this.doTankMove(tid, MATRIX_TANK_SPEED)
			if err != nil {
				sch.BeError(err)
			} else {
				sch.ReplyOK()
			}
		}
		return nil
	})
	if err != nil {
		sch.BeError(err)
	}
}

func (this *Matrix) CommandStop(sch *ServiceChannel) {
	err := this.executor.DoNow("cmdStop", func() error {
		if !this.isBegin {
			sch.ReplyOK()
			return nil
		}
		tid := this.doFindTank(sch)
		if tid == 0 {
			sch.BeError(fmt.Errorf("not found player tank"))
		} else {
			err := this.doTankMove(tid, 0)
			if err != nil {
				sch.BeError(err)
			} else {
				sch.ReplyOK()
			}
		}
		return nil
	})
	if err != nil {
		sch.BeError(err)
	}
}

func (this *Matrix) CommandTurn(sch *ServiceChannel, dir DIR) {
	err := this.executor.DoNow("cmdTurn", func() error {
		if !this.isBegin {
			sch.ReplyOK()
			return nil
		}
		tid := this.doFindTank(sch)
		if tid == 0 {
			sch.BeError(fmt.Errorf("not found player tank"))
		} else {
			err := this.doTankTurn(tid, dir)
			if err != nil {
				sch.BeError(err)
			} else {
				sch.ReplyOK()
			}
		}
		return nil
	})
	if err != nil {
		sch.BeError(err)
	}
}

func (this *Matrix) CommandFire(sch *ServiceChannel) {
	err := this.executor.DoNow("cmdTurn", func() error {
		if !this.isBegin {
			sch.ReplyOK()
			return nil
		}
		tid := this.doFindTank(sch)
		if tid == 0 {
			sch.BeError(fmt.Errorf("not found player tank"))
		} else {
			err := this.doTankFire(tid)
			if err != nil {
				sch.BeError(err)
			} else {
				sch.ReplyOK()
			}
		}
		return nil
	})
	if err != nil {
		sch.BeError(err)
	}
}

func (this *Matrix) CommandKillMe(sch *ServiceChannel) {
	err := this.executor.DoNow("cmdTurn", func() error {
		if !this.isBegin {
			sch.ReplyOK()
			return nil
		}
		tid := this.doFindTank(sch)
		if tid == 0 {
			sch.BeError(fmt.Errorf("not found player tank"))
		} else {
			err := this.doTankKillMe(tid)
			if err != nil {
				sch.BeError(err)
			} else {
				sch.ReplyOK()
			}
		}
		return nil
	})
	if err != nil {
		sch.BeError(err)
	}
}

func (this *Matrix) IsClosing() bool {
	return this.executor.IsClosing()
}

func (this *Matrix) AskClose() bool {
	this.executor.Stop()
	return true
}

func (this *Matrix) doSendEvent(info interface{}) {
	bs, err := json.Marshal(info)
	if err != nil {
		logger.Error(mtag, "doSendEvent(%v) fail - %s", info, err)
		return
	}
	this.eventId++
	s := fmt.Sprintf("EVENT %d %s\n", this.eventId, string(bs))
	this.doSendAll0(s)
}

func (this *Matrix) doSendAll(cmd string, str string) {
	s := ""
	if str != "" {
		s = fmt.Sprintf("%s %s\n", cmd, str)
	} else {
		s = fmt.Sprintf("%s\n", cmd)
	}
	this.doSendAll0(s)
}

func (this *Matrix) doSendAll0(s string) {
	logger.Debug(mtag, "sendAll ==> %s", s)
	for _, pl := range this.players {
		if pl.sch != nil {
			err := pl.sch.Send(s)
			if err != nil {
				logger.Debug(mtag, "send %s fail - %s", pl.sch, err)
			}
		}
	}
	for _, w := range this.watchers {
		err := w.Send(s)
		if err != nil {
			logger.Debug(mtag, "send %s fail - %s", w, err)
		}
	}
}

func (this *Matrix) doSend(sch *ServiceChannel, cmd string, info interface{}) {
	bs, err := json.Marshal(info)
	if err != nil {
		logger.Error(mtag, "doSend(%s, %v) fail - %s", cmd, info, err)
		return
	}
	s := fmt.Sprintf("%s %s\n", cmd, string(bs))
	logger.Debug(mtag, "doSend %s ==> %s", sch, s)
	err = sch.Send(s)
	if err != nil {
		logger.Debug(mtag, "send %s fail - %s", sch, err)
	}
}
