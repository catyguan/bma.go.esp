package eloader

import (
	"esp/shell"
)

const (
	commandNameLoaders = "loaders"
)

type cmdLoaders struct {
	service *LoaderCache
}

func (this *cmdLoaders) Name() string {
	return commandNameLoaders
}

func (this *cmdLoaders) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	var loaders []string

	f := func() error {
		if this.service.loaders != nil {
			loaders = make([]string, 0, len(this.service.loaders))
			for _, l := range this.service.loaders {
				loaders = append(loaders, l.config.String())
			}
		}
		return nil
	}

	exec := this.service.executor
	if exec == nil {
		f()
	} else {
		err := exec.DoSync("cmd_loaders", f)
		if err != nil {
			s.Writeln("ERROR: " + err.Error())
			return true
		}
	}
	s.Writeln("[loaders]")
	if loaders != nil {
		for _, str := range loaders {
			s.Writeln(str)
		}
	}
	s.Writeln("[loaders end]")
	return true
}
