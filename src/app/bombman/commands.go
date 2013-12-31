package bombman

import (
	"errors"
	"esp/espnet"
	"logger"
	"strings"
)

func (this *Service) doChannelMessage(sch *ServiceChanel, msg *espnet.Message) {
	b := msg.GetPayloadB()
	if b == nil {
		return
	}
	line := strings.ToLower(strings.TrimSpace(string(b)))
	logger.Debug(tag, "%s >> %s", sch.channel, line)

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

func (this *Service) processCommand(sch *ServiceChanel, cmd, params string) {
	// only for test
	switch cmd {
	case "start":
		if this.matrix != nil {
			err := errors.New("matrix started")
			sch.BeError(err)
			return
		}
		m := NewMatrix(32)
		m.service = this
		m.Run(this.initMatrix)
		this.matrix = m
		sch.Replay("OK\n")

		m.executor.DoNow("addPlayer", func(m *Matrix) error {
			m.DoAttachPlayer(1, sch)
			return nil
		})
		return
	case "v":
		if this.matrix == nil {
			err := errors.New("matrix nil")
			sch.BeError(err)
			return
		}
		this.matrix.executor.DoNow("view", func(m *Matrix) error {
			s := m.DoView()
			sch.Replay(s)
			return nil
		})
		return
	}

	switch cmd {
	case "join":
	}

	if this.matrix == nil {
		err := errors.New("matrix nil")
		sch.BeError(err)
		return
	}
	dirf := func(s string) int {
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
			this.matrix.PostAction(sch, ACTION_MOVE, dir)
			return
		}
	case "bomb":
		dir := dirf(params)
		if dir != DIR_NONE {
			this.matrix.PostAction(sch, ACTION_BOMB, dir)
			return
		}
	}
	logger.Info(mtag, "unknow command - %s %s", cmd, params)
	sch.BeError(errors.New("unknow command"))
}

func (this *Service) initMatrix(m *Matrix) error {
	m.DoInit(4, 11, 11)
	return nil
}
