package tankbat

import (
	"bmautil/valutil"
	"errors"
	"esp/espnet"
	"fmt"
	"logger"
	"strings"
	"time"
)

func (this *Service) doChannelMessage(sch *ServiceChannel, msg *espnet.Message) {
	if !sch.cmdTime.IsZero() {
		if time.Now().Sub(sch.cmdTime) < COMMAND_LIMIT_MS*time.Millisecond {
			return
		}
	}
	sch.cmdTime = time.Now()
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
	case "v":
		if this.matrix == nil {
			err := errors.New("matrix nil")
			sch.BeError(err)
			return
		}
		this.matrix.executor.DoNow("view", func(m *Matrix) error {
			m.world.BuildDumpMap(m.dmap)
			s := m.dmap.View()
			sch.Reply(s)
			return nil
		})
		return
	}

	switch cmd {
	case "join":
		if sch.name != "" {
			sch.BeError(fmt.Errorf("has join"))
			return
		}
		sp := strings.Split(params, " ")
		nick := sp[0]
		if nick == "" {
			sch.BeError(fmt.Errorf("join miss name"))
			return
		}
		teamId := 0
		if len(sp) > 1 {
			switch sp[1] {
			case "a":
				teamId = 1
			case "b":
				teamId = 2
			default:
				sch.BeError(fmt.Errorf("unknow teamId %s", sp[1]))
				return
			}
		}

		for _, c := range this.channels {
			if c.name == nick {
				sch.BeError(fmt.Errorf("name '%s' exists", nick))
				return
			}
		}

		sch.ReplyOK()

		sch.name = nick
		sch.joinTeamId = teamId
		this.doJoin(sch)
		this.doCheckWaiting()
		return
	}

	if sch.playing {
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
				v := DIR(valutil.ToInt(s, 0))
				if v == DIR_LEFT || v == DIR_RIGHT || v == DIR_UP || v == DIR_DOWN {
					return v
				}
				return DIR_NONE
			}
		}

		switch cmd {
		case "stop":
			this.matrix.CommandStop(sch)
			return
		case "move":
			this.matrix.CommandMove(sch)
			return
		case "turn":
			dir := dirf(params)
			if dir == DIR_NONE {
				sch.BeError(fmt.Errorf("turn where?"))
				return
			}
			this.matrix.CommandTurn(sch, dir)
			return
		case "fire":
			this.matrix.CommandFire(sch)
			return
		case "killme":
			this.matrix.CommandKillMe(sch)
			return
		}
	}
	logger.Info(mtag, "unknow command - %s %s", cmd, params)
	sch.BeError(fmt.Errorf("unknow command '%s'", cmd))
}
