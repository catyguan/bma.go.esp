package protauth

import (
	"esp/espnet"
)

type IAuthToken interface {
	IsValid() bool
	Bind(msg *espnet.Message)
}

type AuthToken struct {
	token string
}

func NewAuthToken(t string) *AuthToken {
	this := new(AuthToken)
	this.token = t
	return this
}

func (this *AuthToken) IsValid() bool {
	return this.token != ""
}

func (this *AuthToken) Bind(msg *espnet.Message) {
	msg.Headers().Set("token", this.token)
}
