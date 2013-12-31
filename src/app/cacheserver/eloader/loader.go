package eloader

import (
	// "logger"
	// "fmt"
	"app/cacheserver"
	"bytes"
	"fmt"
	"uprop"
)

type LoaderConfig struct {
	Name string
	Type string
	prop LoaderProperty
}

func (this *LoaderConfig) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(this.Name)
	buf.WriteString("/")
	buf.WriteString(this.Type)
	buf.WriteString(" : ")

	if this.prop != nil {
		props := this.prop.GetUProperties()
		for _, p := range props {
			buf.WriteString(fmt.Sprintf("%s=%v; ", p.Name, p.Value))
		}
	}

	return buf.String()
}

type Loader interface {
	Load(cache *LoaderCache, req *cacheserver.GetRequest) LoadTask
}

// LoaderProvider
type LoaderProvider interface {
	Type() string

	CreateProperty() LoaderProperty

	CreateLoader(cfg *LoaderConfig, prop LoaderProperty) (Loader, error)
}

type LoaderProperty interface {
	GetUProperties() []*uprop.UProperty

	ToMap() map[string]interface{}
	FromMap(vs map[string]interface{}) error
}

type LoadTask interface {
	Cancel()
}

var loaderProviders = make(map[string]LoaderProvider)

func RegLoaderProvider(p LoaderProvider) {
	if p == nil {
		panic("cacheserver: Register provider is nil")
	}
	name := p.Type()
	if _, dup := loaderProviders[name]; dup {
		panic("cacheserver: Register called twice for provider " + name)
	}
	loaderProviders[name] = p
}

func GetLoaderProvider(name string) LoaderProvider {
	if p, ok := loaderProviders[name]; ok {
		return p
	}
	return nil
}

// NONE
type loaderProviderNONE struct {
}

func (this *loaderProviderNONE) Type() string {
	return "none"
}

func (this *loaderProviderNONE) CreateProperty() LoaderProperty {
	return nil
}

func (this *loaderProviderNONE) CreateLoader(cfg *LoaderConfig, prop LoaderProperty) (Loader, error) {
	r := new(loaderNONE)
	r.name = cfg.Name
	return r, nil
}

type loaderNONE struct {
	name string
}

func (this *loaderNONE) Load(cache *LoaderCache, req *cacheserver.GetRequest) LoadTask {
	// DO nothing
	go func() {
		cache.LoadEnd(this.name, false, req.Key, nil, nil, []string{"NONE: do nothing"})
	}()
	return nil
}

func init() {
	RegLoaderProvider(&loaderProviderNONE{})
	RegLoaderProvider(&loaderProviderTEST{})
}
