package clumem

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
	Services map[string]interface{}
	Remotes  map[string]map[string]interface{}
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
	if cfg.Remotes != nil {
	}
	if cfg.Services != nil {
	}
	return true
}

func (this *Service) storeRuntimeConfig(cfg *runtimeConfig) error {
	return this.database.StoreRuntimeConfig(tableName, 1, cfg)
}

func (this *Service) buildRuntimeConfig() *runtimeConfig {
	r := new(runtimeConfig)
	return r
}

func (this *Service) doSave() error {
	cfg := this.buildRuntimeConfig()
	return this.storeRuntimeConfig(cfg)
}
