package auth

import (
	"acl"
	"bmautil/valutil"
	"esp/espnet/aclesnpmux"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
)

const (
	tag = "auth"
)

type Service struct {
	auths  []NodeAuth
	handle espservice.ServiceHandler
}

func NewService() *Service {
	r := new(Service)
	return r
}

func (this *Service) Bind(h espservice.ServiceHandler) {
	this.handle = h
}

func (this *Service) Add(auth NodeAuth) {
	this.auths = append(this.auths, auth)
}

func (this *Service) Auth(sock *espsocket.Socket, msg *esnp.Message) (bool, *acl.User, error) {
	for _, auth := range this.auths {
		done, user, err := auth(sock, msg)
		if done {
			return done, user, err
		}
		if err != nil {
			return true, nil, err
		}
		if user != nil {
			return true, user, nil
		}
	}
	return false, nil, nil
}

func (this *Service) DoServe(sock *espsocket.Socket, msg *esnp.Message) error {
	user := aclesnpmux.GetUser(sock)
	if user == nil {
		addr := msg.GetAddress()
		if addr.GetService() != SN_AUTH {
			rmsg := msg.ReplyMessage()
			hs := rmsg.Headers()
			hs.Set(HEADER_AUTH_FLAG, false)
			sock.PostMessage(rmsg)
			return nil
		}
		if addr.GetOp() != OP_AUTH {
			return espservice.Miss(msg)
		}
		ok, user, err := this.Auth(sock, msg)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if user == nil {
			return nil
		}
		aclesnpmux.SetUser(sock, user)
		if true {
			rmsg := msg.ReplyMessage()
			hs := rmsg.Headers()
			hs.Set(HEADER_AUTH_FLAG, true)
			sock.PostMessage(rmsg)
		}
		return nil
	}
	return this.handle(sock, msg)
}

func GetAuthType(msg *esnp.Message) string {
	hs := msg.Headers()
	r, _ := hs.GetString(HEADER_AUTH_TYPE, "")
	return r
}

func SetAuthType(msg *esnp.Message, t string) {
	hs := msg.Headers()
	hs.Set(HEADER_AUTH_TYPE, t)
}

func CheckAuth(msg *esnp.Message) (bool, bool) {
	hs := msg.Headers()
	v, _ := hs.Get(HEADER_AUTH_FLAG)
	if v == nil {
		return false, false
	}
	return true, valutil.ToBool(v, false)
}
