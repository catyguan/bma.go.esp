package config

import (
	"bmautil/valutil"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ConfigVar func(name string) (string, bool)

type ConfigObject map[string]interface{}

var (
	Global       ConfigObject
	AppConfigVar ConfigVar
)

type context struct {
	filename string
	cvar     ConfigVar
	vars     map[string]string
}

func include(file string, ctx context) (interface{}, error) {
	fname := file
	if !filepath.IsAbs(fname) {
		fname = filepath.Join(filepath.Dir(ctx.filename), file)
	}
	fmt.Printf("config '%s' including\n", fname)
	nctx := context{fname, ctx.cvar, ctx.vars}
	return loadAndParse(nctx)
}

func var2string(name string, ctx context) string {
	if v, ok := ctx.vars[name]; ok {
		return v
	}
	if ctx.cvar != nil {
		if v, ok := ctx.cvar(name); ok {
			return v
		}
	}
	if name == "CWD" {
		var wd, _ = os.Getwd()
		return wd
	}
	fmt.Printf("WARN: var '%s' not found", name)
	return ""
}

func parse(v interface{}, ctx context) (interface{}, error) {
	switch o := v.(type) {
	case []interface{}:
		for i, v2 := range o {
			nv, err := parse(v2, ctx)
			if err != nil {
				return v, err
			}
			o[i] = nv
		}
		return v, nil
	case map[string]interface{}:
		if ifilev := o["INCLUDE"]; ifilev != nil {
			if ifile, ok := ifilev.(string); ok {
				// include file
				return include(ifile, ctx)
			}
		}
		for k, v2 := range o {
			nv, err := parse(v2, ctx)
			if err != nil {
				return v, err
			}
			o[k] = nv
		}
		return v, nil
	case string:
		out := bytes.NewBuffer(make([]byte, 0))
		var c1 rune = 0
		word := bytes.NewBuffer(make([]byte, 0))
		for _, c := range []rune(o) {
			switch c1 {
			case 0:
				if c == '$' {
					c1 = c
				} else {
					out.WriteRune(c)
				}
			case '$':
				if c == '{' {
					c1 = '{'
				} else {
					out.WriteRune(c1)
					out.WriteRune(c)
					c1 = 0
				}
			case '{':
				if c == '}' {
					varname := word.String()
					word.Reset()
					nv := var2string(varname, ctx)
					out.WriteString(nv)
					c1 = 0
				} else {
					word.WriteRune(c)
				}
			}
		}
		switch c1 {
		case '$':
			out.WriteByte('$')
		case '{':
			out.WriteString("${")
			out.WriteString(word.String())
		}
		return out.String(), nil
	}
	return v, nil
}

func loadAndParse(ctx context) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(ctx.filename)
	if err != nil {
		fmt.Printf("ERROR: config '%s' load fail => %s\n", ctx.filename, err)
		return nil, err
	}
	var temp map[string]interface{}
	if err = json.Unmarshal(file, &temp); err != nil {
		fmt.Printf("ERROR: config '%s' format error => %s\n", ctx.filename, err)
		return nil, err
	}
	if varo, ok := temp["VAR"]; ok {
		if varm := valutil.ToStringMap(varo); varm != nil {
			for k, v := range varm {
				nv, err2 := parse(v, ctx)
				if err2 != nil {
					fmt.Printf("ERROR: %s=%v parse fail => %s\n", k, v, err)
					return nil, err
				}
				ctx.vars[k] = valutil.ToString(nv, "")
			}
		}
		delete(temp, "VAR")
	}
	r, err := parse(temp, ctx)
	return r.(map[string]interface{}), err
}

func InitConfig(fileName string) (ConfigObject, error) {
	return InitAndParseConfig(fileName, AppConfigVar)
}

func InitAndParseConfig(fileName string, cvar ConfigVar) (ConfigObject, error) {
	if fileName == "" {
		fileName = "esp-config.json"
	}
	fileName, _ = filepath.Abs(fileName)
	fmt.Printf("config '%s' loading\n", fileName)
	ctx := context{fileName, cvar, make(map[string]string)}
	temp, err := loadAndParse(ctx)
	if err != nil {
		return nil, err
	}
	return ConfigObject(temp), nil
}

func (this ConfigObject) GetBoolConfig(name string, defv bool) bool {
	if v, f := this.GetConfig(name); f {
		return valutil.ToBool(v, defv)
	}
	return defv
}

func (this ConfigObject) GetIntConfig(name string, defv int) int {
	if v, f := this.GetConfig(name); f {
		return valutil.ToInt(v, defv)
	}
	return defv
}

func (this ConfigObject) GetFloatConfig(name string, defv float64) float64 {
	if v, f := this.GetConfig(name); f {
		return valutil.ToFloat64(v, defv)
	}
	return defv
}

func (this ConfigObject) GetArrayConfig(name string) []interface{} {
	if v, f := this.GetConfig(name); f {
		return valutil.ToArray(v)
	}
	return nil
}

func (this ConfigObject) GetMapConfig(name string) map[string]interface{} {
	if v, f := this.GetConfig(name); f {
		return valutil.ToStringMap(v)
	}
	return nil
}

func (this ConfigObject) SubConfig(name string) ConfigObject {
	if v, f := this.GetConfig(name); f {
		o := valutil.ToStringMap(v)
		if o != nil {
			return ConfigObject(o)
		}
	}
	return nil
}

func (this ConfigObject) GetStringConfig(name string, defv string) string {
	if v, f := this.GetConfig(name); f {
		return valutil.ToString(v, defv)
	}
	return defv
}

func (this ConfigObject) GetBeanConfig(name string, bean interface{}) bool {
	m := this.GetMapConfig(name)
	if m == nil {
		return false
	}
	return valutil.ToBean(m, bean)
}

func (this ConfigObject) GetConfig(name string) (interface{}, bool) {
	nlist := strings.Split(name, ".")
	var thisv interface{}
	thisv = map[string]interface{}(this)
	for _, key := range nlist {
		switch thisv.(type) {
		case map[string]interface{}:
			cfgValue, f := thisv.(map[string]interface{})[key]
			if !f {
				return nil, false
			}
			thisv = cfgValue
		default:
			return nil, false
		}
	}
	return thisv, true
}
