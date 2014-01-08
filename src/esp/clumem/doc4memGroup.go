package clumem

import (
	"bytes"
	"esp/shell"
	"uprop"
)

type doc4MemGroup struct {
	service *Service
	name    string
	config  *MemGroupConfig
	edit    bool
}

func newDoc4Cache(s *Service, name string, cfg *MemGroupConfig) *doc4MemGroup {
	r := new(doc4MemGroup)
	r.service = s
	r.name = name
	r.edit = name != ""
	r.config = cfg
	return r
}

func (this *doc4MemGroup) Title() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("MemGroup - ")
	if this.name == "" {
		buf.WriteString("*")
	} else {
		buf.WriteString(this.name)
	}
	return buf.String()
}

func (this *doc4MemGroup) commands() map[string]func(s *shell.Session, cmd string) bool {
	r := make(map[string]func(s *shell.Session, cmd string) bool)
	// r["test"] = this.commandTest
	return r
}

func (this *doc4MemGroup) ListCommands() []string {
	r := make([]string, 0)
	for k, _ := range this.commands() {
		r = append(r, k)
	}
	return r
}

func (this *doc4MemGroup) HandleCommand(session *shell.Session, cmdline string) (bool, bool) {
	cmd := shell.CommandWord(cmdline)
	f := this.commands()[cmd]
	if f != nil {
		return true, f(session, cmdline)
	}
	return false, false
}

func (this *doc4MemGroup) OnCloseDoc(session *shell.Session) {

}
func (this *doc4MemGroup) GetUProperties() []*uprop.UProperty {
	return this.config.GetProperties()
}

func (this *doc4MemGroup) CommitDoc(session *shell.Session) error {
	if err := this.config.Valid(); err != nil {
		return err
	}

	if this.edit {
		this.config.Name = this.name
		// cache, err := this.service.GetCache(this.name, false)
		// if err != nil {
		// 	return err
		// }
		// err = cache.UpdateConfig(this.config)
		// if err != nil {
		// 	return err
		// }
	} else {
		err := this.service.CreateMemGroup(this.config)
		if err != nil {
			return err
		}
	}
	return this.service.Save()
}
