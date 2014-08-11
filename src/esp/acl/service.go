package acl

import (
	"logger"
	"net"
)

const tag = "acl"

type Service struct {
	name   string
	config *configInfo
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	return this
}

func (this *Service) GetUserIp(id string, token string, ip net.IP) (*User, error) {
	return this.GetUser(id, token, ip.String())
}

func (this *Service) GetUser(id string, token string, ip string) (*User, error) {
	cfg := this.config
	for _, o := range cfg.Users {
		if o.Id == id {
			if o.Token == token {
				match := false
				for _, host := range o.Host {
					if host == "*" {
						match = true
					} else if host == ip {
						match = true
					}
					if match {
						// logger.Debug(tag, "user('%s') ip '%s' match", id, ip)
						break
					}
				}
				if !match {
					err := logger.Error(tag, "user('%s') ip '%s' not accept", id, ip)
					return nil, err
				}
				r := new(User)
				r.Id = o.Id
				r.Name = o.GetName(ip)
				r.Groups = o.Group
				return r, nil
			} else {
				err := logger.Error(tag, "user('%s') token '%s' no match", id, token)
				return nil, err
			}
		}
	}
	err := logger.Error(tag, "user '%s' not exists", id)
	return nil, err
}

func (this *Service) CheckOp(user *User, op string) (bool, error) {
	cfg := this.config
	for _, o := range cfg.Ops {
		if o.Op == op {
			match := false
			for _, who := range o.Who {
				if user.IsWho(who) {
					match = true
					break
				}
			}
			if !match {
				err := logger.Error(tag, "user('%s') can't access '%s'", user, op)
				return false, err
			}
			return true, nil
		}
	}
	err := logger.Error(tag, "op '%s' invalid", op)
	return false, err
}
