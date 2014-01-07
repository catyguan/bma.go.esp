package cacheserver

import (
	"bytes"
	"errors"
	"esp/shell"
	"uprop"
)

type doc4Cache struct {
	service *CacheService
	name    string
	kind    string
	config  ICacheConfig
	edit    bool
}

func newDoc4Cache(s *CacheService, name, kind string, cfg ICacheConfig) *doc4Cache {
	r := new(doc4Cache)
	r.service = s
	r.name = name
	r.edit = name != ""
	r.kind = kind
	r.config = cfg
	return r
}

func (this *doc4Cache) Title() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("Cache - ")
	if this.name == "" {
		buf.WriteString("*")
	} else {
		buf.WriteString(this.name)
	}
	buf.WriteString("/")
	buf.WriteString(this.kind)
	return buf.String()
}

func (this *doc4Cache) commands() map[string]func(s *shell.Session, cmd string) bool {
	r := make(map[string]func(s *shell.Session, cmd string) bool)
	// r["test"] = this.commandTest
	return r
}

func (this *doc4Cache) ListCommands() []string {
	r := make([]string, 0)
	for k, _ := range this.commands() {
		r = append(r, k)
	}
	return r
}

func (this *doc4Cache) HandleCommand(session *shell.Session, cmdline string) (bool, bool) {
	cmd := shell.CommandWord(cmdline)
	f := this.commands()[cmd]
	if f != nil {
		return true, f(session, cmdline)
	}
	return false, false
}

func (this *doc4Cache) OnCloseDoc(session *shell.Session) {

}
func (this *doc4Cache) GetUProperties() []*uprop.UProperty {
	p := this.config.GetProperties()
	r := make([]*uprop.UProperty, 2+len(p))
	r[0] = uprop.NewUProperty("name", this.name, false, "module name", func(v string) error {
		if this.edit {
			return errors.New("can't edit name")
		}
		this.name = v
		return nil
	})
	r[1] = uprop.NewUProperty("type", this.kind, true, "cache type", func(v string) error {
		return errors.New("can't edit type")
	})
	copy(r[2:], p)
	return r
}

func (this *doc4Cache) CommitDoc(session *shell.Session) error {
	if this.name == "" {
		return errors.New("cache name empty")
	}
	if err := this.config.Valid(); err != nil {
		return err
	}

	if this.edit {
		cache, err := this.service.GetCache(this.name, false)
		if err != nil {
			return err
		}
		err = cache.UpdateConfig(this.config)
		if err != nil {
			return err
		}
	} else {
		ctype := this.kind
		fac := GetCacheFactory(ctype)
		cache, err := fac.CreateCache(this.config)
		if err != nil {
			return err
		}
		err = this.service.SetCache(this.name, this.kind, cache)
		if err != nil {
			return err
		}
	}
	return this.service.save()
}
