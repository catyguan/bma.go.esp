package clumem

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
	this.database.InitRuntmeConfigTable(tableName, []int{1, 2})
}

func (this *Service) loadServiceConfig() bool {
	cfg := make(map[string]interface{})
	err := this.database.LoadRuntimeConfig(tableName, 2, cfg)
	if err != nil {
		return false
	}
	scfg := new(serviceConfig)
	err = scfg.FromMap(cfg)
	if err != nil {
		logger.Warn(tag, "load runtime serviceConfig fail - %s", err)
		return false
	}
	this.serviceConfig = scfg
	this.isServiceConfigLoaded = true
	return true
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
			this.doCreateMemGroup(gcfg)
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
	// connect to seed
	return nil
}

func (this *Service) doCreateMemGroup(cfg *MemGroupConfig) error {
	if _, ok := this.memgroups[cfg.Name]; ok {
		return fmt.Errorf("memory group '%s' already exists", cfg.Name)
	}

	mg := newLocalMemGroup(cfg.Name)
	item := new(serviceItem)
	item.config = cfg
	item.group = mg
	this.memgroups[cfg.Name] = item

	return nil
}
