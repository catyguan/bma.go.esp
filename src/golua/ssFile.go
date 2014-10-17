package golua

import (
	"bmautil/valutil"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	AddScriptSourceFactory("file", FileScriptSourceFactory(0))
}

type FileScriptSource struct {
	Dirs []string
}

func (this *FileScriptSource) Load(script string) (bool, string, error) {
	for _, dir := range this.Dirs {
		fn := path.Join(dir, script)
		_, err := os.Stat(fn)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return false, "", err
		}
		content, err2 := ioutil.ReadFile(fn)
		if err != nil {
			return false, "", err2
		}
		return true, string(content), nil
	}
	return false, "", nil
}

type fssConfig struct {
	Dirs []string
}

type FileScriptSourceFactory int

func (this FileScriptSourceFactory) Valid(cfg map[string]interface{}) error {
	var co fssConfig
	if valutil.ToBean(cfg, &co) {
		if len(co.Dirs) == 0 {
			return fmt.Errorf("Dirs empty")
		}
		for _, dir := range co.Dirs {
			if dir == "" {
				continue
			}
			s, err := os.Stat(dir)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("'%s' not exists", dir)
				}
			}
			if !s.IsDir() {
				return fmt.Errorf("'%s' not dir", dir)
			}
		}
		return nil
	}
	return fmt.Errorf("invalid FileScriptSource config")
}

func (this FileScriptSourceFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo fssConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if len(co.Dirs) != len(oo.Dirs) {
		return false
	}
	tmp := make(map[string]bool)
	for _, dir := range oo.Dirs {
		tmp[dir] = true
	}
	for _, dir := range co.Dirs {
		if !tmp[dir] {
			return false
		}
	}
	return true
}

func (this FileScriptSourceFactory) Create(cfg map[string]interface{}) (ScriptSource, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co fssConfig
	valutil.ToBean(cfg, &co)
	r := new(FileScriptSource)
	r.Dirs = co.Dirs
	return r, nil
}
