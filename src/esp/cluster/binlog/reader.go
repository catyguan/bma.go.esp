package binlog

import (
	"bmautil/qpushpull"
	"encoding/binary"
	"errors"
	"esp/cluster/clusterbase"
	"fmt"
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
	seeking bool
	puller  *qpushpull.QueuePushPull
}

func (this *Reader) initReader(fd *os.File, info *binlogInfo, ver clusterbase.OpVer) {
	this.info = info
	if this.rfd != nil {
		this.rfd.Close()
	}
	this.rfd = fd

	this.lastver = ver
	this.rbuffer = make([]byte, 8)
	this.seeking = true
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

func (this *Reader) doRead() (bool, clusterbase.OpVer, []byte, error) {
	if this.lastver == this.service.lastver {
		// end of binlog
		return false, 0, nil, nil
	}
	rver := this.lastver + 1
	if rver > this.info.lastVer {
		// next binlog file
		ok, err := this.service.doReaderOpen(this, rver)
		if err != nil {
			return false, 0, nil, err
		}
		if !ok {
			// end of binlog
			return false, 0, nil, nil
		}
	}

	l, err := this.doReadLength()
	if err != nil {
		return false, 0, nil, err
	}
	data := make([]byte, l)
	_, err2 := this.rfd.Read(data)
	if err2 != nil {
		return false, 0, nil, err2
	}
	ver, err3 := this.doReadVer()
	if err3 != nil {
		return false, 0, nil, err3
	}
	this.lastver = ver
	return true, ver, data, nil

}

func (this *Reader) Read() (bool, clusterbase.OpVer, []byte, error) {
	r := false
	var ver clusterbase.OpVer
	var data []byte
	err := this.service.executor.DoSync("read", func() error {
		if this.seeking {
			return fmt.Errorf("seeking, can't read")
		}
		var errv error
		r, ver, data, errv = this.doRead()
		return errv
	})
	return r, ver, data, err
}

func (this *Reader) doSeekFile(fver clusterbase.OpVer) (bool, error) {
	if this.info != nil {
		if this.lastver == this.service.lastver {
			// end of binlog
			return false, nil
		}
		if fver >= this.info.beginVer && fver <= this.info.lastVer {
			return true, nil
		}
	}
	// open correct binlog file
	ok, err := this.service.doReaderOpen(this, fver)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (this *Reader) doSeekProcess(fver clusterbase.OpVer, num int) (bool, error) {
	if fver == this.info.beginVer {
		return true, nil
	}
	for i := 0; i < num; i++ {
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
	return false, nil
}

func (this *Reader) goSeek(ver clusterbase.OpVer, ch chan error) {
	ok, err := this.doSeekProcess(ver, 10)
	if err != nil {
		ch <- err
		close(ch)
		return
	}
	if ok {
		this.seeking = false
		close(ch)
		return
	}
	// logger.Debug(tag, "seek more %d->%d", this.lastver, ver)
	go this.service.executor.DoNow("goseek", func() error {
		this.goSeek(ver, ch)
		return nil
	})
}

func (this *Reader) Seek(ver clusterbase.OpVer) error {
	ch := make(chan error, 1)
	err := this.service.executor.DoNow("seek", func() error {
		ok, err := this.doSeekFile(ver)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("version[%d] seek binlog file fail")
		}
		this.goSeek(ver, ch)
		return nil
	})
	if err != nil {
		return err
	}
	err = <-ch
	if err != nil {
		return err
	}
	return nil
}

type Listener func(ver clusterbase.OpVer, data []byte, closed bool)

func (this *Reader) Follow(ver clusterbase.OpVer, lis Listener) error {
	if this.listener != nil {
		return fmt.Errorf("already follow")
	}
	err := this.service.executor.DoSync("setlis", func() error {
		this.listener = lis
		if this.puller == nil {
			this.puller = qpushpull.NewQueuePushPull(8, this.qppHandler)
			this.puller.Run()
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = this.Seek(ver)
	if err != nil {
		return err
	}
	for {
		more := false
		err = this.service.executor.DoSync("pick", func() error {
			ok, ver, data, err := this.doRead()
			if err != nil {
				return err
			}
			if ok {
				more = true
				this.pushData(ver, data, false)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if !more {
			logger.Debug(tag, "reader attach writer '%d'", this.lastver)
			if this.rfd != nil {
				this.rfd.Close()
				this.rfd = nil
			}
			this.info = nil
			break
		}
	}
	return nil
}

type pushReq struct {
	ver    clusterbase.OpVer
	data   []byte
	closed bool
}

func (this *Reader) pushData(ver clusterbase.OpVer, data []byte, closed bool) {
	this.lastver = ver
	req := &pushReq{ver, data, closed}
	this.puller.Push(req)
}

func (this *Reader) qppHandler(req interface{}) {
	o := req.(*pushReq)
	this.listener(o.ver, o.data, o.closed)
}

func (this *Reader) Close() {
	this.service.executor.DoNow("closeReader", func() error {
		this.service.doCloseReader(this)
		return nil
	})
}
