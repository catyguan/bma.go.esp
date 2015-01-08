package boot

import (
	"os"
	"path/filepath"
	"time"
)

var (
	DevMode         bool
	Debug           bool
	WorkDir         string
	TempDir         string
	StartConfigFile string
	StartTime       time.Time
	LoadTime        time.Time
	Args            []string
)

func TempFile(fn string) (string, error) {
	f := filepath.Join(TempDir, fn)
	dir := filepath.Dir(f)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return f, nil
}

func WorkFile(fn string) (string, error) {
	f := filepath.Join(WorkDir, fn)
	dir := filepath.Dir(f)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return f, nil
}
