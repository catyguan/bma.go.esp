package tankbat

import (
	"fmt"
	"logger"
	"sort"
)

func (this *Service) doCheckWaiting() {
	if this.matrix != nil {
		// running
		logger.Debug(tag, "matrix running, skip check waiting")
		return
	}
	c := 0
	for _, sch := range this.channels {
		if sch.waiting {
			c++
		}
	}
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

	pl := make([]*ServiceChannel, 0)
	for _, sch := range this.channels {
		if sch.waiting {
			pl = append(pl, sch)
		}
	}
	if len(pl) < this.config.GamePlayerMin {
		return nil
	}
	if len(pl) > this.config.GamePlayerMax {
		sort.Sort(ByPlayTime{pl})
		pl = pl[:this.config.GamePlayerMin]
	}

	m := NewMatrix(this, 32)
	this.matrix = m
	m.Run()

	// for _, sch := range pl {
	// 	func(sch *ServiceChannel) {
	// 		m.executor.DoNow("addPlayer", func(m *Matrix) error {
	// 			m.DoAttachPlayer(sch)
	// 			return nil
	// 		})
	// 	}(sch)
	// }

	// m.executor.DoNow("go", func(m *Matrix) error {
	// 	m.DoGo()
	// 	return nil
	// })
	return nil
}
