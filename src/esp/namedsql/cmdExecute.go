package namedsql

import (
	"database/sql"
	"esp/shell"
	"fmt"
)

type execCommand struct {
	server *NamedSQL
	name   string
}

func newExecCommand(server *NamedSQL, name string) *execCommand {
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

	defer func() {
		s.Writeln(fmt.Sprintf("[%s exec end]", this.name))
	}()

	db, err := this.server.Get(this.name)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: get database(%s) fail %s", this.name, err.Error()))
		return
	}

	var result int64
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
	err = act(db)

	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: %s", err))
	} else {
		s.Writeln(fmt.Sprintf("%d", result))
	}

}
