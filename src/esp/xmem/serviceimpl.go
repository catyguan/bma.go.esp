package xmem

import (
	"bmautil/byteutil"
	xcoder "bmautil/coder"
	"bmautil/sqlutil"
	"bmautil/valutil"
	"bytes"
	"database/sql"
	"esp/sqlite"
	"fmt"
	"io/ioutil"
	"logger"
)

const (
	tableName  = "tbl_xmem_service"
	tableName2 = "tbl_xmem_snapshot"
)

// impl
func (this *Service) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *Service) stopHandler() {

}

type runtimeConfig struct {
	Group map[string]interface{}
}

func (this *Service) initDatabase() {
	this.database.InitRuntmeConfigTable(tableName, []int{1})
	sqlstr := make([]string, 0)
	sqlstr = append(sqlstr, fmt.Sprintf("create table %s (name text not null primary key, content blob)", tableName2))
	this.database.AddInit(sqlite.InitTable("local", tableName2, sqlstr))
}

func (this *Service) loadRuntimeConfig() (*runtimeConfig, bool) {
	var cfg runtimeConfig
	err := this.database.LoadRuntimeConfig(tableName, 1, &cfg)
	if err != nil {
		return nil, false
	}
	return &cfg, true
}

func (this *Service) setupByConfig(cfg *runtimeConfig) bool {
	if cfg.Group != nil {
		for n, g := range cfg.Group {
			gcfg := new(MemGroupConfig)
			err := gcfg.FromMap(valutil.ToStringMap(g))
			if err != nil {
				logger.Warn(tag, "setup memory group '%s' fail - %s", n, err)
				if this.config.SafeMode {
					continue
				}
				return false
			}
			this.doCreateMemGroup(n, gcfg)
		}
	}
	return true
}

func (this *Service) storeRuntimeConfig(cfg *runtimeConfig) error {
	return this.database.StoreRuntimeConfig(tableName, 1, cfg)
}

func (this *Service) buildRuntimeConfig() *runtimeConfig {
	r := new(runtimeConfig)
	r.Group = make(map[string]interface{})
	for n, item := range this.memgroups {
		r.Group[n] = item.config.ToMap()
	}
	return r
}

func (this *Service) doSave() error {
	cfg := this.buildRuntimeConfig()
	return this.storeRuntimeConfig(cfg)
}

func (this *Service) doRun() error {
	return nil
}

func (this *Service) doListMemGroupName() []string {
	r := []string{}
	for k, _ := range this.memgroups {
		r = append(r, k)
	}
	return r
}

func (this *Service) doCreateMemGroup(name string, cfg *MemGroupConfig) (*serviceItem, error) {
	if _, ok := this.memgroups[name]; ok {
		return nil, fmt.Errorf("memory group '%s' already exists", name)
	}

	mg := newLocalMemGroup(name)
	item := new(serviceItem)
	item.config = cfg
	item.group = mg
	this.memgroups[name] = item

	return item, nil
}

func (this *Service) doUpdateMemGroupConfig(name string, cfg *MemGroupConfig) error {
	item, err := this.doGetGroup(name)
	if err != nil {
		return err
	}
	item.config = cfg
	return nil
}

