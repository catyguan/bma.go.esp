package goluaserv

import (
	"golua"
	"smmapi"
	"time"
)

type smmObject struct {
	s  *Service
	gl *golua.GoLua
}

func (this smmObject) getContent() string {
	this.s.lock.RLock()
	defer this.s.lock.RUnlock()
	gli := this.s.gli[this.gl.GetName()]
	r := ""
	if gli.status != 1 {
		r += "start=0;"
	}
	if gli.startErr != nil {
		r += "startErr=" + gli.startErr.Error()
	}
	return r
}

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Content = this.getContent()

	r.Actions = make([]*smmapi.SMAction, 0)
	if true {
		a := new(smmapi.SMAction)
		a.Id = "goluaserv.reset"
		a.Title = "Reset"
		a.Type = smmapi.SMA_API
		r.Actions = append(r.Actions, a)
	}

	return r, nil
}

// Result, refreshInfo, error
func (this smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	switch aid {
	case "goluaserv.reset":
		msg := "OK"
		time.AfterFunc(100*time.Millisecond, func() {
			this.s.ResetGoLua(this.gl.GetName())
		})
		return msg, nil
	}
	return smmapi.MISS(0), nil
}
