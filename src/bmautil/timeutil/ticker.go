package timeutil

import "time"

type ClosableTicker struct {
	ticker *time.Ticker
	C      <-chan *time.Time
	stop   chan bool
}

func NewClosableTicker(d time.Duration) *ClosableTicker {
	t := time.NewTicker(d)
	r := &ClosableTicker{ticker: t}
	r.ticker = t
	c := make(chan *time.Time, 1)
	r.C = c
	s := make(chan bool, 1)
	r.stop = s
	go func() {
		for {
			select {
			case tm := <-t.C:
				c <- &tm
			case <-s:
				c <- nil
				return
			}
		}
	}()
	return r
}

func (this *ClosableTicker) Stop() {
	this.ticker.Stop()
	if this.stop != nil {
		close(this.stop)
	}
}
