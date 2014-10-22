package fileloader

import (
	"bmautil/valutil"
	"fmt"
	"sync"
)

const (
	VAR_F = "$F"
	VAR_M = "$M"
)

type FileLoader interface {
	Load(script string) ([]byte, error)
}

type FileLoaderFactory interface {
	Valid(cfg map[string]interface{}) error
	Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool)
	Create(cfg map[string]interface{}) (FileLoader, error)
}

var (
	flslibs map[string]FileLoaderFactory = make(map[string]FileLoaderFactory)
)

func AddFileLoaderFactory(n string, fac FileLoaderFactory) {
	flslibs[n] = fac
}

func GetFileLoaderFactory(n string) FileLoaderFactory {
	return flslibs[n]
}

func GetFileLoaderFactoryByType(cfg map[string]interface{}) (FileLoaderFactory, string, error) {
	xt, ok := cfg["Type"]
	if !ok {
		return nil, "", fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, "", fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := GetFileLoaderFactory(vxt)
	if fac == nil {
		return nil, "", fmt.Errorf("invalid FileLoader Type(%s)", xt)
	}
	return fac, vxt, nil
}

var (
	mflLock sync.RWMutex
	mfllibs map[string]FileLoader = make(map[string]FileLoader)
)

func SetModuleFileLoader(n string, fl FileLoader) {
	mflLock.Lock()
	defer mflLock.Unlock()
	mfllibs[n] = fl
}

func GetModuleFileLoader(n string) FileLoader {
	mflLock.RLock()
	defer mflLock.RUnlock()
	return mfllibs[n]
}

func RemoveModuleFileLoader(n string) {
	mflLock.Lock()
	defer mflLock.Unlock()
	delete(mfllibs, n)
}
