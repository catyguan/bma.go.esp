package tankbat

import (
	"fmt"
	"logger"
	"sort"
	"time"
)

func (this *Service) doJoin(sch *ServiceChannel) {
	wui := new(WUserInfo)
	wui.channel = sch
	wui.waitTime = time.Now()
	this.waitingRoom[sch.Id()] = wui
}

func (this *Service) doCheckWaiting() {
	if this.matrix != nil {
		// running
		logger.Debug(tag, "matrix running, skip check waiting")
		return
	}
	c := len(this.waitingRoom)
	logger.Debug(tag, "waiting user(%d)", c)
	if c < this.config.GamePlayerMin {
		return
	}

	logger.Debug(tag, "game starting")
	this.doStartGame()
}

func (this *Service) doStartGame() error {
	if this.matrix != nil {
		logger.Debug(tag, "matrix running, skip start game")
		return fmt.Errorf("matrix running")
	}

	pl := make(WUserInfoList, 0)
	for _, sch := range this.waitingRoom {
		pl = append(pl, sch)
	}
	if len(pl) < this.config.GamePlayerMin {
		return nil
	}
	sort.Sort(pl)
	tnum := len(pl)/2 + len(pl)%2
	if tnum > this.config.TeamPlayerMax {
		tnum = this.config.TeamPlayerMax
	}
	if len(pl) > tnum*2 {
		pl = pl[:tnum*2]
	}
	ta := make([]*ServiceChannel, 0)
	tb := make([]*ServiceChannel, 0)
	for _, p := range pl {
		k := p.channel.joinTeamId
		switch k {
		case 0:
			if len(ta) > len(tb) {
				k = 2
			} else {
				k = 1
			}
		case 1:
			if len(ta) >= tnum {
				k = 2
			}
		case 2:
			if len(tb) >= tnum {
				k = 1
			}
		}
		if k == 1 {
			ta = append(ta, p.channel)
		} else {
			tb = append(tb, p.channel)
		}
		p.channel.playing = true
		delete(this.waitingRoom, p.channel.Id())
	}

	m := NewMatrix(this, 32)
	this.matrix = m
	m.Run(1, ta, tb)
	return nil
}

func (this *Service) doMatrixEnd(m *Matrix) {
	if this.matrix != m {
		return
	}
	this.matrix = nil
	for _, pl := range m.players {
		pl.sch.playing = false
		if _, ok := this.channels[pl.Id()]; ok {
			this.doJoin(pl.sch)
		}
	}
	this.doCheckWaiting()
}
