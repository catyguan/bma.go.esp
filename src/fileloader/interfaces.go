package fileloader

import (
	"bmautil/valutil"
	"fmt"
	"strings"
	"sync"
)

const (
	VAR_F    = "$F"
	VAR_NO_F = "$!F"
	VAR_M    = "$M"
	VAR_NO_M = "$!M"
)

type FileLoader interface {
	Load(script string) ([]byte, error)
	Check(script string) (uint64, error)
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
	fldefLock sync.RWMutex
	fldeflibs map[string]FileLoader = make(map[string]FileLoader)
)

func DefineFileLoader(n string, fl FileLoader) {
	fldefLock.Lock()
	defer fldefLock.Unlock()
	fldeflibs[n] = fl
}

func GetDefinedFileLoader(n string) FileLoader {
	fldefLock.RLock()
	defer fldefLock.RUnlock()
	return fldeflibs[n]
}

func UndefineFileLoader(n string) {
	fldefLock.Lock()
	defer fldefLock.Unlock()
	delete(fldeflibs, n)
}

func DoValid(cfg map[string]interface{}) error {
	fac, _, err := GetFileLoaderFactoryByType(cfg)
	if err != nil {
		return err
	}
	return fac.Valid(cfg)
}

func DoCompare(cfg map[string]interface{}, old map[string]interface{}) bool {
	fac1, xt1, err1 := GetFileLoaderFactoryByType(cfg)
	if err1 != nil {
		return false
	}
	_, xt2, err2 := GetFileLoaderFactoryByType(old)
	if err2 != nil {
		return false
	}
	if xt1 != xt2 {
		return false
	}
	return fac1.Compare(cfg, old)
}

func DoCreate(cfg map[string]interface{}) (FileLoader, error) {
	fac, _, err := GetFileLoaderFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	return fac.Create(cfg)
}

func SplitModuleScript(n string) (string, string) {
	ns := strings.SplitN(n, ":", 2)
	if len(ns) == 1 {
		return "", n
	} else {
		return ns[0], ns[1]
	}
}

func BuildModuleScript(m, n string) string {
	if m == "" {
		return n
	}
	return m + ":" + n
}
