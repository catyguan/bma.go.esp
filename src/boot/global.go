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
	err := os.MkdirAll(TempDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	f := filepath.Join(TempDir, fn)
	return f, nil
}

func WorkFile(fn string) (string, error) {
	err := os.MkdirAll(WorkDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	f := filepath.Join(WorkDir, fn)
	return f, nil
}