func (this *Service) doEnableMemGroup(prof *MemGroupProfile) error {
	item, ok := this.memgroups[prof.Name]
	if !ok {
		cfg := new(MemGroupConfig)
		item, _ = this.doCreateMemGroup(prof.Name, cfg)
	}
	if item.profile != nil {
		return fmt.Errorf("memory group '%s' already enable", prof.Name)
	}
	item.profile = prof

	if prof.Coder == nil {
		item.config.NoSave = true
	}

	if !item.config.NoSave {
		err := this.doMemLoad(prof.Name, "", nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Service) doGetGroup(name string) (*serviceItem, error) {
	item, ok := this.memgroups[name]
	if !ok {
		return nil, fmt.Errorf("'%s' not exists", name)
	}
	return item, nil
}

func (this *Service) doMemSave(name string, fileName string, buf *bytes.Buffer) error {
	item, err := this.doGetGroup(name)
	if err != nil {
		return err
	}
	if item.profile == nil || item.profile.Coder == nil {
		return fmt.Errorf("'%s' no coder", name)
	}
	logger.Debug(tag, "doMemSave(%s,%s)", name, fileName)
	bs, err2 := this.doExecMemEncode(name, item.group, item.profile.Coder)
	if err2 != nil {
		return err2
	}
	if buf != nil {
		buf.Write(bs)
		return nil
	}
	if fileName != "" {
		logger.Debug(tag, "write '%s' snapshot to '%s'", name, fileName)
		return ioutil.WriteFile(fileName, bs, 0664)
	}
	if item.config.NoSave {
		return fmt.Errorf("'%s' disable save", name)
	}
	logger.Debug(tag, "write '%s' snapshot to local", name)
	return this.doSnapshotSave(name, bs)
}

func (this *Service) doExecMemSaveFile(fileName string, mg *localMemGroup, coder XMemCoder) error {
	slist, err2 := mg.Snapshot(coder)
	if err2 != nil {
		return err2
	}
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	xcoder.Int.DoEncode(w, len(slist))
	for _, s := range slist {
		s.Encode(w)
	}
	w.End()
	bs := buf.ToBytes()

	return ioutil.WriteFile(fileName, bs, 0644)
}

func (this *Service) doExecMemEncode(name string, mg *localMemGroup, coder XMemCoder) ([]byte, error) {
	slist, err2 := mg.Snapshot(coder)
	if err2 != nil {
		return nil, err2
	}
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	xcoder.Int.DoEncode(w, len(slist))
	for _, s := range slist {
		err := s.Encode(w)
		if err != nil {
			return nil, err
		}
	}
	return w.End().ToBytes(), nil
}

func (this *Service) doSnapshotSave(name string, bs []byte) error {
	logger.Debug(tag, "store memory snapshot = %s, %d", name, len(bs))

	delsql := fmt.Sprintf("DELETE FROM %s WHERE name = ?", tableName2)
	delact := sqlutil.ExecuteAction(nil, delsql, name)
	err3 := this.database.Do("local", delact, nil)
	if err3 != nil {
		return err3
	}

	inssql := fmt.Sprintf("INSERT INTO %s (name, content) VALUES(?,?)", tableName2)
	insact := sqlutil.ExecuteAction(nil, inssql, name, bs)
	return this.database.Do("local", insact, nil)
}

func (this *Service) doMemLoad(name string, fileName string, data []byte) error {
	item, err := this.doGetGroup(name)
	if err != nil {
		return err
	}
	if item.profile == nil || item.profile.Coder == nil {
		return fmt.Errorf("'%s' no coder", name)
	}
	if data == nil {
		if fileName != "" {
			logger.Debug(tag, "load '%s' snapshot from '%s'", name, fileName)
			data, err = ioutil.ReadFile(fileName)
		} else {
			logger.Debug(tag, "load '%s' snapshot from local", name)
			data, err = this.doSnapshotLoad(name)
		}
		if err != nil {
			return err
		}
	}
	return this.doExecMemDecode(name, item.group, item.profile.Coder, data)
}

func (this *Service) doSnapshotLoad(name string) ([]byte, error) {
	var content []byte
	rowScan := func(rows *sql.Rows) error {
		if rows.Next() {
			return rows.Scan(&content)
		}
		return nil
	}
	sqlstr := fmt.Sprintf("SELECT content FROM %s WHERE name = ?", tableName2)
	action := sqlutil.QueryAction(rowScan, sqlstr, name)
	event := make(chan error)
	defer close(event)
	this.database.Do("local", action, event)
	if err := <-event; err != nil {
		return nil, err
	}
	return content, nil
}

func (this *Service) doExecMemDecode(name string, mg *localMemGroup, coder XMemCoder, content []byte) error {
	logger.Debug(tag, "memory snapshot size = %d", len(content))
	if len(content) == 0 {
		logger.Debug(tag, "'%s' no snapshot", name)
		return nil
	}

	buf := byteutil.NewBytesBufferB(content)
	r := buf.NewReader()
	l, err1 := xcoder.Int.DoDecode(r)
	if err1 != nil {
		return err1
	}
	slist := []*XMemSnapshot{}
	for i := 0; i < l; i++ {
		ss, err2 := DecodeSnapshot(r)
		if err2 != nil {
			return err2
		}
		slist = append(slist, ss)
	}
	return mg.BuildFromSnapshot(coder, slist)
}

func (this *Service) doStoreAllMemGroup() error {
	for n, item := range this.memgroups {
		if !item.config.NoSave {
			err := this.doMemSave(n, "", nil)
			if err != nil {
				logger.Warn(tag, "store '%s' fail -%s", n, err)
			}
		}
	}
	return nil
}
