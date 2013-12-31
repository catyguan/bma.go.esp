package sqlite

import (
	"boot"
	"config"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSqliteServer(t *testing.T) {
	cfile := "../../../test/sqlite-config.json"

	config.AppConfigVar = func(name string) (string, bool) {
		if name == "CWD" {
			var wd, _ = os.Getwd()
			return wd + "/../../../bin", true
		}
		return "", false
	}

	sqliteServer := NewSqliteServer("sqliteServer")
	sqlstr := "create table foo (id integer not null primary key, name text);"
	sqliteServer.AddInit(InitTable("", "foo", []string{sqlstr}))

	if sqliteServer != nil {
		boot.Define(boot.STOP, "test1", func() bool {
			// do after server stop
			event := make(chan error, 1)
			sqliteServer.Do("test", func(db *sql.DB) error { return nil }, event)
			err := <-event
			if err != nil {
				fmt.Println("stop error", err)
			}
			return true
		})
	}

	sqliteServer.DefaultBoot()

	if sqliteServer != nil {
		act1 := func(db *sql.DB) error {
			rows, err := db.Query("select id, name from foo1")
			if err != nil {
				return err
			}
			defer rows.Close()
			return nil
		}
		act2 := func(db *sql.DB) error {
			rows, err := db.Query("select id, name from foo")
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id int
				var name string
				rows.Scan(&id, &name)
				fmt.Println("query =>", id, name)
			}
			return nil
		}

		boot.Define(boot.RUN, "test1", func() bool {
			event := make(chan error, 1)
			sqliteServer.Do("test", act1, event)
			err := <-event
			if err != nil {
				fmt.Println("do error", err)
			}
			sqliteServer.Do("test", act2, event)
			err = <-event
			if err != nil {
				fmt.Println("do error", err)
			}
			return true
		})
	}

	boot.Define(boot.RUN, "shutdonw", func() bool {
		time.Sleep(1 * time.Second)
		boot.Shutdown()
		return true
	})

	boot.Go(cfile)
}
