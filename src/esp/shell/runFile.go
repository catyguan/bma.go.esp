package shell

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type RunFileCommand struct {
	shell *Shell
}

func NewRunFileCommand(shell *Shell) *RunFileCommand {
	r := new(RunFileCommand)
	r.shell = shell
	return r
}

func (this *RunFileCommand) Name() string {
	return "runfile"
}

func (this *RunFileCommand) Process(s *Session, command string) bool {
	if CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "fileName"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}
	ExecuteFile(this.shell, s, fs.Arg(0))
	return true
}

func ExecuteFile(shell *Shell, s *Session, file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: file '%s' load fail => %s\n", file, err))
		return
	}

	clist := strings.Split(string(content), "\n")
	for _, line := range clist {
		str := strings.TrimSpace(line)
		if str == "" {
			continue
		}
		s.Writeln(" >>> " + str)
		shell.doProcess(s, str, true)
	}
}
