package eloader

import (
	// "logger"
	// "fmt"
	"app/cacheserver"
	"bytes"
	"errors"
	"fmt"
	"uprop"
)

var (
	loadeConfigId int32
)

type LoaderConfig struct {
	Name string
	Type string
	prop LoaderProperty
}

func (this *LoaderConfig) Desc() string {
	buf := bytes.NewBuffer([]byte{})
	if this.Name == "" {
		buf.WriteString("???")
	} else {
		buf.WriteString(this.Name)
	}
	buf.WriteString(" - ")
	if this.Type == "" {
		buf.WriteString("???")
	} else {
		buf.WriteString(this.Type)
	}
	return buf.String()
}

func (this *LoaderConfig) Valid() error {
	if this.Name == "" {
		return errors.New("loader name is empty")
	}
	if this.Type == "" {
		return fmt.Errorf("loader(%s) type empty", this.Name)
	}
	if this.prop != nil {
		err := this.prop.Valid()
		if err != nil {
			return fmt.Errorf("loader(%s) - %s", this.Name, err.Error())
		}
	}
	return nil
}

func (this *LoaderConfig) GetProperties() []*uprop.UProperty {
	b := new(uprop.UPropertyBuilder)
	ltypes := bytes.NewBuffer([]byte{})
	for k, _ := range loaderProviders {
		if ltypes.Len() > 0 {
			ltypes.WriteString(", ")
		}
		ltypes.WriteString(k)
	}
	b.NewProp("name", "loader name").Optional(false).BeValue(this.Name, func(v string) error {
		this.Name = v
		return nil
	})
	b.NewProp("type", "loader type, ["+ltypes.String()+"]").Optional(false).BeValue(this.Type, func(v string) error {
		p := GetLoaderProvider(v)
		if p == nil {
			return fmt.Errorf("invalid loader type '%s'", v)
		}
		this.Type = v
		this.prop = p.CreateProperty()
		return nil
	})
	if this.prop != nil {
		b.Merge(this.prop.GetUProperties())
	}
	return b.AsList()
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
	Valid() error
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
