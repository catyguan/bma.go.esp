package namedsql

import (
	"bmautil/valutil"
	"esp/shell"
)

const (
	commandNameNewSQL   = "new"
	sessionKeySQLConfig = "@namedSQL.config"
)

type cmdNewSQL struct {
	server *NamedSQL
}

func (this *cmdNewSQL) Name() string {
	return commandNameNewSQL
}

func (this *cmdNewSQL) cur(s *shell.Session, c bool) *SQLConfig {
	var f func() interface{}
	if c {
		f = func() interface{} {
			return new(SQLConfig)
		}
	}
	r := s.Get(sessionKeySQLConfig, f)
	if r != nil {
		return r.(*SQLConfig)
	}
	return nil
}

func (this *cmdNewSQL) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[start|reset|commit|varname varval]"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 2) {
		return true
	}

	if fs.NArg() == 0 {
		this.showCurrent(s)
	}
	if fs.NArg() == 1 {
		subcmd := fs.Arg(0)
		switch subcmd {
		case "start":
			this.cur(s, true)
			this.showCurrent(s)
		case "reset":
			delete(s.Vars, sessionKeySQLConfig)
			s.Writeln("reset done")
		case "commit":
			cfg := this.cur(s, false)
			if cfg == nil {
				s.Writeln("ERROR: no working")
			} else {
				err := this.server.CreateSQL(cfg)
				if err != nil {
					s.Writeln("ERROR: " + err.Error())
				} else {
					s.Writeln("commit done")
					delete(s.Vars, sessionKeySQLConfig)
				}
			}
		default:
			s.Writeln("ERROR: unknow action '" + subcmd + "'")
		}
	}
	if fs.NArg() == 2 {
		varn := fs.Arg(0)
		v := fs.Arg(1)
		cfg := this.cur(s, true)
		done := this.editConfig(s, cfg, varn, v)
		if done {
			this.showCurrent(s)
		}
	}
	return true
}

func (this *cmdNewSQL) showCurrent(s *shell.Session) {
	cfg := this.cur(s, false)
	if cfg == nil {
		s.Writeln("no working, use: new start")
		return
	}
	this.showConfig(s, cfg)
}

func (this *cmdNewSQL) showConfig(s *shell.Session, cfg *SQLConfig) {
	s.Writeln("[sql config]")
	shell.PrintProp(s, "name", cfg.Name, false, "sql name")
	shell.PrintProp(s, "driver", cfg.Driver, false, "sql driver")
	shell.PrintProp(s, "dsn", cfg.DataSource, false, "sql datasource name")
	shell.PrintProp(s, "delay", cfg.DelayOpen, true, "delay open")
	shell.PrintProp(s, "maxidle", cfg.MaxIdleConns, true, "max idle connections, <=0 no limit")
	shell.PrintProp(s, "maxopen", cfg.MaxOpenConns, true, "max open connections, <=0 no limit")
	s.Writeln("[sql config end]")
}

func (this *cmdNewSQL) editConfig(s *shell.Session, cfg *SQLConfig, n, v string) bool {
	done := true
	switch n {
	case "name":
		cfg.Name = v
	case "driver":
		cfg.Driver = v
	case "dsn":
		cfg.DataSource = v
	case "delay":
		cfg.DelayOpen = valutil.ToBool(v, false)
	case "maxidle":
		cfg.MaxIdleConns = valutil.ToInt(v, 0)
	case "maxopen":
		cfg.MaxOpenConns = valutil.ToInt(v, 0)
	default:
		done = false
		s.Writeln("ERROR: unknow var '" + n + "'")
	}
	return done
}
