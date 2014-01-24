package binlog

import (
	"bmautil/byteutil"
	"bmautil/qexec"
	"encoding/binary"
	"errors"
	"io"
	"logger"
	"os"
	"time"
)

const (
	tag        = "binlog"
	headerSize = 8 + 4
)

type BinlogVer int64

type BinlogVerCoder int

func (this BinlogVerCoder) DoEncode(w *byteutil.BytesBufferWriter, v BinlogVer) {
	binary.Write(w, binary.BigEndian, int64(v))
}

func (this BinlogVerCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	this.DoEncode(w, v.(BinlogVer))
	return nil
}

func (this BinlogVerCoder) DoDecode(r *byteutil.BytesBufferReader) (BinlogVer, error) {
	var v BinlogVer
	err := binary.Read(r, binary.BigEndian, &v)
	return BinlogVer(v), err
}

func (this BinlogVerCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := this.DoDecode(r)
	return v, err
}

// Service
type Service struct {
	name     string
	config   *BinlogConfig
	executor *qexec.QueueExecutor

	wfd     *os.File
	seq     BinlogVer
	opc     int
	wbuffer []byte
	readers map[*Reader]bool
}

func NewBinLog(n string, bufsize int, cfg *BinlogConfig) *Service {
	this := new(Service)
	this.name = n
	this.config = cfg
	this.seq = BinlogVer(time.Now().UnixNano())
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
	rd.service = this
	rd.rfd = fd
	rd.rbuffer = make([]byte, headerSize)

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

func (this *Service) doWrite(bs []byte) {
	if this.wfd == nil {
		logger.Debug(tag, "'%s' binlog closed when write", this.name)
		return
	}
	old := this.seq
	this.seq++
	len := uint32(len(bs))
	binary.BigEndian.PutUint64(this.wbuffer, uint64(this.seq))
	binary.BigEndian.PutUint32(this.wbuffer[8:], len)
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
		if rd.listener != nil && rd.lastseq == old {
			rd.listener(this.seq, bs, false)
			rd.lastseq = this.seq
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

// Writer
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

func (this *Writer) GerVersion() (BinlogVer, error) {
	rv := BinlogVer(0)
	err := this.service.executor.DoSync("getver", func() error {
		rv = this.service.seq
		return nil
	})
	return rv, err
}

// Reader
type Reader struct {
	service  *Service
	rfd      *os.File
	listener Listener

	lastseq BinlogVer
	rbuffer []byte
	readed  int
	data    []byte
	remain  int
}

func (this *Reader) doReadHead() (bool, error) {
	fd := this.rfd
	if fd == nil {
		return false, errors.New("closed")
	}
	if this.readed >= headerSize {
		return true, nil
	}
	n, err := this.rfd.Read(this.rbuffer[this.readed:])
	this.readed += n
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (this *Reader) header() (BinlogVer, int) {
	tm := binary.BigEndian.Uint64(this.rbuffer)
	l := binary.BigEndian.Uint32(this.rbuffer[8:])
	return BinlogVer(tm), int(l)
}

func (this *Reader) doRead() (BinlogVer, []byte, error) {
	ok, err := this.doReadHead()
	if !ok {
		return 0, nil, err
	}
	tm, l := this.header()
	if this.data == nil {
		this.data = make([]byte, l)
		this.remain = int(l)
	}
	n, err2 := this.rfd.Read(this.data[l-this.remain:])
	this.remain += n
	if err2 != nil {
		if err2 == io.EOF {
			return 0, nil, nil
		}
		return 0, nil, err
	}
	data := this.data
	this.readed = 0
	this.data = nil
	return tm, data, nil
}

func (this *Reader) Read() (BinlogVer, []byte, error) {
	var seq BinlogVer
	var data []byte
	err := this.service.executor.DoSync("read", func() error {
		var err error
		seq, data, err = this.doRead()
		return err
	})
	return seq, data, err
}

func (this *Reader) doSeek(seq BinlogVer) (bool, error) {
	for {
		ok, err := this.doReadHead()
		if !ok {
			if err != nil {
				return false, err
			}
			return false, nil
		}
		tm, l := this.header()
		if err != nil {
			return false, err
		}
		if tm > seq {
			return false, nil
		}
		this.readed = 0
		_, err = this.rfd.Seek(int64(l), os.SEEK_CUR)
		if tm == seq {
			return true, nil
		}
	}
}

func (this *Reader) Seek(seq BinlogVer) (bool, error) {
	r := false
	var rerr error
	err2 := this.service.executor.DoSync("seek", func() error {
		ok, err := this.doSeek(seq)
		r = ok
		rerr = err
		return nil
	})
	if err2 != nil {
		rerr = err2
	}
	return r, rerr
}

type Listener func(seq BinlogVer, data []byte, closed bool)

func (this *Reader) SetListener(lis Listener) bool {
	if this.listener != nil {
		return false
	}
	this.service.executor.DoNow("setlis", func() error {
		this.listener = lis
		this.doPeek()
		return nil
	})
	return true
}

func (this *Reader) Peek() {
	if this.rfd != nil {
		this.service.executor.DoNow("peek", func() error {
			this.doPeek()
			return nil
		})
	}
}

func (this *Reader) doPeek() {
	if this.listener == nil {
		return
	}
	seq, data, err := this.doRead()
	if err != nil {
		logger.Debug(tag, "peek error - %s", err)
		return
	}
	if data == nil {
		this.lastseq = this.service.seq
		return
	}
	this.listener(seq, data, false)
	go this.Peek()
}

func (this *Reader) Reset() {
	this.service.executor.DoNow("reset", func() error {
		this.rfd.Seek(0, os.SEEK_SET)
		return nil
	})
}

func (this *Reader) SeekAndListen(seq BinlogVer, lis Listener) bool {
	_, err := this.Seek(seq)
	if err != nil {
		return false
	}
	return this.SetListener(lis)
}

func (this *Reader) Close() {
	this.service.executor.DoNow("closeReader", func() error {
		this.service.doCloseReader(this)
		return nil
	})
}
