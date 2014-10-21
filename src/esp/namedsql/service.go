package namedsql

import (
	"database/sql"
	"fmt"
	"logger"
	"sync"
)

const (
	tag = "namedsql"
)

type dbitem struct {
	db     *sql.DB
	config *sqlConfig
}

type Service struct {
	name   string
	config *configInfo
	mutex  sync.RWMutex
	dbs    map[string]*dbitem
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.dbs = make(map[string]*dbitem)
	return this
}

func (this *Service) open(k string, cfg *sqlConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DataSource)
	if err != nil {
		return nil, err
	}
	err2 := db.Ping()
	if err2 != nil {
		return nil, err2
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	dbi, ok := this.dbs[k]
	if !ok {
		db.Close()
		return nil, fmt.Errorf("invalid database[%s]", k)
	}
	if dbi.db != nil {
		db.Close()
		return dbi.db, nil
	}
	dbi.db = db
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	logger.Debug(tag, "open database[%s] done", k)
	return db, nil
}

func (this *Service) _create(k string, cfg *sqlConfig) {
	dbi := new(dbitem)
	dbi.config = cfg
	if !cfg.DelayOpen {
		go func() {
			_, err := this.open(k, cfg)
			if err != nil {
				logger.Warn(tag, "open database(%s) fail - %s", k, err)
			}
		}()
	}
	this.dbs[k] = dbi
}

func (this *Service) _remove(k string) {
	if dbi, ok := this.dbs[k]; ok {
		logger.Debug(tag, "close database(%s)", k)
		delete(this.dbs, k)
		if dbi.db != nil {
			dbi.db.Close()
			dbi.db = nil
		}
	}
}

func (this *Service) Get(name string) (*sql.DB, error) {
	this.mutex.RLock()
	dbi, ok := this.dbs[name]
	this.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("database[%s] not exists", name)
	}

	if dbi.db != nil {
		return dbi.db, nil
	}
	return this.open(name, dbi.config)
}
