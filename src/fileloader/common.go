package fileloader

import (
	"bmautil/valutil"
	"fmt"
)

func init() {
	AddFileLoaderFactory("c", CFileLoaderFactory)
}

type cFileLoader struct {
	deny   []string
	accept []string
	loader []FileLoader
}

func (this *cFileLoader) check(m string) bool {
	if this.accept != nil {
		for _, s := range this.deny {
			if s == m {
				return true
			}
		}
	}
	if this.deny != nil {
		for _, s := range this.deny {
			if s == "*" {
				return false
			}
			if s == m {
				return false
			}
		}
	}
	return true
}

func (this *cFileLoader) Load(script string) ([]byte, error) {
	module, _ := SplitModuleScript(script)
	if !this.check(module) {
		return nil, nil
	}
	for _, fl := range this.loader {
		bs, err := fl.Load(script)
		if err != nil {
			return nil, err
		}
		if bs != nil {
			return bs, nil
		}
	}
	return nil, nil
}

func (this cFileLoader) Check(script string) (uint64, error) {
	module, _ := SplitModuleScript(script)
	if !this.check(module) {
		return 0, nil
	}
	for _, fl := range this.loader {
		c, err := fl.Check(script)
		if err != nil {
			return 0, err
		}
		if c > 0 {
			return c, nil
		}
	}
	return 0, nil
}

type cflConfig struct {
	Deny   []string
	Accept []string
	FL     []string
}

type cFileLoaderFactory int

const (
	CFileLoaderFactory = cFileLoaderFactory(0)
)

func (this cFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	var co cflConfig
	if valutil.ToBean(cfg, &co) {
		if len(co.FL) == 0 {
			return fmt.Errorf("FL empty")
		}
		return nil
	}
	return fmt.Errorf("invalid CFileLoader config")
}

func (this cFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo cflConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if true {
		if len(co.Deny) != len(oo.Deny) {
			return false
		}
		tmp := make(map[string]bool)
		for _, s := range oo.Deny {
			tmp[s] = true
		}
		for _, s := range co.Deny {
			if !tmp[s] {
				return false
			}
		}
	}
	if true {
		if len(co.Accept) != len(oo.Accept) {
			return false
		}
		tmp := make(map[string]bool)
		for _, s := range oo.Accept {
			tmp[s] = true
		}
		for _, s := range co.Accept {
			if !tmp[s] {
				return false
			}
		}
	}
	if true {
		if len(co.FL) != len(oo.FL) {
			return false
		}
		for i, s := range oo.FL {
			if oo.FL[i] != s {
				return false
			}
		}
	}
	return true
}

func (this cFileLoaderFactory) Create(cfg map[string]interface{}) (FileLoader, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co cflConfig
	valutil.ToBean(cfg, &co)

	r := new(cFileLoader)
	r.deny = co.Deny
	r.accept = co.Accept
	r.loader = make([]FileLoader, 0)
	for _, n := range co.FL {
		fl := GetDefinedFileLoader(n)
		if fl == nil {
			return nil, fmt.Errorf("DefinedFileLoader('%s') miss", n)
		}
		r.loader = append(r.loader, fl)
	}
	return r, nil
}
