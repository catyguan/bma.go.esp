package tankbat

import (
	"errors"
	"esp/espnet"
	"fmt"
	"logger"
	"strings"
	"time"
)

func (this *Service) doChannelMessage(sch *ServiceChannel, msg *espnet.Message) {
	b := msg.GetPayloadB()
	if b == nil {
		return
	}
	line := strings.ToLower(strings.TrimSpace(string(b)))
	logger.Debug(tag, "%s >> %s", sch, line)

	larray := strings.SplitN(line, " ", 2)
	cmd := larray[0]
	params := ""
	if len(larray) > 1 {
		params = larray[1]
	}
	if cmd == "" {
		return
	}
	this.processCommand(sch, cmd, params)
}

func (this *Service) processCommand(sch *ServiceChannel, cmd, params string) {
	// only for test
	switch cmd {
	case "start":
		this.doStartGame()
		return
	case "v":
		if this.matrix == nil {
			err := errors.New("matrix nil")
			sch.BeError(err)
			return
		}
		this.matrix.executor.DoNow("view", func(m *Matrix) error {
			s := ""
			sch.Reply(s)
			return nil
		})
		return
	}

	switch cmd {
	case "join":
		if params == "" {
			err := fmt.Errorf("join miss name")
			sch.BeError(err)
			return
		}
		sch.ReplyOK()

		if sch.name == "" {
			sch.name = params
			sch.waiting = true
			sch.waitTime = time.Now()
			this.doCheckWaiting()
		}
		return
	}

	if this.matrix == nil {
		err := errors.New("matrix nil")
		sch.BeError(err)
		return
	}
	dirf := func(s string) DIR {
		switch s {
		case "left":
			return DIR_LEFT
		case "right":
			return DIR_RIGHT
		case "up":
			return DIR_UP
		case "down":
			return DIR_DOWN
		default:
			return DIR_NONE
		}
	}

	switch cmd {
	case "move":
		dir := dirf(params)
		if dir != DIR_NONE {
			return
		}
	case "bomb":
		dir := dirf(params)
		if dir != DIR_NONE {
			return
		}
	}
	logger.Info(mtag, "unknow command - %s %s", cmd, params)
	sch.BeError(errors.New("unknow command"))
}
