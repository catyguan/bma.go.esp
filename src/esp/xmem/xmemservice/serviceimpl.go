package xmemservice

import (
	"bmautil/binlog"
	"bmautil/sqlutil"
	"bmautil/valutil"
	"bytes"
	"database/sql"
	"esp/sqlite"
	"esp/xmem/xmemprot"
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
	for _, si := range this.memgroups {
		if si.group.blservice != nil {
			si.group.blservice.Stop()
		}
	}
	for _, si := range this.memgroups {
		if si.group.blservice != nil {
			si.group.blservice.WaitStop()
		}
	}
	for _, si := range this.memgroups {
		si.group.root.Clear()
	}
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

	if cfg.IsEnableBinlog() {
		this.doStartBinlog(name, mg, cfg)
	}

	return item, nil
}

func (this *Service) doUpdateMemGroupConfig(name string, cfg *MemGroupConfig) error {
	item, err := this.doGetGroup(name)
	if err != nil {
		return err
	}

	if item.config.IsEnableBinlog() {
		doStop := false
		if !cfg.IsEnableBinlog() {
			doStop = true
		} else {
			if item.config.BLConfig.FileName != cfg.BLConfig.FileName {
				doStop = true
			} else {
				doStop = item.config.BLConfig.Readonly != cfg.BLConfig.Readonly
			}
		}

		if doStop {
			err := this.doStopBinlog(name, item.group)
			if err != nil {
				logger.Debug(tag, "stop binlog fail - %s", err)
				return err
			}
		}
	}

	item.config = cfg

	if cfg.IsEnableBinlog() && item.group.blservice == nil {
		err := this.doStartBinlog(name, item.group, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Service) doEnableMemGroup(prof *MemGroupProfile) error {
	if err := prof.Valid(); err != nil {
		return err
	}

	item, ok := this.memgroups[prof.Name]
	if !ok {
		cfg := new(MemGroupConfig)
		item, _ = this.doCreateMemGroup(prof.Name, cfg)
	}
	if item.profile != nil {
		return fmt.Errorf("memory group '%s' already enable", prof.Name)
	}
	item.profile = prof

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

func (this *Service) doExecMemEncode(name string, mg *localMemGroup, coder XMemCoder) ([]byte, error) {
	gss, err2 := mg.Snapshot(coder)
	if err2 != nil {
		return nil, err2
	}
	return gss.Encode()
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
	gss := new(XMemGroupSnapshot)
	err := gss.Decode(content)
	if err != nil {
		return err
	}
	return mg.BuildFromSnapshot(coder, gss)
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

func (this *Service) doSetOp(group string, key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer, isAbsent bool) (xmemprot.MemVer, error) {
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "setOp(%s, %v, %v, %d, %s, %v)", group, key, val, sz, ver, isAbsent)
	}
	si, err := this.doGetGroup(group)
	if err != nil {
		return xmemprot.VERSION_INVALID, err
	}
	nver := xmemprot.VERSION_INVALID
	if ver == xmemprot.VERSION_INVALID {
		if isAbsent {
			nver = si.group.SetIfAbsent(key, val, sz)
		} else {
			nver = si.group.Set(key, val, sz)
		}
	} else {
		nver = si.group.CompareAndSet(key, val, sz, ver)
	}

	if si.config.IsBinlogWrite() {
		// binlog
		bl := new(XMemBinlog)
		bl.Op = OP_SET
		bl.Key = key.ToString()
		bl.Value = val
		bl.Version = ver
		bl.IsAbsent = isAbsent
		this.doWriteBinlog(group, si, bl)
	}

	return nver, nil
}

func (this *Service) doDeleteOp(group string, key xmemprot.MemKey, ver xmemprot.MemVer) (bool, error) {
	if logger.EnableDebug(tag) {
		logger.Debug(tag, "deleteOp(%s, %v, %s)", group, key, ver)
	}
	si, err := this.doGetGroup(group)
	if err != nil {
		return false, err
	}
	nver := si.group.CompareAndDelete(key, ver)
	if si.config.IsBinlogWrite() {
		// binlog
		bl := new(XMemBinlog)
		bl.Op = OP_DELETE
		bl.Key = key.ToString()
		bl.Version = ver
		this.doWriteBinlog(group, si, bl)
	}
	return nver != xmemprot.VERSION_INVALID, nil
}

func (this *Service) doSlaveJoin(g string, ver binlog.BinlogVer, lis binlog.Listener) (*binlog.Reader, error) {
	si, err := this.doGetGroup(g)
	if err != nil {
		return nil, err
	}
	if !si.config.IsEnableBinlog() {
		return nil, fmt.Errorf("'%s' binlog disable", g)
	}
	if si.profile == nil {
		return nil, fmt.Errorf("'%s' no profile", g)
	}
	if si.group.blservice == nil {
		return nil, fmt.Errorf("'%s' binlog not start", g)
	}
	logger.Info(tag, "'%s' slave join %d", g, ver)
	rd, err2 := si.group.blservice.NewReader()
	if err2 != nil {
		logger.Warn(tag, "'%s' slave join fail - %s", g, err2)
		return nil, err2
	}
	if !rd.SeekAndListen(ver, lis) {
		rd.Close()
		return nil, logger.Warn(tag, "'%s' slave join fail - seek fail", g)
	}
	return rd, nil
}
