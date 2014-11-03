package vmmsql

import (
	"database/sql"
	"fmt"
	"logger"
	"strings"
	"sync"
)

type smartDBGlobal struct {
	lock     sync.RWMutex
	dbtables map[string][]string
}

func updateGlobal(dr, ds string, ts []string) {
	n := fmt.Sprintf("%s_%s", dr, ds)
	sdbGlobal.lock.Lock()
	defer sdbGlobal.lock.Unlock()
	sdbGlobal.dbtables[n] = ts
}

func queryGlobal(dr, ds string) []string {
	n := fmt.Sprintf("%s_%s", dr, ds)
	sdbGlobal.lock.RLock()
	defer sdbGlobal.lock.RUnlock()
	return sdbGlobal.dbtables[n]
}

var (
	sdbGlobal smartDBGlobal
)

func init() {
	sdbGlobal.dbtables = make(map[string][]string)
}

type dbInfo struct {
	Name         string
	Driver       string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	ReadOnly     bool
	Priority     int
}

func (this *dbInfo) String() string {
	return fmt.Sprintf("%s(%s)", this.Name, this.Driver)
}

type smartDB struct {
	lock    sync.RWMutex
	dbinfos []*dbInfo
	tables  map[string][]*dbInfo
}

func (this *smartDB) Update(dbi *dbInfo, tbs []string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, tb := range tbs {
		tb = strings.ToLower(tb)
		dbl, ok := this.tables[tb]
		if ok {
			dbl = append(dbl, dbi)
		} else {
			dbl = []*dbInfo{dbi}
		}
		this.tables[tb] = dbl
	}
}

func (this *smartDB) Add(dbi *dbInfo) {
	this.lock.Lock()
	this.dbinfos = append(this.dbinfos, dbi)
	this.lock.Unlock()

	tbs := queryGlobal(dbi.Driver, dbi.DataSource)
	this.Update(dbi, tbs)

	go func() {
		err := this.Refresh(dbi)
		if err != nil {
			logger.Warn(tag, "refresh %s fail - %s", dbi, err)
		}
	}()
}

func (this *smartDB) Remove(n string) {

}

func (this *smartDB) Select(tableName string, write bool) *dbInfo {
	tableName = strings.ToLower(tableName)
	this.lock.RLock()
	dbiList, ok := this.tables[tableName]
	this.lock.RUnlock()
	if !ok {
		return nil
	}
	var dbi *dbInfo
	pri := -1
	for _, o := range dbiList {
		if o.ReadOnly && write {
			continue
		}
		if o.Priority > pri {
			dbi = o
			pri = o.Priority
		}
	}
	return dbi
}

func (this *smartDB) Refresh(dbi *dbInfo) error {
	logger.Debug(tag, "refresh %s", dbi)
	db, err := sql.Open(dbi.Driver, dbi.DataSource)
	if err != nil {
		return err
	}
	defer db.Close()
	if dbi.Driver == "mysql" {
		rows, err1 := db.Query("show tables")
		if err1 != nil {
			return err1
		}
		defer rows.Close()
		tbs := make([]string, 0)
		for rows.Next() {
			var name string
			err2 := rows.Scan(&name)
			if err2 != nil {
				return err2
			}
			tbs = append(tbs, name)
		}
		logger.Debug(tag, "refresh %s end - %d", dbi, len(tbs))
		updateGlobal(dbi.Driver, dbi.DataSource, tbs)
		this.Update(dbi, tbs)
	} else {
		return fmt.Errorf("unknow database driver(%s) for SmartDB refresh", dbi.Driver)
	}
	return nil
}

func createSmartDB() (interface{}, error) {
	r := new(smartDB)
	r.dbinfos = make([]*dbInfo, 0)
	r.tables = make(map[string][]*dbInfo)
	return r, nil
}
