package authclient

import (
	"errors"
	"esp/espnet"
	"esp/espnet/msgagent"
	"esp/protauth"
	"logger"
)

const (
	tag = "authclient"
)

type AuthClient struct {
	serviceAddress espnet.Address
	channel        espnet.Channel
	Agent          *msgagent.Agent
}

func NewAuthClient(ch espnet.Channel, addr espnet.Address) *AuthClient {
	this := new(AuthClient)
	this.channel = ch
	this.serviceAddress = addr
	return this
}

func (this *AuthClient) Login(user string, certificate string) (protauth.IAuthToken, error) {
	msg := espnet.NewRequestMessage()

	msg.SetAddress(this.serviceAddress)

	hs := msg.Headers()
	hs.Set("method", "Login")

	req := msg.Datas()
	req.Set("user", user)
	req.Set("token", certificate)

	ag := msgagent.S(this.Agent)
	rmsg, err := ag.SendMessage(this.channel, msg, 0)
	if err != nil {
		logger.Debug(tag, "Login(%s) send error - %s", user, err)
		return nil, err
	}
	err2 := rmsg.ToError()
	if err2 != nil {
		logger.Debug(tag, "Login(%s) fail - %s", user, err2)
		return nil, err2
	}
	token, _ := rmsg.Datas().GetString("token", "")
	if token == "" {
		logger.Debug(tag, "Login(%s) invalid", user)
		return nil, errors.New("invalid")
	}
	return protauth.NewAuthToken(token), nil
}
