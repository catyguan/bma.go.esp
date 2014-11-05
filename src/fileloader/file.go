package fileloader

import (
	"bmautil/valutil"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	AddFileLoaderFactory("file", FileFileLoaderFactory)
}

type FileFileLoader struct {
	Dirs []string
}

func (this *FileFileLoader) vfile(script string) (string, error) {
	module, n := SplitModuleScript(script)
	n = filepath.Clean("/" + n)[1:]
	for _, dir := range this.Dirs {
		fn := dir
		if strings.Contains(fn, VAR_M) {
			fn = strings.Replace(fn, VAR_M, module, -1)
		} else {
			if module != "" {
				fn = path.Join(fn, module)
			}
		}
		if strings.Contains(fn, VAR_F) {
			fn = strings.Replace(fn, VAR_F, n, -1)
		} else {
			fn = path.Join(fn, n)
		}
		var err0 error
		fn, err0 = filepath.Abs(fn)
		if err0 != nil {
			return "", err0
		}
		// logger.Debug("file.fl", "%s, %s -> %s", dir, script, fn)
		_, err := os.Stat(fn)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		return fn, nil
	}
	return "", nil
}

func (this *FileFileLoader) Load(script string) ([]byte, error) {
	fn, err := this.vfile(script)
	if err != nil {
		return nil, err
	}
	content, err2 := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err2
	}
	return content, nil
}

func (this *FileFileLoader) Check(script string) (uint64, error) {
	fn, err := this.vfile(script)
	if err != nil {
		return 0, err
	}
	info, err2 := os.Stat(fn)
	if err2 != nil {
		return 0, err2
	}
	return uint64(info.ModTime().Unix()), nil
}

type fflConfig struct {
	Dirs []string
}

type fileFileLoaderFactory int

const (
	FileFileLoaderFactory = fileFileLoaderFactory(0)
)

func (this fileFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	var co fflConfig
	if valutil.ToBean(cfg, &co) {
		if len(co.Dirs) == 0 {
			return fmt.Errorf("Dirs empty")
		}
		for _, dir := range co.Dirs {
			if dir == "" {
				continue
			}
			if strings.Contains(dir, VAR_M) {
				continue
			}
			if strings.Contains(dir, VAR_F) {
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
	return fmt.Errorf("invalid FileFileLoader config")
}

func (this fileFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo fflConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if len(co.Dirs) != len(oo.Dirs) {
		return false
	}
	if true {
		tmp := make(map[string]bool)
		for _, dir := range oo.Dirs {
			tmp[dir] = true
		}
		for _, dir := range co.Dirs {
			if !tmp[dir] {
				return false
			}
		}
	}

	return true
}

func (this fileFileLoaderFactory) Create(cfg map[string]interface{}) (FileLoader, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co fflConfig
	valutil.ToBean(cfg, &co)
	r := new(FileFileLoader)
	r.Dirs = make([]string, 0, len(co.Dirs))
	for _, dir := range co.Dirs {
		dir = filepath.Clean(dir)
		r.Dirs = append(r.Dirs, dir)
	}
	return r, nil
}
