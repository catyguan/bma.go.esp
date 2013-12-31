package retryst

import (
	"testing"
	"time"
)

func TestRetryBase(t *testing.T) {
	cfg := new(RetryConfig)
	cfg.DelayMin = 100
	cfg.DelayIncrease = 100
	cfg.DelayMax = 400
	cfg.Max = 3

	rst := new(RetryState)
	rst.Config = cfg
	rst.Begin(func(lastTry bool) bool {

		t.Error(time.Now())
		if lastTry {
			t.Error("lastTry!")
		}
		if true {
			panic("test")
		}
		return false
	})
	time.Sleep(time.Duration(2) * time.Second)
	rst.Cancel()
	t.Error(rst)
}
