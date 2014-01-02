package estable

import (
	"bmautil/valutil"
	"esp/shell"
	"uprop"
)

const (
	commandNameEdit = "edit"
)

type cmdEdit struct {
	service *StableCache
}

func (this *cmdEdit) Name() string {
	return commandNameEdit
}

func (this *cmdEdit) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[reset|commit][varname varval]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}

	if fs.NArg() == 0 {
		shell.EditorHelper.DoList(s, "config", &this.service.editConfig)
	}
	if fs.NArg() == 1 {
		subcmd := fs.Arg(0)
		switch subcmd {
		case "reload", "reset":
			if this.service.config != nil {
				this.service.editConfig = *this.service.config
			} else {
				this.service.editConfig = CacheConfig{}
			}
			s.Writeln(subcmd + " done")
		case "commit":
			if err := this.service.editConfig.Valid(); err != nil {
				s.Writeln("ERROR: config invalid - " + err.Error())
				return true
			}
			cfg := new(CacheConfig)
			*cfg = this.service.editConfig
			this.service.config = cfg
			s.Writeln("commit done")
		default:
			s.Writeln("ERROR: unknow action '" + subcmd + "'")
		}
	}
	if fs.NArg() == 2 {
		varn := fs.Arg(0)
		v := fs.Arg(1)
		done := shell.EditorHelper.DoSet(s, &this.service.editConfig, varn, v)
		if done {
			shell.EditorHelper.DoList(s, "config", &this.service.editConfig)
		}
	}
	return true
}

func (this *CacheConfig) GetUProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	r = append(r, uprop.NewUProperty("maxsize", this.MaxSize, true, "cache max size", func(v string) error {
		this.MaxSize = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("queuesize", this.QueueSize, true, "executor queue size", func(v string) error {
		this.QueueSize = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("recover", this.RecoverTime, true, "recover time when invalidResponse sended, SEC", func(v string) error {
		this.RecoverTime = valutil.ToInt(v, 0)
		return nil
	}))
	return r
}
