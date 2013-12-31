package namedsql

import (
	"boot"
	"config"
	"database/sql"
	"errors"
	"fmt"
	"logger"
	"sync"
)

const (
	tag = "NamedSQL"
)

const (
	actionSQL = iota
)

type SQLConfig struct {
	Name         string
	Driver       string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DelayOpen    bool
}

type configInfo struct {
	Database map[string]*SQLConfig
}

type sqlInfo struct {
	name   string
	db     *sql.DB
	config *SQLConfig
}

type NamedSQL struct {
	name      string
	databases map[string]*sqlInfo
	lock      sync.Mutex
}

func NewNamedSQL(name string) *NamedSQL {
	this := new(NamedSQL)
	this.name = name
	this.databases = make(map[string]*sqlInfo)
	return this
}

func (this *NamedSQL) Name() string {
	return this.name
}

func (this *NamedSQL) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.Database != nil {
			for key, scfg := range cfg.Database {
				if key == "" {
					logger.Error(tag, "database name empty")
					return false
				}
				scfg.Name = key
				if scfg.Driver == "" {
					logger.Error(tag, "database(%s) driver empty", key)
					return false
				}
				if scfg.DataSource == "" {
					logger.Error(tag, "database(%s) dataSource empty", key)
					return false
				}

				dobj := new(sqlInfo)
				dobj.name = key
				dobj.config = scfg
				this.databases[key] = dobj
			}
		}
	}
	return true
}

func (this *NamedSQL) open(dobj *sqlInfo) error {
	if dobj.db != nil {
		return nil
	}
	db, err := sql.Open(dobj.config.Driver, dobj.config.DataSource)
	if err != nil {
		logger.Error(tag, "open %s: %s, %s fail - %s", dobj.name, dobj.config.Driver, dobj.config.DataSource, err.Error())
		return err
	}
	db.SetMaxIdleConns(dobj.config.MaxIdleConns)
	db.SetMaxOpenConns(dobj.config.MaxOpenConns)
	dobj.db = db
	logger.Info(tag, "open %s: %s, %s", dobj.name, dobj.config.Driver, dobj.config.DataSource)

	return nil
}

func (this *NamedSQL) Start() bool {
	for _, dobj := range this.databases {
		if dobj.config.DelayOpen {
			logger.Info(tag, "'%s' skip open", dobj.name)
			continue
		}
		if this.open(dobj) != nil {
			return false
		}
	}
	return true
}

func (this *NamedSQL) Stop() bool {
	for _, dobj := range this.databases {
		dobj.db.Close()
		dobj.db = nil
	}
	return true
}

func (this *NamedSQL) DefaultBoot() {
	if this.name == "" {
		panic("SqliteServer name not set")
	}
	boot.Define(boot.INIT, this.name, this.Init)
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.CLEANUP, this.name, this.Stop)

	boot.Install(this.name, this)
}

func (this *NamedSQL) CreateSQL(cfg *SQLConfig) error {

	this.lock.Lock()
	defer this.lock.Unlock()

	name := cfg.Name
	dobj, ok := this.databases[name]
	if ok {
		return errors.New(fmt.Sprintf("sql[%s] exists", name))
	}

	logger.Info(tag, "create sql[%s, %s]", name, cfg.Driver)

	dobj = new(sqlInfo)
	dobj.config = cfg
	dobj.name = name
	if !cfg.DelayOpen {
		err := this.open(dobj)
		if err != nil {
			return err
		}
	}

	ndbs := make(map[string]*sqlInfo)
	for k, dobj := range this.databases {
		ndbs[k] = dobj
	}
	ndbs[name] = dobj
	this.databases = ndbs

	return nil
}

func (this *NamedSQL) CloseSQL(name string) error {
	dobj, ok := this.databases[name]
	if !ok {
		return nil
	}

	logger.Info(tag, "close sql[%s]", name)

	ndbs := make(map[string]*sqlInfo)
	for k, dobj := range this.databases {
		if k == name {
			continue
		}
		ndbs[k] = dobj
	}
	this.databases = ndbs

	dobj.db.Close()
	dobj.db = nil

	return nil
}

func (this *NamedSQL) Get(name string) (*sql.DB, error) {
	dobj, ok := this.databases[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("sql[%s] not exists", name))
	}
	if dobj.db != nil {
		return dobj.db, nil
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.open(dobj)
	if err != nil {
		return nil, err
	}
	return dobj.db, nil
}

func (this *NamedSQL) CheckConfig(name string) *SQLConfig {
	dobj, ok := this.databases[name]
	if ok {
		return dobj.config
	}
	return nil
}
