package namedsql

import (
	"esp/shell"
	"fmt"
)

type sqlDir struct {
	server *NamedSQL
}

func (this *sqlDir) Name() string {
	return this.server.Name()
}

func (this *sqlDir) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case commandNameNewSQL:
		return &cmdNewSQL{this.server}
	case commandNameDeleteSQL:
		return &cmdDeleteSQL{this.server, ""}
	}
	return nil
}

func (this *sqlDir) GetDir(name string) shell.ShellDir {
	if _, ok := this.server.databases[name]; ok {
		return &databaseDir{this.server, name}
	}
	return nil
}

func (this *sqlDir) List() []string {
	r := make([]string, 0)
	for k, _ := range this.server.databases {
		r = append(r, shell.DirName(k))
	}
	cmds := []string{
		commandNameNewSQL,
		commandNameDeleteSQL,
	}
	for _, k := range cmds {
		r = append(r, k)
	}
	return r
}

type databaseDir struct {
	server *NamedSQL
	name   string
}

func (this *databaseDir) Name() string {
	return this.name
}

func (this *databaseDir) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case "exec":
		return newExecCommand(this.server, this.name)
	case "query":
		return newQueryCommand(this.server, this.name)
	}
	return nil
}

func (this *databaseDir) GetDir(name string) shell.ShellDir {
	return nil
}

func (this *databaseDir) List() []string {
	r := []string{
		"exec",
		"query",
	}
	return r
}

func (this *databaseDir) DirInfo() string {
	if db, ok := this.server.databases[this.name]; ok {
		return fmt.Sprintf("%s, %s, %d", db.config.Driver, db.config.DataSource, db.config.MaxIdleConns)
	}
	return ""
}

func (this *NamedSQL) NewShellDir() shell.ShellDir {
	return &sqlDir{this}
}
