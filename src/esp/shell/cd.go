package shell

import (
	"fmt"
	"strings"
)

type CDCommand struct {
}

func NewCDCommand() *CDCommand {
	r := new(CDCommand)
	return r
}

func (this *CDCommand) Name() string {
	return "cd"
}

func (this *CDCommand) Process(s *Session, command string) bool {
	if CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "groupName"
	fs := NewFlagSet(name)
	if DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}
	ExecuteCD(s, fs.Args())
	return true
}

func ExecuteCD(s *Session, args []string) bool {
	o := s.Get("@DIR", nil).(*dirInfo)
	r := true
	if len(args) > 0 {
		name := args[0]

		if name == "/" {
			for {
				l := len(o.stack)
				if l > 0 {
					o.pwd = o.stack[l-1]
					o.stack = o.stack[:l-1]
				} else {
					break
				}
			}
		} else {
			ps := strings.Split(name, "/")
			for _, n := range ps {
				switch n {
				case ".": // do thing
				case "..":
					l := len(o.stack)
					if l > 0 {
						o.pwd = o.stack[l-1]
						o.stack = o.stack[:l-1]
					}
				default:
					if o.pwd != nil {
						ndir := o.pwd.GetDir(n)
						if ndir == nil {
							s.Writeln(fmt.Sprintf("ERROR: '%s' not found '%s'", o.CurName(), n))
							break
						}
						o.stack = append(o.stack, o.pwd)
						o.pwd = ndir
					}
				}
			}
		}
	}

	s.Writeln(fmt.Sprintf("-> %s", o.CurName()))
	return r
}
