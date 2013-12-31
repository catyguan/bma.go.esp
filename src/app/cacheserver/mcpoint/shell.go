package mcpoint

import (
	"esp/shell"
)

type dirService struct {
	service *MemcachePoint
}

func (this *dirService) Name() string {
	return "routers"
}

func (this *dirService) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case commandNameAdd:
		return &cmdAdd{this.service}
	case commandNameDelete:
		return &cmdDelete{this.service}
	case commandNameList:
		return &cmdList{this.service}
	case commandNameMove:
		return &cmdMove{this.service}
	}
	return nil
}

func (this *dirService) GetDir(name string) shell.ShellDir {
	return nil
}

func (this *dirService) List() []string {
	r := make([]string, 0)
	cmds := []string{
		commandNameAdd,
		commandNameDelete,
		commandNameList,
		commandNameMove,
	}
	for _, k := range cmds {
		r = append(r, k)
	}
	return r
}

func (this *MemcachePoint) BuildShellDir(pdir shell.ShellDir) {
	dir := pdir.(*shell.ShellDirCommon)
	d1 := &dirService{this}
	dir.AddDir(d1)
}
