package tbus

import "esp/shell"

type dirService struct {
	service *TBusService
}

func (this *dirService) Name() string {
	return this.service.Name()
}

func (this *dirService) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case commandNameSave:
		return &cmdSave{this.service}
	}
	return nil
}

func (this *dirService) GetDir(name string) shell.ShellDir {
	switch name {
	case "module":
		r := new(dirModule)
		r.InitDir(this.service)
		return r
	case "remote":
		r := new(dirRemote)
		r.InitDir(this.service)
		return r
	}
	return nil
}

func (this *dirService) List() []string {
	dirs := []string{
		"module",
		"remote",
	}
	r := make([]string, 0)
	for _, k := range dirs {
		r = append(r, shell.DirName(k))
	}
	cmds := []string{
		commandNameSave,
	}
	for _, k := range cmds {
		r = append(r, k)
	}
	return r
}

func (this *TBusService) NewShellDir() shell.ShellDir {
	return &dirService{this}
}
