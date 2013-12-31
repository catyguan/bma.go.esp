package eloader

import "esp/shell"

const (
	commandNameDeploy = "deploy"
)

type cmdDeploy struct {
	service *LoaderCache
}

func (this *cmdDeploy) Name() string {
	return commandNameDeploy
}

func (this *cmdDeploy) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}
	err := this.service.Deploy()
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}
	s.Writeln("deploy done")
	return true
}
