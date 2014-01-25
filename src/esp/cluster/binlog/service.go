package binlog

import (
	"bmautil/qexec"
	"encoding/binary"
	"esp/cluster/clusterbase"
	"logger"
	"os"
)

const (
	tag        = "binlog"
	headerSize = 8 + 4
)

// Service
type Service struct {
	name     string
	config   *BinlogConfig
	executor *qexec.QueueExecutor

	wfd     *os.File
	lastver clusterbase.OpVer
	opc     int
	wbuffer []byte
	readers map[*Reader]bool
}

func NewBinLog(n string, bufsize int, cfg *BinlogConfig) *Service {
	this := new(Service)
	this.name = n
	this.config = cfg
	this.executor = qexec.NewQueueExecutor(n, bufsize, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.wbuffer = make([]byte, headerSize)
	this.readers = make(map[*Reader]bool)
	return this
}

func (this *Service) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *Service) stopHandler() {
	this.doClose()
}

func (this *Service) Run() bool {
	if this.executor.Run() {
		this.executor.DoNow("setup", this.doSetup)
		return true
	}
	return false
}

func (this *Service) Stop() bool {
	return this.executor.Stop()
}

func (this *Service) WaitStop() {
	if this.executor.IsRun() {
		this.executor.WaitStop()
	}
}

func (this *Service) doSetup() error {
	if !this.config.Readonly {
		err := this.doOpenWrite()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Service) doOpenWrite() error {
	// Open the log file
	fd, err := os.OpenFile(this.config.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	this.wfd = fd
	return nil
}

func (this *Service) doOpenRead() (*Reader, error) {
	// Open the log file
	fd, err := os.OpenFile(this.config.FileName, os.O_RDONLY, 0664)
	if err != nil {
		return nil, err
	}
	rd := new(Reader)
	rd.initReader(this, fd)
	this.readers[rd] = true
	return rd, nil
}

func (this *Service) doClose() {
	for rd, _ := range this.readers {
		this.doCloseReader(rd)
	}
	if this.wfd != nil {
		this.wfd.Close()
		this.wfd = nil
	}
}

func (this *Service) doCloseReader(rd *Reader) {
	delete(this.readers, rd)
	if rd.rfd != nil {
		rd.rfd.Close()
		rd.rfd = nil
	}
	if rd.listener != nil {
		rd.listener(0, nil, true)
		rd.listener = nil
	}
}

func (this *Service) doWrite(ver clusterbase.OpVer, bs []byte) {
	if this.wfd == nil {
		logger.Debug(tag, "'%s' binlog closed when write", this.name)
		return
	}
	if ver <= this.lastver {
		return
	}
	old := this.lastver
	this.lastver = ver
	l := uint32(len(bs))
	binary.BigEndian.PutUint64(this.wbuffer, uint64(ver))
	binary.BigEndian.PutUint32(this.wbuffer[8:], l)
	_, err := this.wfd.Write(this.wbuffer)
	if err != nil {
		logger.Warn(tag, "'%s' write '%s' fail - %s", this.name, this.config.FileName, err)
		return
	}
	_, err = this.wfd.Write(bs)
	if err != nil {
		logger.Warn(tag, "'%s' write '%s' fail - %s", this.name, this.config.FileName, err)
		return
	}
	if this.config.SyncOpNum > 0 {
		this.opc++
		if this.opc >= this.config.SyncOpNum {
			this.opc = 0
			this.wfd.Sync()
		}
	}

	// push to waiting reader
	for rd, _ := range this.readers {
		if rd.listener != nil && rd.lastver == old {
			rd.listener(ver, bs, false)
			rd.lastver = ver
		}
	}
}

func (this *Service) NewWriter() (*Writer, error) {
	w := new(Writer)
	w.service = this
	return w, nil
}

func (this *Service) NewReader() (*Reader, error) {
	var rd *Reader
	err := this.executor.DoSync("reader", func() error {
		o, err := this.doOpenRead()
		if err != nil {
			return err
		}
		rd = o
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rd, nil
}
