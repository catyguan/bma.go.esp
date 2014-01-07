package eloader

import (
	"app/cacheserver"
	"bmautil/valutil"
	"errors"
	"fmt"
	"logger"
	"math/rand"
	"time"
	"uprop"
)

type testProperties struct {
	// prop
	MinIdle int
	MaxIdle int
}

func (this *testProperties) Valid() error {
	if this.MaxIdle < this.MinIdle {
		return errors.New("min > max")
	}
	return nil
}

func (this *testProperties) GetUProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	p1 := uprop.NewUProperty("maxidle", this.MaxIdle, true, "max idle time, ms", func(v string) error {
		this.MaxIdle = valutil.ToInt(v, 0)
		return nil
	})
	r = append(r, p1)
	p2 := uprop.NewUProperty("minidle", this.MinIdle, true, "min idle time, ms", func(v string) error {
		this.MinIdle = valutil.ToInt(v, 0)
		return nil
	})
	r = append(r, p2)
	return r
}

func (this *testProperties) ToMap() map[string]interface{} {
	return valutil.BeanToMap(this)
}

func (this *testProperties) FromMap(vs map[string]interface{}) error {
	if !valutil.ToBean(vs, this) {
		return errors.New("config invalid")
	}
	return nil
}

type loaderProviderTEST struct {
}

func (this *loaderProviderTEST) Type() string {
	return "test"
}

func (this *loaderProviderTEST) CreateProperty() LoaderProperty {
	return new(testProperties)
}

func (this *loaderProviderTEST) CreateLoader(cfg *LoaderConfig, prop LoaderProperty) (Loader, error) {
	r := new(loaderTEST)
	r.prop = prop.(*testProperties)
	r.name = cfg.Name
	return r, nil
}

type loaderTEST struct {
	prop *testProperties
	name string
}

type loaderTaskTEST struct {
	timer *time.Timer
}

func (this *loaderTaskTEST) Cancel() {
	this.timer.Stop()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (this *loaderTEST) Load(cache *LoaderCache, req *cacheserver.GetRequest) LoadTask {
	f := func() {
		v := rand.Intn(999999-100000) + 100000
		val := fmt.Sprintf("%d", v)
		cache.LoadEnd(this.name, true, req.Key, []byte(val), nil, []string{"TEST: random value"})
	}
	sp := this.prop.MaxIdle - this.prop.MinIdle
	if sp < 0 {
		sp = 1000
	}
	wt := 0
	if sp > 0 {
		wt = rand.Intn(sp)
	}
	wt = wt + this.prop.MinIdle
	du := time.Duration(wt) * time.Millisecond
	timer := time.AfterFunc(du, f)
	logger.Debug(tag, "loaderTEST idle %s", du.String())
	return &loaderTaskTEST{timer}
}
