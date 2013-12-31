package eloader

import (
	"esp/shell"
	"fmt"
	"uprop"
)

const (
	commandNameNewLoader   = "load"
	sessionKeyLoaderConfig = "loaderCache.load.session"
)

type cmdNewLoader struct {
	service *LoaderCache
}

func (this *cmdNewLoader) Name() string {
	return commandNameNewLoader
}

func (this *cmdNewLoader) cur(s *shell.Session, c bool) *LoaderConfig {
	var f func() interface{}
	if c {
		f = func() interface{} {
			return new(LoaderConfig)
		}
	}
	r := s.Get(sessionKeyLoaderConfig, f)
	if r != nil {
		return r.(*LoaderConfig)
	}
	return nil
}

func (this *cmdNewLoader) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	list := false

	name := this.Name()
	args := "[start|reset|commit|varname varval]"
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&list, "l", false, "list all loader type")
	if shell.DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}

	if fs.NArg() == 0 {
		if list {
			this.listTypes(s)
		} else {
			this.showCurrent(s)
		}
	}
	if fs.NArg() == 1 {
		subcmd := fs.Arg(0)
		switch subcmd {
		case "start":
			this.cur(s, true)
			this.showCurrent(s)
		case "reset":
			delete(s.Vars, sessionKeyLoaderConfig)
			s.Writeln("reset done")
		case "commit":
			this.commit(s)
		default:
			s.Writeln("ERROR: unknow action '" + subcmd + "'")
		}
	}
	if fs.NArg() == 2 {
		varn := fs.Arg(0)
		v := fs.Arg(1)
		cfg := this.cur(s, true)
		done := false
		if varn == "type" {
			p := GetLoaderProvider(v)
			if p == nil {
				s.Writeln(fmt.Sprintf("ERROR: unknow loader type '%s'", v))
				this.listTypes(s)
				return true
			}
			done = true
			cfg.Type = v
			cfg.prop = p.CreateProperty()
			varn = ""
		}
		if varn != "" {
			var editor shell.Editor
			done = editor.Edit(s, cfg, varn, v)
		}
		if done {
			this.showCurrent(s)
		}
	}
	return true
}

func (this *cmdNewLoader) commit(s *shell.Session) {
	cfg := this.cur(s, false)
	if cfg == nil {
		s.Writeln("ERROR: no working")
		return
	}

	p := GetLoaderProvider(cfg.Type)
	if p == nil {
		s.Writeln("ERROR: loader type invalid")
		this.listTypes(s)
		return
	}

	err := this.service.AddLoader(cfg)
	if err != nil {
		s.Writeln("ERROR: " + err.Error())
		return
	}
	s.Writeln("commit done")
	delete(s.Vars, sessionKeyLoaderConfig)
}

func (this *cmdNewLoader) listTypes(s *shell.Session) {
	first := true
	s.Writeln("loader types -> ")
	for k, _ := range loaderProviders {
		if !first {
			s.Write(", ")
		}
		s.Write(k)
		first = false
	}
	s.Writeln("")
}

func (this *cmdNewLoader) showCurrent(s *shell.Session) {
	cfg := this.cur(s, false)
	if cfg == nil {
		s.Writeln("no working, use: load start")
		return
	}
	var editor shell.Editor
	editor.List(s, "loader", cfg)
}

func (this *LoaderConfig) GetUProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	r = append(r, uprop.NewUProperty("name", this.Name, false, "loader name", func(v string) error {
		this.Name = v
		return nil
	}))
	r = append(r, uprop.NewUProperty("type", this.Type, false, "loader type", func(v string) error {
		this.Type = v
		return nil
	}))
	if this.prop != nil {
		plist := this.prop.GetUProperties()
		for _, p := range plist {
			r = append(r, p)
		}
	}
	return r
}
