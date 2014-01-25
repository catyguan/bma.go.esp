package binlog

import (
	"encoding/binary"
	"errors"
	"io"
	"logger"
	"os"
)

type Reader struct {
	service      *Service
	rfd          *os.File
	listener     Listener
	attachWriter bool

	lastseq BinlogVer
	rbuffer []byte
	readed  int
	data    []byte
	remain  int
}

func (this *Reader) initReader(s *Service, fd *os.File) {
	this.service = s
	this.rfd = fd
	this.rbuffer = make([]byte, headerSize)
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

func (this *Reader) Reset() {
	this.rfd.Seek(0, os.SEEK_SET)
}

func (this *Reader) Seek(seq BinlogVer) (bool, error) {
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
	if this.attachWriter {
		return
	}
	seq, data, err := this.doRead()
	if err != nil {
		logger.Debug(tag, "peek read error - %s", err)
		return
	}
	if data == nil {
		this.attachWriter = true
		this.lastseq = this.service.seq
		return
	}
	this.listener(seq, data, false)
	go this.Peek()
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
