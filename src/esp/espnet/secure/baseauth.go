package secure

import (
	"bmautil/valutil"
	"crypto/md5"
	"encoding/hex"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"time"
)

type BaseAuthRequest struct {
	App   string
	Token string
}

func (this *BaseAuthRequest) Encode(msg *esnp.Message) error {
	if this.Token != "" {
		esnp.MessageLineCoders.XData.Add(msg, 1, this.Token, esnp.Coders.String)
	}
	if this.App != "" {
		esnp.MessageLineCoders.XData.Add(msg, 2, this.App, esnp.Coders.String)
	}
	return nil
}

func (this *BaseAuthRequest) Decode(msg *esnp.Message) error {
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 1, esnp.Coders.String)
		if err != nil {
			return err
		}
		this.Token = valutil.ToString(v, "")
	}
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 2, esnp.Coders.String)
		if err != nil {
			return err
		}
		this.App = valutil.ToString(v, "")
	}
	return nil
}

func (this *BaseAuthRequest) Valid() error {
	return nil
}

func (this *BaseAuthRequest) Reset() {
	this.Token = ""
	this.App = ""
}

type BaseAuthResponse struct {
	Token string
	Done  int
}

func (this *BaseAuthResponse) Encode(msg *esnp.Message) error {
	if this.Token != "" {
		esnp.MessageLineCoders.XData.Add(msg, 1, this.Token, esnp.Coders.String)
	}
	if this.Done != 0 {
		esnp.MessageLineCoders.XData.Add(msg, 1, this.Done == 1, esnp.Coders.Bool)
	}
	return nil
}

func (this *BaseAuthResponse) Decode(msg *esnp.Message) error {
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 1, esnp.Coders.String)
		if err != nil {
			return err
		}
		this.Token = valutil.ToString(v, "")
	}
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 1, esnp.Coders.Bool)
		if err != nil {
			return err
		}
		if v != nil && valutil.ToBool(v, false) {
			this.Done = 1
		}
	}
	return nil
}

func (this *BaseAuthResponse) IsDone() bool {
	return this.Done == 1
}

func (this *BaseAuthResponse) Valid() error {
	return nil
}

func (this *BaseAuthResponse) Reset() {
	this.Token = ""
	this.Done = 1
}

func CreateAuthToken(tk string, k string) string {
	h := md5.New()
	h.Write([]byte(tk))
	h.Write([]byte(k))
	return hex.EncodeToString(h.Sum(nil))
}

func DoBaseAuth(sock espsocket.Socket, app, key string, timeout time.Duration) error {
	var req BaseAuthRequest
	var rep BaseAuthResponse

	espsocket.SetDeadline(sock, time.Now().Add(timeout))
	defer espsocket.ClearDeadline(sock)

	msg1 := esnp.NewRequestMessage()
	req.App = app
	err1 := req.Encode(msg1)
	if err1 != nil {
		return err1
	}
	rmsg1, err2 := espsocket.Call(sock, msg1)
	if err2 != nil {
		return err2
	}
	err2 = rep.Decode(rmsg1)
	if err2 != nil {
		return err2
	}
	err2 = rep.Valid()
	if err2 != nil {
		return err2
	}
	autk := CreateAuthToken(rep.Token, key)
	logger.Debug(tag, "BaseAuth request app=%s,key=%s, token=%s, auth=%s", app, key, rep.Token, autk)
	req.Reset()
	req.App = app
	req.Token = autk
	msg2 := esnp.NewRequestMessage()
	err2 = req.Encode(msg2)
	if err2 != nil {
		return err2
	}
	_, err3 := espsocket.Call(sock, msg2)
	if err3 != nil {
		return err3
	}
	return nil
}

// BaseAuthEntry
type AppKeyProvider func(app string, addr string) (string, error)

func SimpleAppKeyProvider(key string) AppKeyProvider {
	return func(app string, addr string) (string, error) {
		return key, nil
	}
}

type BaseAuthEntry struct {
	BaseSecureConfig
	akp AppKeyProvider
}

func NewBaseAuthEntry(akp AppKeyProvider, e espservice.ServiceEntry) *BaseAuthEntry {
	r := new(BaseAuthEntry)
	r.InitDefault()
	r.Entry = e
	r.akp = akp
	return r
}

func (this *BaseAuthEntry) AuthEntry(sock espsocket.Socket) {
	defer sock.AskFinish()
	this.Begin(sock)
	if !this.DoAuth(sock) {
		return
	}
	this.DoNext(sock)
}

func (this *BaseAuthEntry) DoAuth(sock espsocket.Socket) bool {
	var req BaseAuthRequest
	var rep BaseAuthResponse

	msg1, err1 := sock.ReadMessage(true)
	if err1 != nil {
		logger.Debug(tag, "BaseAuth read request 1 fail - %s", err1)
		return false
	}
	err1 = req.Decode(msg1)
	if err1 != nil {
		logger.Debug(tag, "BaseAuth decode request 1 fail - %s", err1)
		return false
	}
	addr, _ := espsocket.GetProperty(sock, espsocket.PROP_SOCKET_REMOTE_ADDR)
	key, errK := this.akp(req.App, valutil.ToString(addr, ""))
	if errK != nil {
		logger.Debug(tag, "BaseAuth app(%s) key fail - %s", req.App, errK)
		return false
	}
	tk := fmt.Sprintf("%d", time.Now().UnixNano())
	logger.Debug(tag, "BaseAuth request app=%s, token=%s", req.App, tk)
	rep.Token = tk
	rmsg1 := msg1.ReplyMessage()
	err1 = rep.Encode(rmsg1)
	if err1 != nil {
		logger.Debug(tag, "BaseAuth encode response 1 fail - %s", err1)
		return false
	}
	err1 = sock.WriteMessage(rmsg1)
	if err1 != nil {
		logger.Debug(tag, "BaseAuth write response 1 fail - %s", err1)
		return false
	}

	msg2, err2 := sock.ReadMessage(true)
	if err2 != nil {
		logger.Debug(tag, "BaseAuth read request 2 fail - %s", err2)
		return false
	}
	req.Reset()
	err2 = req.Decode(msg2)
	if err2 != nil {
		logger.Debug(tag, "BaseAuth decode request 2 fail - %s", err2)
		return false
	}
	err2 = req.Valid()
	if err2 != nil {
		logger.Debug(tag, "BaseAuth valid request 2 fail - %s", err2)
		return false
	}
	atki := req.Token
	atkm := CreateAuthToken(tk, key)
	if atki != atkm {
		logger.Warn(tag, "BaseAuth %s invalid auth token (in=%s, me=%s)", sock, atki, atkm)
		return false
	}
	rep.Reset()
	rmsg2 := msg2.ReplyMessage()
	err2 = rep.Encode(rmsg2)
	if err2 != nil {
		logger.Debug(tag, "BaseAuth encode response 2 fail - %s", err2)
		return false
	}
	err2 = sock.WriteMessage(rmsg2)
	if err2 != nil {
		logger.Debug(tag, "BaseAuth write response 2 fail - %s", err2)
		return false
	}
	return true
}
