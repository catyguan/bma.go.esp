package fileloader

import (
	"fmt"
	"strings"
)

type wrap struct {
	l FileLoader
}

func (this *wrap) ModuleLoad(module string, script string) ([]byte, error) {
	return this.l.Load(script)
}

func init() {
	AddFileLoaderFactory("m", MFileLoaderFactory)
}

func SplitModuleScript(n string) (string, string) {
	ns := strings.SplitN(n, ":", 2)
	if len(ns) == 1 {
		return "", n
	} else {
		return ns[0], ns[1]
	}
}

type mFileLoader int

func (this mFileLoader) Load(script string) ([]byte, error) {
	module, n := SplitModuleScript(script)
	var mfl FileLoader
	if module != "" {
		mfl = GetModuleFileLoader(module)
	}
	if mfl == nil {
		mfl = GetModuleFileLoader("*")
	}
	if mfl == nil {
		if module == "" {
			return nil, fmt.Errorf("FileLoader module empty - %s", script)
		}
		return nil, fmt.Errorf("FileLoader module(%s) invalid", module)
	}
	return mfl.Load(n)
}

type mFileLoaderFactory int

const (
	MFileLoaderFactory = mFileLoaderFactory(0)
)

func (this mFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	return nil
}

func (this mFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	return true
}

func (this mFileLoaderFactory) Create(cfg map[string]interface{}) (FileLoader, error) {
	return mFileLoader(0), nil
}
