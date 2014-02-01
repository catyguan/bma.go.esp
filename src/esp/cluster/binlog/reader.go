package binlog

import (
	"encoding/binary"
	"errors"
	"esp/cluster/clusterbase"
	"logger"
	"os"
)

type Reader struct {
	service  *Service
	info     *binlogInfo
	rfd      *os.File
	listener Listener

	lastver clusterbase.OpVer
	rbuffer []byte
}

func (this *Reader) initReader(s *Service, fd *os.File, info *binlogInfo, ver clusterbase.OpVer) {
	this.service = s
	this.info = info
	if this.rfd != nil {
		this.rfd.Close()
	}
	this.rfd = fd

	this.lastver = ver
	this.rbuffer = make([]byte, 8)
}

func (this *Reader) doReadLength() (uint32, error) {
	fd := this.rfd
	if fd == nil {
		return 0, errors.New("closed")
	}
	_, err := this.rfd.Read(this.rbuffer[0:4])
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(this.rbuffer)
	return v, nil
}

func (this *Reader) doReadVer() (clusterbase.OpVer, error) {
	fd := this.rfd
	if fd == nil {
		return clusterbase.OpVer(0), errors.New("closed")
	}
	_, err := this.rfd.Read(this.rbuffer)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(this.rbuffer)
	return clusterbase.OpVer(v), nil
}

func (this *Reader) doRead() (clusterbase.OpVer, []byte, error) {
	l, err := this.doReadLength()
	if err != nil {
		return clusterbase.OpVer(0), nil, err
	}
	data := make([]byte, l)
	_, err2 := this.rfd.Read(data)
	if err2 != nil {
		return clusterbase.OpVer(0), nil, err2
	}
	ver, err3 := this.doReadVer()
	if err3 != nil {
		return clusterbase.OpVer(0), nil, err3
	}
	this.lastver = ver
	return ver, data, nil
}

func (this *Reader) Read() (bool, clusterbase.OpVer, []byte, error) {
	r := false
	var ver clusterbase.OpVer
	var data []byte
	err := this.service.executor.DoSync("read", func() error {
		if this.lastver == this.service.lastver {
			// end of binlog
			return nil
		}
		var err error
		for {
			rver := this.lastver + 1
			if rver <= this.info.lastVer {
				r = true
				ver, data, err = this.doRead()
				break
			}
			// next binlog file
			ok := false
			ok, err = this.service.doReaderOpen(this, rver)
			if err != nil {
				break
			}
			if !ok {
				// end of binlog
				break
			}
		}
		return err
	})
	return r, ver, data, err
}

func (this *Reader) Reset() {
	this.rfd.Seek(0, os.SEEK_SET)
}

func (this *Reader) doSeek(fver clusterbase.OpVer) (bool, error) {
	if this.lastver == this.service.lastver {
		// end of binlog
		return false, nil
	}
	if fver == this.info.beginVer {
		return true, nil
	}
	if fver > this.info.lastVer {
		// open correct binlog file
		ok, err := this.service.doReaderOpen(this, fver)
		if err != nil {
			return false, err
		}
		if !ok {
			// end of binlog
			return false, nil
		}
	}
	for {
		l, err := this.doReadLength()
		if err != nil {
			return false, err
		}
		_, err = this.rfd.Seek(int64(l), os.SEEK_CUR)
		if err != nil {
			return false, err
		}
		ver, err2 := this.doReadVer()
		if err2 != nil {
			return false, err2
		}
		this.lastver = ver
		if ver >= fver {
			return true, nil
		}
	}
}

type Listener func(ver clusterbase.OpVer, data []byte, closed bool)

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
	// TODO
	// if this.attachWriter {
	// 	return
	// }
	ver, data, err := this.doRead()
	if err != nil {
		logger.Debug(tag, "peek read error - %s", err)
		return
	}
	if data == nil {
		// TODO
		// this.attachWriter = true
		// this.lastver = this.service.lastver
		return
	}
	this.listener(ver, data, false)
	go this.Peek()
}

func (this *Reader) SeekAndListen(ver clusterbase.OpVer, lis Listener) bool {
	// _, err := this.Seek(ver)
	// if err != nil {
	// 	return false
	// }
	return this.SetListener(lis)
}

func (this *Reader) Close() {
	this.service.executor.DoNow("closeReader", func() error {
		this.service.doCloseReader(this)
		return nil
	})
}
