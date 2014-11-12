package espmux4smmapi

import (
	"crypto/md5"
	"encoding/json"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"smmapi"
	"strings"
)

type Service struct {
	name   string
	config *configInfo
}

func NewService(n string) *Service {
	r := new(Service)
	r.name = n
	return r
}

func (this *Service) InitMuxInvoke(mux *espservice.ServiceMux, serName, opName string) {
	mux.AddHandler(serName, opName, func(sock *espsocket.Socket, msg *esnp.Message) error {
		form := msg.Datas()
		param := make(map[string]interface{})
		id, err1 := form.GetString("id", "")
		if err1 != nil {
			return err1
		}
		aid, err2 := form.GetString("aid", "")
		if err2 != nil {
			return err2
		}
		strparam, err3 := form.GetString("param", "")
		if err3 != nil {
			return err3
		}
		tmpcode, err4 := form.GetString("code", "")
		if err4 != nil {
			return err4
		}
		code := strings.ToLower(tmpcode)
		mycode := ""
		cfgcode := this.config.Code
		if cfgcode != "" {
			tmp := fmt.Sprintf("%d/%s/%s/%s", id, aid, strparam, cfgcode)
			h := md5.New()
			h.Write([]byte(tmp))
			mycode = fmt.Sprintf("%x", h.Sum(nil))
		}
		if code != mycode {
			return logger.Warn(tag, "code invalid(in=%s, my=%s)", code, mycode)
		}
		if strparam != "" {
			err1 := json.Unmarshal([]byte(strparam), &param)
			if err1 != nil {
				return err1
			}
		}

		result, err := smmapi.Invoke(id, aid, param)

		r := make(map[string]interface{})
		if err != nil {
			r["Status"] = 500
			r["Error"] = err.Error()
		} else {
			r["Status"] = 200
			r["Result"] = result
		}

		rmsg := msg.ReplyMessage()
		datas := rmsg.Datas()
		datas.Set("Content", r)

		return sock.SendMessage(rmsg, nil)
	})
}
