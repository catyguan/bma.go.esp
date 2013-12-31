package cacheserver

import (
	"bmautil/valutil"
	"testing"
)

type myconfig struct {
	A int
	LruCacheConfig
}

func TestConfigBean(t *testing.T) {

	cfg := new(myconfig)
	cfg.A = 1
	cfg.InvalidHolder = true
	r := valutil.BeanToMap(cfg)
	t.Errorf("%v", r)

	ncfg := new(myconfig)
	valutil.ToBean(r, ncfg)
	t.Errorf("%v", ncfg)
}
