package xmem

import (
	"bmautil/binlog"
	"bmautil/valutil"
	"esp/shell"
	"fmt"
	"strings"
)

type dirGroup struct {
	shell.ShellDirBase
	service *Service
	name    string
}

func (this *dirGroup) InitDir(s *Service, n string) {
	this.service = s
	this.name = n
	this.DirName = this.service.Name()
	this.Commands = this.MakeCommands
}

func (this *dirGroup) MakeCommands() map[string]shell.ShellProcessor {
	r := make(map[string]shell.ShellProcessor)
	r["save"] = this.CF(this.commandSave)
	r["binlog"] = this.CF(this.commandBinlog)
	r["load"] = this.CF(this.commandLoad)
	r["dump"] = this.CF(this.commandDump)
	return r
}

func (this *dirGroup) commandSave(s *shell.Session, command string) bool {
	name := "save"
	args := ""
	fileName := ""
	fs := shell.NewFlagSet(name)
	fs.StringVar(&fileName, "f", "", "file to store snapshot")
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.SaveMemGroup(this.name, fileName)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("save done")
	}

	return true
}

func (this *dirGroup) commandLoad(s *shell.Session, command string) bool {
	name := "load"
	args := ""
	fileName := ""
	fs := shell.NewFlagSet(name)
	fs.StringVar(&fileName, "f", "", "file to load snapshot")
	if shell.DoParse(s, command, fs, name, args, 0, 0) {
		return true
	}

	err := this.service.LoadMemGroup(this.name, fileName)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("load done")
	}

	return true
}

func (this *dirGroup) commandDump(s *shell.Session, command string) bool {
	name := "dump"
	args := "[key]"
	all := false
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&all, "a", false, "dump all sub items")
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}

	key := ""
	if fs.NArg() > 0 {
		key = fs.Arg(0)
	}

	k := MemKeyFromString(key)
	str, err := this.service.Dump(this.name, k, all)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
	} else {
		s.Writeln("dump done >>")
		s.Writeln(str)
	}

	return true
}

func (this *dirGroup) commandBinlog(s *shell.Session, command string) bool {
	name := "binlog"
	args := "save filename OR ver [slaveVer] OR start OR stop"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 1, 2) {
		return true
	}

	act := strings.ToLower(fs.Arg(0))
	p1 := ""
	if fs.NArg() > 1 {
		p1 = fs.Arg(1)
	}

	switch act {
	case "save":
		if p1 == "" {
			s.Writeln("ERROR: save file name empty")
		} else {
			err := this.service.SaveBinlogSnapshot(this.name, p1)
			if err != nil {
				s.Writeln(fmt.Sprintf("ERROR: save binlog fail - %s", err))
			} else {
				s.Writeln("save binlog snapshot done")
			}
		}
	case "ver":
		if p1 == "" {
			var mv, sv binlog.BinlogVer
			err := this.service.executor.DoSync("ver", func() error {
				var err error
				mv, sv, err = this.service.doGetBinogVersion(this.name)
				return err
			})
			if err != nil {
				s.Writeln(fmt.Sprintf("ERROR: binlog version fail - %s", err))
			} else {
				s.Writeln(fmt.Sprintf("binlog master version : %d", mv))
				s.Writeln(fmt.Sprintf("binlog slave version : %d", sv))
			}
		} else {
			ver := valutil.ToInt64(p1, -1)
			if ver < 0 {
				s.Writeln("ERROR: version value invalid")
			} else {
				err := this.service.executor.DoSync("setver", func() error {
					si, err := this.service.doGetGroup(this.name)
					if err != nil {
						return err
					}
					si.group.blver = binlog.BinlogVer(ver)
					return nil
				})
				if err != nil {
					s.Writeln(fmt.Sprintf("ERROR: binlog version fail - %s", err))
				} else {
					s.Writeln(fmt.Sprintf("set binlog slave version done => %d", ver))
				}
			}
		}
	default:
		s.Writeln(fmt.Sprintf("ERROR: unknow action '%s'", act))
	}
	return true
}
