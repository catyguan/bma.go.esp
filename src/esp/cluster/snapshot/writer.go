package snapshot

import (
	"esp/cluster/clusterbase"
	"fmt"
	"io/ioutil"
	"logger"
	"os"
)

type Writer struct {
	service   *Service
	ver       clusterbase.OpVer
	wfileName string
	wfd       *os.File
}

func (this *Writer) InitWriter(s *Service, ver clusterbase.OpVer) error {
	this.service = s
	this.ver = ver
	var err error
	this.wfd, err = ioutil.TempFile(s.config.DataDir, s.name+"_snapshot_")
	if err != nil {
		return err
	}
	this.wfileName = this.wfd.Name()
	logger.Debug(tag, "open snapshot temp file %s", this.wfileName)
	return nil
}

func (this *Writer) Write(p []byte) (n int, err error) {
	if this.wfd == nil {
		return 0, fmt.Errorf("writer close")
	}
	return this.wfd.Write(p)
}

func (this *Writer) Commit() error {
	if this.wfileName != "" {
		if this.wfd != nil {
			this.wfd.Close()
			this.wfd = nil
		}
		fn := this.service.config.GetFileName(this.ver)
		logger.Debug(tag, "commit snapshot file %s -> %s", this.wfileName, fn)
		os.Remove(fn)
		return os.Rename(this.wfileName, fn)
	}
	return nil
}

func (this *Writer) Close() {
	if this.wfd != nil {
		this.wfd.Close()
		this.wfd = nil
	}
	if this.wfileName != "" {
		logger.Debug(tag, "close snapshot temp file %s", this.wfileName)
		os.Remove(this.wfileName)
	}
}
