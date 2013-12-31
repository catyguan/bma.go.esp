package sqlite

import (
	"esp/shell"
)

type sqliteDir struct {
	sqlite *SqliteServer
}

func (this *sqliteDir) Name() string {
	return this.sqlite.Name()
}

func (this *sqliteDir) GetCommand(name string) shell.ShellProcessor {
	return nil
}

func (this *sqliteDir) GetDir(name string) shell.ShellDir {
	if _, ok := this.sqlite.databases[name]; ok {
		return &databaseDir{this.sqlite, name}
	}
	return nil
}

func (this *sqliteDir) List() []string {
	r := make([]string, 0)
	for k, _ := range this.sqlite.databases {
		r = append(r, shell.DirName(k))
	}
	return r
}

type databaseDir struct {
	sqlite *SqliteServer
	name   string
}

func (this *databaseDir) Name() string {
	return this.name
}

func (this *databaseDir) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case "exec":
		return newExecCommand(this.sqlite, this.name)
	case "query":
		return newQueryCommand(this.sqlite, this.name)
	}
	return nil
}

func (this *databaseDir) GetDir(name string) shell.ShellDir {
	return nil
}

func (this *databaseDir) List() []string {
	r := []string{"exec", "query"}
	return r
}

func (this *databaseDir) DirInfo() string {
	if db, ok := this.sqlite.databases[this.name]; ok {
		return db.dsn
	}
	return ""
}

func (this *SqliteServer) NewShellDir() shell.ShellDir {
	return &sqliteDir{this}
}
