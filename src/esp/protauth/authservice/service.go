package authservice

import (
	"errors"
	"esp/espnet"
	"fmt"
)

type AuthExecutor interface {
	//
	Login(user string, token string) (string, error)
}

type AuthServiceHandler struct {
	executor AuthExecutor
}

func NewAuthServiceHandler(exec AuthExecutor) *AuthServiceHandler {
	this := new(AuthServiceHandler)
	this.executor = exec
	return this
}

func (this *AuthServiceHandler) Serve(msg *espnet.Message, rep espnet.ServiceResponser) error {
	m, _ := msg.Headers().GetString("method", "")
	switch m {
	case "Login":
		return this.login(msg, rep)
	}
	return errors.New(fmt.Sprintf("method '%s' not exists", m))
}

func (this *AuthServiceHandler) login(msg *espnet.Message, rep espnet.ServiceResponser) error {
	form := msg.Datas()
	user, _ := form.GetString("user", "")
	token, _ := form.GetString("token", "")

	// do it
	ntoken, err := this.executor.Login(user, token)

	// reply
	rmsg := msg.ReplyMessage()
	if err != nil {
		rmsg.BeError(err)
		return rep(rmsg)
	}
	rmsg.Datas().Set("token", ntoken)
	return rep(rmsg)
}
