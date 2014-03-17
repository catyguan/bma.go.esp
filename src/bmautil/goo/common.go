package goo

const (
	STATE_INVALID = uint32(0)
	STATE_INIT    = uint32(1)
	STATE_START   = uint32(2)
	STATE_RUN     = uint32(3)
	STATE_STOP    = uint32(4)
	STATE_CLOSE   = uint32(5)
	STATE_END     = uint32(6)
	STATE_IDLE    = uint32(7)
)

const (
	STR_STATE_INVALID = "invalid"
	STR_STATE_INIT    = "init"
	STR_STATE_START   = "start"
	STR_STATE_RUN     = "run"
	STR_STATE_STOP    = "stop"
	STR_STATE_CLOSE   = "close"
	STR_STATE_END     = "end"
	STR_STATE_IDLE    = "idle"
)
