package boot

import "time"

var (
	DevMode         bool
	Debug           bool
	WorkDir         string
	TempDir         string
	StartConfigFile string
	StartTime       time.Time
	LoadTime        time.Time
)
