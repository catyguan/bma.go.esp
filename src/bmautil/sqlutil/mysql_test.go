package sqlutil

import (
	"database/sql"
	// "github.com/ziutek/mymysql/godrv"
	"bmautil/valutil"
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestMysql(t *testing.T) {

	db, err := sql.Open("mysql", "root:root@tcp(172.19.16.195:3306)/db_live2")
	if err != nil {
		t.Error(err)
		return
	}

	sqlstr := "show tables"

	var result []map[string]interface{}
	act := func(db *sql.DB) error {
		rows, err := db.Query(sqlstr)
		if err != nil {
			return err
		}
		defer rows.Close()
		result, err = FetchMap(rows, 0, 10)
		return err
	}
	err = act(db)
	if err != nil {
		t.Error(err)
		return
	}
	if result != nil {
		buf := bytes.NewBuffer(make([]byte, 0))
		for i, res := range result {
			if i > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(fmt.Sprintf("%d: ", i+0))
			for k, v := range res {
				buf.WriteString(fmt.Sprintf("%s=%s; ", k, valutil.ToString(v, "<unknow>")))
			}
		}
		t.Error(buf.String())
	}

}
