package goo

const (
	STATE_INIT  = uint32(0)
	STATE_START = uint32(1)
	STATE_RUN   = uint32(2)
	STATE_STOP  = uint32(3)
	STATE_CLOSE = uint32(4)
	STATE_END   = uint32(5)
	STATE_IDLE  = uint32(6)
)

var (
	INFO_STATE_INIT  = NewStateInfO(STATE_INIT, "init")
	INFO_STATE_START = NewStateInfO(STATE_START, "start")
	INFO_STATE_RUN   = NewStateInfO(STATE_RUN, "run")
	INFO_STATE_STOP  = NewStateInfO(STATE_STOP, "stop")
	INFO_STATE_CLOSE = NewStateInfO(STATE_CLOSE, "close")
	INFO_STATE_END   = NewStateInfO(STATE_END, "end")
	INFO_STATE_IDLE  = NewStateInfO(STATE_IDLE, "idle")
)
