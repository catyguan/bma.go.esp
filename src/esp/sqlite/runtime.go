package sqlite

import (
	"bmautil/sqlutil"
	"database/sql"
	"encoding/json"
	"fmt"
	"logger"
)

const (
	runtimeTableName = "tbl_runtime_data"
)

func (this *SqliteServer) InitRuntmeConfigTable() {
	str := fmt.Sprintf("create table %s (id text not null primary key, content text)", runtimeTableName)
	this.AddInit(InitTable("local", runtimeTableName, []string{str}))
}

func (this *SqliteServer) LoadRuntimeConfig(key string, cfgPtr interface{}) error {
	content := ""
	rowScan := func(rows *sql.Rows) error {
		if rows.Next() {
			return rows.Scan(&content)
		}
		return nil
	}
	sqlstr := fmt.Sprintf("SELECT content FROM %s WHERE id = ?", runtimeTableName)
	action := sqlutil.QueryAction(rowScan, sqlstr, key)
	event := make(chan error)
	defer close(event)
	this.Do("local", action, event)
	if err := <-event; err != nil {
		return logger.Error(tag, "load local data fail %s", err)
	}
	logger.Debug(tag, "load runtime config = %s", content)
	if content != "" {
		if err := json.Unmarshal([]byte(content), cfgPtr); err != nil {
			return logger.Error(tag, "runtime config parse error => %s", err)
		}
	}
	return nil
}

func (this *SqliteServer) StoreRuntimeConfig(key string, cfgPtr interface{}) error {
	data, err := json.Marshal(cfgPtr)
	if err != nil {
		logger.Error("ERROR: runtime config format error => %s", err)
		return err
	}

	sqlstr1 := fmt.Sprintf("DELETE FROM %s WHERE id = ?", runtimeTableName)
	action1 := sqlutil.ExecuteAction(nil, sqlstr1, key)
	err = this.Do("local", action1, nil)
	if err != nil {
		return err
	}

	f := func(r int64) {
		logger.Debug(tag, "store runtime config => %d", r)
	}
	content := string(data)
	logger.Debug(tag, "store runtime config = %s", content)
	sqlstr2 := fmt.Sprintf("INSERT INTO %s VALUES(?,?)", runtimeTableName)
	action2 := sqlutil.ExecuteAction(f, sqlstr2, key, content)
	logger.Info(tag, "do store runtime config")
	return this.Do("local", action2, nil)
}
