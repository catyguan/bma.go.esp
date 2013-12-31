package authservice

import (
	"config"
	"errors"
	"strings"
	"uuid"
)

type SimpleAuthExecutor struct {
	name       string
	userTokens map[string]string
}

func NewSimpleAuthExecutor(name string) *SimpleAuthExecutor {
	this := new(SimpleAuthExecutor)
	this.name = name
	this.userTokens = make(map[string]string)
	return this
}

func (this *SimpleAuthExecutor) InitUser(user, token string) {
	this.userTokens[user] = token
}

func (this *SimpleAuthExecutor) Login(user string, token string) (string, error) {
	t, ok := this.userTokens[user]
	if !ok {
		return "", errors.New("user not exists")
	}
	if t != token {
		return "", errors.New("certificate fail")
	}
	id, err2 := uuid.NewV4()
	if err2 != nil {
		return "", err2
	}
	return id.String(), nil
}

type simpleAuthConfig struct {
	User string
}

func (this *SimpleAuthExecutor) Init() bool {
	cfg := simpleAuthConfig{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.User != "" {
			slist := strings.Split(cfg.User, ",")
			for _, utword := range slist {
				utp := strings.SplitN(utword, ":", 2)
				if len(utp) != 2 {
					continue
				}
				user := strings.TrimSpace(utp[0])
				token := strings.TrimSpace(utp[1])
				this.InitUser(user, token)
			}
		}
	}
	return true
}
