package objfac

import (
	"bmautil/valutil"
	"fmt"
)

type ObjectFactory interface {
	Valid(cfg map[string]interface{}, ofp ObjectFactoryProvider) error
	Compare(cfg map[string]interface{}, old map[string]interface{}, ofp ObjectFactoryProvider) (same bool)
	Create(cfg map[string]interface{}, ofp ObjectFactoryProvider) (interface{}, error)
}

type ObjectFactoryProvider func(kind, typ string) ObjectFactory

var (
	gOFS map[string]map[string]ObjectFactory
)

func GetObjectFactory(kind, typ string) ObjectFactory {
	if gOFS == nil {
		return nil
	}
	m := gOFS[kind]
	if m == nil {
		return nil
	}
	return m[typ]
}

func SetObjectFactory(kind, typ string, of ObjectFactory) {
	if gOFS == nil {
		gOFS = make(map[string]map[string]ObjectFactory)
	}
	m := gOFS[kind]
	if m == nil {
		m = make(map[string]ObjectFactory)
		gOFS[kind] = m
	}
	m[typ] = of
}

func QueryObjectFactory(kind string, cfg map[string]interface{}, ofp ObjectFactoryProvider) (ObjectFactory, string, error) {
	if ofp == nil {
		ofp = GetObjectFactory
	}
	xt, ok := cfg["Type"]
	if !ok {
		return nil, "", fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, "", fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := ofp(kind, vxt)
	if fac == nil {
		return nil, "", fmt.Errorf("invalid Object Type(%s, %s)", kind, xt)
	}
	return fac, vxt, nil
}

func DoValid(kind string, cfg map[string]interface{}, ofp ObjectFactoryProvider) error {
	fac, _, err := QueryObjectFactory(kind, cfg, ofp)
	if err != nil {
		return err
	}
	return fac.Valid(cfg, ofp)
}

func DoCompare(kind string, cfg map[string]interface{}, old map[string]interface{}, ofp ObjectFactoryProvider) bool {
	fac1, xt1, err1 := QueryObjectFactory(kind, cfg, ofp)
	if err1 != nil {
		return false
	}
	_, xt2, err2 := QueryObjectFactory(kind, old, ofp)
	if err2 != nil {
		return false
	}
	if xt1 != xt2 {
		return false
	}
	return fac1.Compare(cfg, old, ofp)
}

func DoCreate(kind string, cfg map[string]interface{}, ofp ObjectFactoryProvider) (interface{}, error) {
	fac, _, err := QueryObjectFactory(kind, cfg, ofp)
	if err != nil {
		return nil, err
	}
	return fac.Create(cfg, ofp)
}
