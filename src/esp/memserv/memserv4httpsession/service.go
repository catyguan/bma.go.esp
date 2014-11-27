package memserv4httpsession

import (
	"acl"
	"esp/memserv"
	"esp/memserv/memserv4session"
	"fmt"
	"net/http"
	"time"
	"uuid"
)

const (
	tag = "memserv4httpsession"
)

type Service struct {
	name string
	s    *memserv.MemoryServ
	cfg  *config
}

func NewService(n string, s *memserv.MemoryServ) *Service {
	r := new(Service)
	r.name = n
	r.s = s
	return r
}

func (this *Service) sessionId(r *http.Request) (string, error) {
	ck, err := r.Cookie(this.cfg.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	return ck.Value, nil
}

func (this *Service) CreateSessionId() string {
	var s string
	u, err := uuid.NewV4()
	if err != nil {
		s = fmt.Sprintf("%d", time.Now().Nanosecond())
	} else {
		s = u.String()
	}
	return this.cfg.SessionPrefix + s
}

func (this *Service) NewSession(w http.ResponseWriter) string {
	sid := this.CreateSessionId()
	ck := new(http.Cookie)
	ck.Name = this.cfg.CookieName
	ck.Value = sid
	ck.Path = "/"
	http.SetCookie(w, ck)
	return sid
}

func (this *Service) GetSession(r *http.Request, key string) (interface{}, error) {
	sid, err := this.sessionId(r)
	if err != nil {
		return nil, err
	}
	if sid == "" {
		return nil, nil
	}
	return memserv4session.GetSession(this.s, sid, key, this.cfg.TimeoutMS)
}

func (this *Service) MGetSession(r *http.Request, keys []string) (interface{}, error) {
	sid, err := this.sessionId(r)
	if err != nil {
		return nil, err
	}
	if sid == "" {
		return nil, nil
	}
	return memserv4session.MGetSession(this.s, sid, keys, this.cfg.TimeoutMS)
}

func (this *Service) SetSession(w http.ResponseWriter, r *http.Request, key string, val interface{}) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		sid = this.NewSession(w)
	}
	return memserv4session.SetSession(this.s, sid, key, val, this.cfg.TimeoutMS)
}

func (this *Service) MSetSession(w http.ResponseWriter, r *http.Request, mv map[string]interface{}) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		sid = this.NewSession(w)
	}
	return memserv4session.MSetSession(this.s, sid, mv, this.cfg.TimeoutMS)
}

func (this *Service) DeleteSession(r *http.Request, key string) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}
	return memserv4session.DeleteSession(this.s, sid, key, this.cfg.TimeoutMS)
}

func (this *Service) MDeleteSession(r *http.Request, keys []string) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}
	return memserv4session.MDeleteSession(this.s, sid, keys, this.cfg.TimeoutMS)
}

func (this *Service) CloseSession(r *http.Request) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}
	return memserv4session.CloseSession(this.s, sid)
}

func (this *Service) GetUser(r *http.Request) (*acl.User, error) {
	sid, err := this.sessionId(r)
	if err != nil {
		return nil, err
	}
	if sid == "" {
		return nil, nil
	}
	return memserv4session.QueryUser(this.s, sid, this.cfg.TimeoutMS)
}

func (this *Service) SetUser(w http.ResponseWriter, r *http.Request, user *acl.User) error {
	sid, err := this.sessionId(r)
	if err != nil {
		return err
	}
	if sid == "" {
		sid = this.NewSession(w)
	}
	return memserv4session.SetUser(this.s, sid, user, this.cfg.TimeoutMS)
}
