package acclog

type AccLogInfo interface {
	Message(cfg map[string]string) string
	TimeDay() int
}
