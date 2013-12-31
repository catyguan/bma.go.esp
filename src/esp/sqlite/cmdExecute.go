package sqlite

import (
	"database/sql"
	"esp/shell"
	"fmt"
)

type execCommand struct {
	server *SqliteServer
	name   string
}

func newExecCommand(server *SqliteServer, name string) *execCommand {
	r := new(execCommand)
	r.server = server
	r.name = name
	return r
}

func (this *execCommand) Name() string {
	return "exec"
}

func (this *execCommand) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[$$]sql"
	returnId := false
	fs := shell.NewFlagSet(name)
	fs.BoolVar(&returnId, "id", false, "return last insert id")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}
	sqlstr := fs.Arg(0)
	this.executeSQL(s, sqlstr, returnId)
	return true
}

func (this *execCommand) executeSQL(s *shell.Session, sqlstr string, returnId bool) {
	s.Writeln(fmt.Sprintf("[%s exec]", this.name))
	s.Writeln(fmt.Sprintf("%s ->", sqlstr))

	var result int64
	done := make(chan error)
	act := func(db *sql.DB) error {
		res, err := db.Exec(sqlstr)
		if err != nil {
			return err
		}
		if returnId {
			result, err = res.LastInsertId()
			if err != nil {
				return err
			}
		} else {
			result, err = res.RowsAffected()
			if err != nil {
				return err
			}
		}
		return nil
	}
	this.server.Do(this.name, act, done)
	err := <-done
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: %s", err))
	} else {
		s.Writeln(fmt.Sprintf("%d", result))
	}
	s.Writeln(fmt.Sprintf("[%s exec end]", this.name))
}
