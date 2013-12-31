package netcmd

import (
	"esp/shell"
	"fmt"
	"strings"
)

type IpLimitCommand struct {
	GetBlackList func() []string
	SetBlackList func([]string)
	GetWhiteList func() []string
	SetWhiteList func([]string)
}

func NewIpLimitCommand() *IpLimitCommand {
	r := new(IpLimitCommand)
	return r
}

func (this *IpLimitCommand) Name() string {
	return "iplimit"
}

func (this *IpLimitCommand) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "ipvalue"
	isWhite := false
	isBlack := false
	isAdd := false
	isRemove := false
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&isWhite, "w", false, "handle white list")
	fs.BoolVar(&isBlack, "b", false, "handle black list")
	fs.BoolVar(&isAdd, "a", false, "add ip")
	fs.BoolVar(&isRemove, "r", false, "remove ip")
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}
	ip := ""
	if fs.NArg() > 0 {
		ip = fs.Arg(0)
	}
	if isAdd {
		if !isBlack {
			isWhite = true
		}
		this.add(s, isWhite, ip)
	} else if isRemove {
		if !isBlack {
			isWhite = true
		}
		this.remove(s, isWhite, ip)
	} else {
		this.list(s, isWhite, isBlack)
	}
	return true
}

func (this *IpLimitCommand) add(s *shell.Session, isWhite bool, ip string) {
	if ip == "" {
		s.Writeln("ERROR: ip empty")
		return
	}

	var vlist []string
	if isWhite {
		if this.GetWhiteList == nil {
			s.Writeln("ERROR: can't get white list")
			return
		}
		if this.SetWhiteList == nil {
			s.Writeln("ERROR: can't set white list")
			return
		}
		vlist = this.GetWhiteList()
	} else {
		if this.GetBlackList == nil {
			s.Writeln("ERROR: can't get black list")
			return
		}
		if this.SetBlackList == nil {
			s.Writeln("ERROR: can't set black list")
			return
		}
		vlist = this.GetBlackList()
	}
	if vlist == nil {
		vlist = make([]string, 0)
	}
	if len(vlist) == 1 && vlist[0] == "" {
		vlist[0] = ip
	} else {
		vlist = append(vlist, ip)
	}
	if isWhite {
		this.SetWhiteList(vlist)
	} else {
		this.SetBlackList(vlist)
	}
	s.Writeln("add done")
	this.list(s, isWhite, !isWhite)
}

func (this *IpLimitCommand) remove(s *shell.Session, isWhite bool, ip string) {
	if ip == "" {
		s.Writeln("ERROR: ip empty")
		return
	}

	var vlist []string
	if isWhite {
		if this.GetWhiteList == nil {
			s.Writeln("ERROR: can't get white list")
			return
		}
		if this.SetWhiteList == nil {
			s.Writeln("ERROR: can't set white list")
			return
		}
		vlist = this.GetWhiteList()
	} else {
		if this.GetBlackList == nil {
			s.Writeln("ERROR: can't get black list")
			return
		}
		if this.SetBlackList == nil {
			s.Writeln("ERROR: can't set black list")
			return
		}
		vlist = this.GetBlackList()
	}
	if vlist == nil {
		vlist = make([]string, 0)
	}
	nlist := make([]string, 0, len(vlist))
	for _, v := range vlist {
		if v == ip {
			continue
		}
		nlist = append(nlist, v)
	}
	if isWhite {
		this.SetWhiteList(nlist)
	} else {
		this.SetBlackList(nlist)
	}
	s.Writeln("remove done")
	this.list(s, isWhite, !isWhite)
}

func printList(s *shell.Session, name string, list []string) {
	var str string = ""
	if list != nil {
		str = strings.Join(list, ",")
	}
	s.Write(fmt.Sprintf("%s -> %s\n", name, str))
}

func (this *IpLimitCommand) list(s *shell.Session, isWhite, isBlack bool) {
	isAll := false
	if !isWhite && !isBlack {
		isAll = true
		isWhite = true
		isBlack = true
	}
	if isBlack {
		if this.GetBlackList == nil {
			if !isAll {
				s.Writeln("ERROR: can't get black list")
				return
			}
		} else {
			printList(s, "black list", this.GetBlackList())
		}
	}
	if isWhite {
		if this.GetWhiteList == nil {
			if !isAll {
				s.Writeln("ERROR: can't get white list")
				return
			}
		} else {
			printList(s, "white list", this.GetWhiteList())
		}
	}
	return
}
