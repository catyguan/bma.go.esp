package cacheserver

import (
	"bmautil/valutil"
	"esp/shell"
	"time"
	//	"fmt"
)

const (
	commandNameSet = "cset"
)

type cmdSet struct {
	service *CacheService
	name    string
}

func (this *cmdSet) Name() string {
	return commandNameSet
}

func (this *cmdSet) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "key value [timeout]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 2, 3) {
		return true
	}

	key := fs.Arg(0)
	val := fs.Arg(1)
	tm := 0
	if fs.NArg() > 2 {
		tm = valutil.ToInt(fs.Arg(2), 0)
	}
	var dtime int64
	if tm > 0 {
		dtime = time.Now().Unix() + int64(tm)
	}

	err := this.service.Put(this.name, key, []byte(val), dtime)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}

	s.Writeln("set done")
	return true
}
