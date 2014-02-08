package snapshot

import (
	"io/ioutil"
	"logger"
	"os"
	"path/filepath"
	"testing"
)

func makeConfig() *SnapshotConfig {
	r := new(SnapshotConfig)
	wd, _ := os.Getwd()
	r.DataDir = filepath.Join(wd, "testdir")
	r.FileFormatter = "test_%04d.ssd"
	return r
}

func TestWriter(t *testing.T) {
	cfg := makeConfig()
	s := NewService("test", cfg)
	s.Setup()

	w, err1 := s.NewWriter(12)
	if err1 != nil {
		t.Error(err1)
		return
	}
	defer w.Close()

	w.Write([]byte("test1"))
	w.Write([]byte("test2"))
	err2 := w.Commit()
	if err2 != nil {
		t.Error(err2)
		return
	}

	logger.Debug("test", "end")
}

func TestReader(t *testing.T) {
	cfg := makeConfig()
	s := NewService("test", cfg)
	s.Setup()

	fd, err1 := s.NewReader(12, 0)
	if err1 != nil {
		t.Error(err1)
		return
	}
	defer fd.Close()

	b, err2 := ioutil.ReadAll(fd)
	if err2 != nil {
		t.Error(err2)
		return
	}
	logger.Debug("test", "data '%s'", string(b))
	logger.Debug("test", "end")
}
