package sql4fileloader

import (
	"bmautil/valutil"
	"database/sql"
	"fileloader"
	"fmt"
	"runtime"
	"sync"
)

const (
	tag = "sqlfileloader"
)

func init() {
	fileloader.AddFileLoaderFactory("sql", SQLFileLoaderFactory)
}

type SQLFileLoader struct {
	db      *sql.DB
	table   string
	lock    sync.Mutex
	modules map[string]map[string]string
}

func (this *SQLFileLoader) Load(script string) ([]byte, error) {
	module, n := fileloader.SplitModuleScript(script)
	this.lock.Lock()
	defer this.lock.Unlock()
	mb, ok := this.modules[module]
	if !ok {
		rows, err := this.db.Query("SELECT script, content FROM "+this.table+" WHERE module=?", module)
		if err != nil {
			return nil, err
		}
		mb = make(map[string]string)
		var content string
		var name string
		for rows.Next() {
			err1 := rows.Scan(&name, &content)
			if err1 != nil {
				return nil, err1
			}
			mb[name] = content
		}
		this.modules[module] = mb
	}
	if str, ok2 := mb[n]; ok2 {
		return []byte(str), nil
	}
	return nil, nil
}

func (this *SQLFileLoader) Check(script string) (uint64, error) {
	return 0, nil
}

type config struct {
	Driver       string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	Table        string
}

type sqlFileLoaderFactory int

const (
	SQLFileLoaderFactory = sqlFileLoaderFactory(0)
)

func (this sqlFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	var co config
	if valutil.ToBean(cfg, &co) {
		if co.Driver == "" {
			return fmt.Errorf("Driver empty")
		}
		if co.DataSource == "" {
			return fmt.Errorf("DataSource empty")
		}
		return nil
	}
	return fmt.Errorf("invalid SQLFileLoader config")
}

func (this sqlFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo config
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.Driver != oo.Driver {
		return false
	}
	if co.DataSource != oo.DataSource {
		return false
	}
	if co.MaxOpenConns != oo.MaxOpenConns {
		return false
	}
	if co.MaxIdleConns != oo.MaxIdleConns {
		return false
	}
	if co.Table != oo.Table {
		return false
	}
	return true
}

func (this sqlFileLoaderFactory) Create(cfg map[string]interface{}) (fileloader.FileLoader, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co config
	valutil.ToBean(cfg, &co)
	db, err := sql.Open(co.Driver, co.DataSource)
	if err != nil {
		return nil, err
	}
	if co.MaxOpenConns > 0 {
		db.SetMaxOpenConns(co.MaxOpenConns)
	}
	if co.MaxIdleConns > 0 {
		db.SetMaxIdleConns(co.MaxIdleConns)
	}
	r := new(SQLFileLoader)
	r.db = db
	r.table = co.Table
	r.modules = make(map[string]map[string]string)
	if r.table == "" {
		r.table = "golua_sqlloader"
	}
	runtime.SetFinalizer(r, func(x *SQLFileLoader) {
		x.db.Close()
	})
	return r, nil
}
