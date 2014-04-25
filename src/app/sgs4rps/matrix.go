package sgs4rps

import (
	"fmt"
	"logger"
)

type matrix struct {
	s       *Service
	psid    int
	players map[int]*Player

	championId int
}

func (this *matrix) start() {
	this.psid = 0
	this.players = make(map[int]*Player)
	this.championId = 0

	c := this.s.config.RobotNum
	for i := 0; i < c; i++ {
		// this.createRobot(i + 1)
	}
}

func (this *matrix) stop() {
	for k, _ := range this.players {
		delete(this.players, k)
	}
}

func (this *matrix) clear() {

}

func (this *matrix) createRobot(idx int) {
	go func() {
		psid, err := this.s.NewPlayer(nil)
		if err != nil {
			logger.Error(tag, "createRobot NewPlayer fail")
			return
		}
		err = this.s.SetNick(psid, fmt.Sprintf("robot%d", idx))
		if err != nil {
			logger.Error(tag, "createRobot SetNick fail")
			return
		}
	}()
}
