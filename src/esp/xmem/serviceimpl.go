package xmem

import (
	"bmautil/valutil"
	"fmt"
	"logger"
)

const (
	tableName = "tbl_clumen_service"
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

func (this *Service) doEnableMemGroup(prof *memGroupProfile) error {
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
		err := this.doMemLoad(prof.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Service) doMemSave(name string) error {
	return nil
}

func (this *Service) doMemLoad(name string) error {
	return nil
}
