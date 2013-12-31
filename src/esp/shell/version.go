package shell

import (
	"bmautil/valutil"
	"fmt"
	"strings"
)

type VersionCommand struct {
	appName string
	version func() string
}

func newVersionCommand(appName string, version func() string) *VersionCommand {
	r := new(VersionCommand)
	r.appName = appName
	r.version = version
	return r
}

func (this *VersionCommand) Name() string {
	return "version"
}

func (this *VersionCommand) checkVersion(ver string) bool {
	v1l := strings.Split(this.version(), ".")
	v2l := strings.Split(ver, ".")
	for i, v2 := range v2l {
		if v2 == "*" {
			continue
		}
		if i >= len(v1l) {
			return false
		}
		i1 := valutil.ToInt(v1l[i], 0)
		i2 := valutil.ToInt(v2, 0)
		if i1 < i2 {
			return false
		}
	}
	return true
}

func (this *VersionCommand) Process(s *Session, command string) bool {
	if CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	check := false
	match := false

	args := "[appName] [version]"
	fs := NewFlagSet(name)
	fs.BoolVar(&check, "c", false, "check version")
	fs.BoolVar(&match, "m", false, "match version")
	if DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}

	if match {
		if fs.NArg() >= 2 {
			appName := fs.Arg(0)
			ver := fs.Arg(1)
			if appName != this.appName {
				s.Writeln(fmt.Sprintf("ERROR: appName '%s' not '%s'", this.appName, appName))
				return true
			}

			ok := false
			if check {
				ok = this.checkVersion(ver)
			} else {
				ok = this.version() == ver
			}

			if !ok {
				s.Writeln(fmt.Sprintf("ERROR: version '%s' not '%s'", this.version(), ver))
			} else {
				s.Writeln("version match ok")
			}
		} else {
			s.Writeln("ERROR: params invalid")
		}
	} else {
		s.Writeln(fmt.Sprintf("version %s %s", this.appName, this.version()))
	}
	return true
}
