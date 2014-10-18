package fileloader

import (
	"bmautil/valutil"
	"fmt"
)

type FileLoader interface {
	Load(script string) (bool, string, error)
}
type FileLoaderFactory interface {
	Valid(cfg map[string]interface{}) error
	Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool)
	Create(cfg map[string]interface{}) (FileLoader, error)
}

var (
	ssflibs map[string]FileLoaderFactory = make(map[string]FileLoaderFactory)
)

func AddFileLoaderFactory(n string, fac FileLoaderFactory) {
	ssflibs[n] = fac
}

func GetFileLoaderFactory(n string) FileLoaderFactory {
	return ssflibs[n]
}

func GetFileLoaderFactoryByType(cfg map[string]interface{}) (FileLoaderFactory, error) {
	xt, ok := cfg["Type"]
	if !ok {
		return nil, fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := GetFileLoaderFactory(vxt)
	if fac == nil {
		return nil, fmt.Errorf("invalid FileLoader Type(%s)", xt)
	}
	return fac, nil
}
