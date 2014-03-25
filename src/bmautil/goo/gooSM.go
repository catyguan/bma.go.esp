package goo

import "logger"

var gooStates = StateCollection{
	STATE_INIT:  STR_STATE_INIT,
	STATE_START: STR_STATE_START,
	STATE_RUN:   STR_STATE_RUN,
	STATE_STOP:  STR_STATE_STOP,
	STATE_CLOSE: STR_STATE_CLOSE,
}

func canEnter4goo(o interface{}, cur uint32, st uint32) bool {
	switch cur {
	case STATE_INIT:
		return st == STATE_START || st == STATE_CLOSE
	case STATE_START:
		return st == STATE_RUN || st == STATE_STOP
	case STATE_RUN:
		return st == STATE_STOP
	case STATE_STOP:
		return st == STATE_CLOSE
	}
	return false
}

func afterEnter4goo(o interface{}, st uint32) {
	obj := o.(*Goo)
	if obj.EDebug {
		logger.Debug(obj.Tag, "enterState(%s)", gooStates.ToString(st))
	}
	switch st {
	case STATE_STOP:
		obj.queue <- nil
	case STATE_START:
		go obj.run()
	}
}
