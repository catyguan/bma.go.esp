package shell

import (
	"fmt"
	"sort"
	"strings"
)

type LSCommand struct {
}

func NewLSCommand() *LSCommand {
	r := new(LSCommand)
	return r
}

func (this *LSCommand) Name() string {
	return "ls"
}

func (this *LSCommand) Process(s *Session, command string) bool {
	if CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := ""
	cmdOnly := false
	dirOnly := false
	dirList := false
	fs := NewFlagSet(name)
	fs.BoolVar(&cmdOnly, "c", false, "list command only")
	fs.BoolVar(&dirOnly, "d", false, "list dir only")
	fs.BoolVar(&dirList, "l", false, "list dir only and show dir info")
	if DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}
	if dirList {
		dirOnly = true
		cmdOnly = false
	}
	ExecuteLS(s, cmdOnly, dirOnly, dirList)
	return true
}

func ExecuteLS(s *Session, cmdOnly, dirOnly, dirList bool) bool {
	o := s.Get("@DIR", nil).(*dirInfo)
	pwd := o.pwd
	slist := pwd.List()
	sort.Strings(slist)
	if cmdOnly || dirOnly {
		tlist := make([]string, 0, len(slist))
		for _, s := range slist {
			if strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">") {
				if cmdOnly {
					continue
				}
				if dirList {
					ns := s
					odir := o.pwd.GetDir(s[1 : len(s)-1])
					if odir != nil {
						if oinfo, ok := odir.(DirInfoSupport); ok {
							ds := oinfo.DirInfo()
							if ds != "" {
								ns += ": " + oinfo.DirInfo()
							}
						}
					}
					s = ns
				}
			} else {
				if dirOnly {
					continue
				}
			}
			tlist = append(tlist, s)
		}
		slist = tlist
	}

	s.Writeln(fmt.Sprintf("%s ->", o.CurName()))
	for i, str := range slist {
		if i > 0 {
			if dirList {
				s.Write("\n")
			} else {
				s.Write(", ")
			}
		}
		s.Write(str)
	}
	s.Writeln("")
	return true
}
