package golua

import (
	"bmautil/valutil"
	"fmt"
)

var (
	ssflibs map[string]ScriptSourceFactory = make(map[string]ScriptSourceFactory)
)

func AddScriptSourceFactory(n string, fac ScriptSourceFactory) {
	ssflibs[n] = fac
}

func GetScriptSourceFactory(n string) ScriptSourceFactory {
	return ssflibs[n]
}

func GetScriptSourceFactoryByType(cfg map[string]interface{}) (ScriptSourceFactory, error) {
	xt, ok := cfg["Type"]
	if !ok {
		return nil, fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := GetScriptSourceFactory(vxt)
	if fac == nil {
		return nil, fmt.Errorf("invalid ScriptSource Type(%s)", xt)
	}
	return fac, nil
}

type CommonScriptSourceFactory int

func (this CommonScriptSourceFactory) Valid(cfg map[string]interface{}) error {
	fac, err := GetScriptSourceFactoryByType(cfg)
	if err != nil {
		return err
	}
	return fac.Valid(cfg)
}

func (this CommonScriptSourceFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	fac1, err1 := GetScriptSourceFactoryByType(cfg)
	if err1 != nil {
		return false
	}
	fac2, err2 := GetScriptSourceFactoryByType(old)
	if err2 != nil {
		return false
	}
	if fac1 != fac2 {
		return false
	}
	return fac1.Compare(cfg, old)
}

func (this CommonScriptSourceFactory) Create(cfg map[string]interface{}) (ScriptSource, error) {
	fac, err := GetScriptSourceFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	return fac.Create(cfg)
}
