package cacheserver

import (
	"esp/shell"
	"fmt"
)

const (
	commandNameGet = "cget"
)

type cmdGet struct {
	service *CacheService
	name    string
}

func (this *cmdGet) Name() string {
	return commandNameGet
}

func (this *cmdGet) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	notLoad := false
	trace := false
	timeout := 5

	name := this.Name()
	args := "key"
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&notLoad, "n", false, "not load")
	fs.BoolVar(&trace, "t", false, "enable trace mode")
	fs.IntVar(&timeout, "to", 5, "timeout in seconds")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}

	req := NewGetRequest(fs.Arg(0))
	req.TimeoutMs = int32(timeout) * 1000
	req.NotLoad = notLoad
	req.Trace = trace

	rep := make(chan *GetResult, 1)
	defer close(rep)

	err := this.service.Get(this.name, req, rep)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return true
	}

	result := <-rep
	if result == nil {
		s.Writeln(fmt.Sprintf("ERROR: '%s:%s' null result return", this.name, req.Key))
		return true
	}
	if result.Err != nil {
		s.Writeln("ERROR: " + result.Err.Error())
		return true
	}

	s.Write(fmt.Sprintf("'%s:%s' -> ", this.name, req.Key))
	if !result.Done {
		s.Writeln("<miss>")
	} else {
		if result.Value != nil {
			s.Writeln(string(result.Value))
		} else {
			s.Writeln("<empty data>")
		}
	}
	if trace {
		s.Writeln("TRACE ->")
		if result.TraceInfo != nil {
			for i, str := range result.TraceInfo {
				s.Writeln(fmt.Sprintf("%d: %s", i+1, str))
			}
		}
	}
	return true
}
