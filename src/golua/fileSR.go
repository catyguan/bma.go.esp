package golua

import (
	"io/ioutil"
	"os"
	"path"
)

type FileGoSourceRepository struct {
	Dirs []string
}

func (this *FileGoSourceRepository) Load(script string, reload bool) (bool, string, error) {
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
