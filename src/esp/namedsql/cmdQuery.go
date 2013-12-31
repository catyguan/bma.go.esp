package namedsql

import (
	"bmautil/sqlutil"
	"bmautil/valutil"
	"bytes"
	"database/sql"
	"esp/shell"
	"fmt"
)

type queryCommand struct {
	server *NamedSQL
	name   string
}

func newQueryCommand(server *NamedSQL, name string) *queryCommand {
	r := new(queryCommand)
	r.server = server
	r.name = name
	return r
}

func (this *queryCommand) Name() string {
	return "query"
}

func (this *queryCommand) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "[$$]sql"
	startPos := 0
	rowCount := 0
	fs := shell.NewFlagSet(name)
	fs.IntVar(&rowCount, "c", 0, "result show row count")
	fs.IntVar(&startPos, "p", 0, "result show start position")
	if shell.DoParse(s, command, fs, name, args, 1, 1) {
		return true
	}
	if startPos < 0 {
		s.Writeln(fmt.Sprintf("ERROR: start position invalid: %d", startPos))
		return true
	}
	if rowCount < 0 {
		s.Writeln(fmt.Sprintf("ERROR: row count invalid: %d", startPos))
		return true
	}
	if rowCount == 0 {
		rowCount = 10
	}
	sqlstr := fs.Arg(0)
	this.executeSQL(s, sqlstr, startPos, rowCount)
	return true
}

func (this *queryCommand) executeSQL(s *shell.Session, sqlstr string, pos, count int) {
	s.Writeln(fmt.Sprintf("[%s query]", this.name))
	s.Writeln(fmt.Sprintf("%s - %d,%d->", sqlstr, pos, count))

	defer func() {
		s.Writeln(fmt.Sprintf("[%s query end]", this.name))
	}()

	db, err := this.server.Get(this.name)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: get database(%s) fail %s", this.name, err.Error()))
		return
	}
	var result []map[string]interface{}
	act := func(db *sql.DB) error {
		rows, err := db.Query(sqlstr)
		if err != nil {
			return err
		}
		defer rows.Close()
		result, err = sqlutil.FetchMap(rows, pos, count)
		return err
	}
	err = act(db)
	if err != nil {
		s.Writeln(fmt.Sprintf("ERROR: %s", err))
		return
	}
	if result != nil {
		buf := bytes.NewBuffer(make([]byte, 0))
		for i, res := range result {
			if i > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%d: ", i+pos))
			for k, v := range res {
				buf.WriteString(fmt.Sprintf("%s=%s; ", k, valutil.ToString(v, "<null>")))
			}
		}
		s.Writeln(buf.String())
	}

}
