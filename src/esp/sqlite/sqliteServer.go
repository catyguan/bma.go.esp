package sqlite

import (
	"bmautil/qexec"
	"bmautil/sqlutil"
	"boot"
	"config"
	"database/sql"
	// "fmt"
	"logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tag = "SqliteServer"
)

const (
	actionSQL = iota
)

type SqliteDatabaseInit func(name string, db *sql.DB) error

type actionInfo struct {
	actionType int
	name       string
	action     sqlutil.SQLAction
}

type sqliteDatabase struct {
	name          string
	db            *sql.DB
	dsn           string
	startOnAction bool
	closeOnError  bool
}

type SqliteServer struct {
	name      string
	databases map[string]*sqliteDatabase
	dbInit    []SqliteDatabaseInit
	executor  *qexec.QueueExecutor
}

type databaseConfig struct {
	Dsn           string
	StartOnAction bool
	CloseOnError  bool
}

type configInfo struct {
	QueueSize int
	Database  map[string]databaseConfig
}

func IsTableExists(db *sql.DB, tableName string) (bool, error) {
	sql := "SELECT count(*) as num FROM sqlite_master WHERE type='table' AND name=?;"
	var ret int
	err := db.QueryRow(sql, tableName).Scan(&ret)
	// fmt.Println(sql, "=>", ret, err)
	if err != nil {
		return false, err
	}
	return ret > 0, nil
}

func NewSqliteServer(name string) *SqliteServer {
	this := new(SqliteServer)
	this.name = name
	this.databases = make(map[string]*sqliteDatabase)
	this.dbInit = make([]SqliteDatabaseInit, 0)
	exec := qexec.NewQueueExecutor(tag, 0, this.requestHandler)
	exec.ErrorHandler = this.errorHandler
	exec.StopHandler = this.stopHandler
	this.executor = exec
	return this
}

func (this *SqliteServer) AddInit(dbInit SqliteDatabaseInit) {
	if dbInit != nil {
		this.dbInit = append(this.dbInit, dbInit)
	}
}

func (this *SqliteServer) Name() string {
	return this.name
}

func (this *SqliteServer) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		sz := cfg.QueueSize
		if sz <= 0 {
			sz = 32
		}
		this.executor.InitRequests(sz)
		if cfg.Database != nil {
			for key, info := range cfg.Database {
				if key == "" {
					logger.Error(tag, "database name empty")
					return false
				}
				if info.Dsn == "" {
					logger.Error(tag, "database(%s) dsn empty", key)
					return false
				}
				dobj := new(sqliteDatabase)
				dobj.name = key
				dobj.dsn = info.Dsn
				dobj.startOnAction = info.StartOnAction
				dobj.closeOnError = info.CloseOnError
				this.databases[key] = dobj
			}
		}
		if len(this.databases) == 0 {
			logger.Warn(tag, "database empty")
		}
		return true
	}
	logger.Error(tag, "GetBeanConfig(%s) fail", this.name)
	return false
}

func (this *SqliteServer) open(dobj *sqliteDatabase) error {
	logger.Debug(tag, "open '%s, %s'", dobj.name, dobj.dsn)
	db, err := sql.Open("sqlite3", dobj.dsn)
	if err != nil {
		logger.Error(tag, "open '%s, %s' fail - %s", dobj.name, dobj.dsn, err)
		return err
	}
	dobj.db = db

	for _, dbInit := range this.dbInit {
		action := new(actionInfo)
		action.actionType = actionSQL
		action.name = dobj.name
		action.action = func(db *sql.DB) error {
			return dbInit(dobj.name, db)
		}
		req := qexec.NewRequest("INIT", action, nil)
		_, err := this.executor.Execute(&req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *SqliteServer) Start() bool {
	for _, dobj := range this.databases {
		if dobj.startOnAction {
			logger.Info(tag, "'%s' skip open", dobj.name)
			continue
		}
		if this.open(dobj) != nil {
			return false
		}
	}
	this.executor.Run()
	return true
}

func (this *SqliteServer) errorHandler(req interface{}, err error) bool {
	action := req.(*actionInfo)
	if dobj, ok := this.databases[action.name]; ok {
		if dobj.closeOnError {
			logger.Info(tag, "close database '%s' after error", dobj.name)
			dobj.db.Close()
			dobj.db = nil
		}
	}
	return true
}

func (this *SqliteServer) stopHandler() {
	for _, dobj := range this.databases {
		if dobj.db != nil {
			logger.Debug(tag, "close database '%s'", dobj.name)
			dobj.db.Close()
			dobj.db = nil
		}
	}
}

func (this *SqliteServer) requestHandler(req interface{}) (bool, error) {
	action := req.(*actionInfo)
	switch action.actionType {
	case actionSQL:
		if dobj, ok := this.databases[action.name]; ok {
			if dobj.db == nil {
				if err := this.open(dobj); err != nil {
					return true, err
				}
			}
			return true, action.action(dobj.db)
		} else {
			err := logger.Error(tag, "action fail, database '%s' not exists", action.name)
			return true, err
		}
	}
	return true, nil
}

func (this *SqliteServer) Exec(name string, action sqlutil.SQLAction, cb func(err error)) error {
	act := new(actionInfo)
	act.actionType = actionSQL
	act.name = name
	act.action = action
	return this.executor.Do("exec", act, cb)
}

func (this *SqliteServer) Do(name string, action sqlutil.SQLAction, event chan error) error {
	return this.Exec(name, action, qexec.SyncCallback(event))
}

func (this *SqliteServer) Stop() bool {
	return this.executor.Stop()
}

func (this *SqliteServer) WaitStop() bool {
	return this.executor.WaitStop()
}

func (this *SqliteServer) DefaultBoot() {
	if this.name == "" {
		panic("SqliteServer name not set")
	}
	boot.Define(boot.INIT, this.name, this.Init)
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.STOP, this.name, this.Stop)
	boot.Define(boot.CLEANUP, this.name, this.WaitStop)

	boot.Install(this.name, this)
}

func InitTable(dbName string, tableName string, sqlstr []string) SqliteDatabaseInit {
	return func(name string, db *sql.DB) error {
		if dbName != "" && name != dbName {
			return nil
		}
		ok, err := IsTableExists(db, tableName)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		logger.Info(tag, "init create table '%s'", tableName)
		for _, str := range sqlstr {
			logger.Debug(tag, "init table %s => %s", tableName, str)
			_, err = db.Exec(str)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
