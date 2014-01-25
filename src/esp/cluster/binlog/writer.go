package binlog

import "logger"

type Writer struct {
	service *Service
}

func (this *Writer) Write(bs []byte) bool {
	err := this.service.executor.DoNow("write", func() error {
		this.service.doWrite(bs)
		return nil
	})
	if err != nil {
		logger.Warn(tag, "'%s' write fail - %s", this.service.name, err)
		return false
	}
	return true
}

func (this *Writer) WriteRetVer(bs []byte) (BinlogVer, error) {
	rv := BinlogVer(0)
	err := this.service.executor.DoSync("write", func() error {
		this.service.doWrite(bs)
		rv = this.service.seq
		return nil
	})
	return rv, err
}

func (this *Writer) GerVersion() (BinlogVer, error) {
	rv := BinlogVer(0)
	err := this.service.executor.DoSync("getver", func() error {
		rv = this.service.seq
		return nil
	})
	return rv, err
}
