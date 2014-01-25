package binlog

import (
	"esp/cluster/clusterbase"
	"fmt"
	"logger"
)

type Writer struct {
	service *Service
}

func (this *Writer) Write(ver clusterbase.OpVer, bs []byte) bool {
	err := this.service.executor.DoNow("write", func() error {
		if ver <= this.service.lastver {
			return fmt.Errorf("invalid ver %d (lastver=%d)", ver, this.service.lastver)
		}
		this.service.doWrite(ver, bs)
		return nil
	})
	if err != nil {
		logger.Warn(tag, "'%s' write fail - %s", this.service.name, err)
		return false
	}
	return true
}

func (this *Writer) GerVersion() (clusterbase.OpVer, error) {
	rv := clusterbase.OpVer(0)
	err := this.service.executor.DoSync("getver", func() error {
		rv = this.service.lastver
		return nil
	})
	return rv, err
}
