package goo

import "logger"

var gooStates = []*StateInfo{
	INFO_STATE_INIT,
	INFO_STATE_START,
	INFO_STATE_RUN,
	INFO_STATE_STOP,
	INFO_STATE_CLOSE,
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
		logger.Debug(obj.Tag, "enterState(%s)", st)
	}
	switch st {
	case STATE_STOP:
		close(obj.queue)
	case STATE_START:
		go obj.run()
	}
}
