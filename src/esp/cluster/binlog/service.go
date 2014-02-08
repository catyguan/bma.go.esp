package binlog

import (
	"bmautil/qexec"
	"bmautil/valutil"
	"bytes"
	"encoding/binary"
	"errors"
	"esp/cluster/clusterbase"
	"fmt"
	"io/ioutil"
	"logger"
	"os"
	"sort"
)

const (
	tag = "binlog"
)

type binlogInfo struct {
	num      int
	beginVer clusterbase.OpVer // ver>beginVer
	lastVer  clusterbase.OpVer // lastVer!=0 && ver<=lastVer
	fileSize int64
}

type sortInfo []*binlogInfo

func (a sortInfo) Len() int           { return len(a) }
func (a sortInfo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortInfo) Less(i, j int) bool { return a[i].num < a[j].num }

// Service
type Service struct {
	name     string
	config   *BinlogConfig
	executor *qexec.QueueExecutor

	infos     []*binlogInfo
	current   *binlogInfo
	wfileName string
	wfd       *os.File
	lastver   clusterbase.OpVer
	opc       int
	wbuffer   []byte
	readers   map[*Reader]bool
}

func NewBinLog(n string, bufsize int, cfg *BinlogConfig) *Service {
	this := new(Service)
	this.name = n
	this.config = cfg
	this.executor = qexec.NewQueueExecutor(n, bufsize, this.requestHandler)
	this.executor.StopHandler = this.stopHandler
	this.wbuffer = make([]byte, 8)
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

func (this *Service) String() string {
	return fmt.Sprintf("Binlog[%s, %d]", this.name, this.lastver)
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
	if true {
		err := this.doSetupInfo()
		if err != nil {
			return err
		}
	}
	if !this.config.Readonly {
		err := this.doOpenWrite(false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Service) doSetupInfo() error {
	if true {
		err := os.MkdirAll(this.config.LogDir, 0664)
		if err != nil {
			return err
		}
	}

	if true {
		fs, err := ioutil.ReadDir(this.config.LogDir)
		if err != nil {
			return err
		}
		this.infos = make([]*binlogInfo, 0)
		for _, f := range fs {
			if f.IsDir() {
				continue
			}
			var num int32
			n, err2 := fmt.Sscanf(f.Name(), this.config.FileFormatter, &num)
			if err2 != nil || n != 1 {
				continue
			}
			info := new(binlogInfo)
			info.num = int(num)
			info.beginVer = 0
			info.lastVer = 0
			info.fileSize = f.Size()
			this.infos = append(this.infos, info)
		}
		sort.Sort(sortInfo(this.infos))
	}

	if true {
		cver := clusterbase.OpVer(0)
		for _, info := range this.infos {
			info.beginVer = cver
			ver, err := this.doReadLastVer(info)
			if err != nil {
				return err
			}
			info.lastVer = ver
			cver = ver
		}
		this.lastver = cver
	}

	if logger.EnableDebug(tag) {
		buf := bytes.NewBuffer([]byte{})
		for i, info := range this.infos {
			if i != 0 {
				buf.WriteString(",")
			}
			sz := valutil.MakeSizeString(uint64(info.fileSize))
			buf.WriteString(fmt.Sprintf("%d[%d-%d,%s]", info.num, info.beginVer, info.lastVer, sz))
		}
		logger.Debug(tag, "%s setupInfo %s", this, buf.String())
	}
	return nil
}

func (this *Service) doReadLastVer(info *binlogInfo) (clusterbase.OpVer, error) {
	fd, err := os.OpenFile(this.config.GetFileName(info.num), os.O_RDONLY, 0664)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	_, err = fd.Seek(-8, 2)
	if err != nil {
		return 0, nil
	}
	b := make([]byte, 8)
	_, err = fd.Read(b)
	if err != nil {
		return 0, nil
	}
	v := binary.BigEndian.Uint64(b)
	return clusterbase.OpVer(v), nil
}

func (this *Service) doOpenWrite(forceNew bool) error {
	if this.wfd != nil {
		this.wfd.Close()
		this.wfd = nil
	}

	var info *binlogInfo
	if !forceNew {
		if len(this.infos) > 0 {
			info = this.infos[len(this.infos)-1]
		}
	}
	if info == nil {
		info = new(binlogInfo)
		info.num = 1
		if len(this.infos) > 0 {
			info.num = this.infos[len(this.infos)-1].num + 1
		}
		info.beginVer = this.lastver
		info.lastVer = 0
		fn := this.config.GetFileName(info.num)
		fi, err := os.Stat(fn)
		if err == nil {
			info.fileSize = fi.Size()
		} else {
			logger.Info(tag, "%s new binlog file '%s'", this, fn)
			info.fileSize = 0
		}
		this.infos = append(this.infos, info)
	}

	// Open the log file
	this.wfileName = this.config.GetFileName(info.num)
	logger.Debug(tag, "%s open binlog file %s", this, this.wfileName)
	fd, err := os.OpenFile(this.wfileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	this.wfd = fd
	this.current = info
	return nil
}

func (this *Service) doSelectOpen(ver clusterbase.OpVer, flag int) (*os.File, *binlogInfo, error) {
	// Seek log file
	var info *binlogInfo
	if ver == clusterbase.OpVer(0) {
		info = this.infos[0]
	} else {
		for i := len(this.infos) - 1; i >= 0; i-- {
			o := this.infos[i]
			if ver > o.beginVer {
				info = o
				break
			}
		}
	}
	if info == nil {
		return nil, nil, nil
	}

	// Open the log file
	fn := this.config.GetFileName(info.num)
	logger.Debug(tag, "%s select(%d) open(%s)", this, ver, fn)
	if flag == 0 {
		flag = os.O_RDONLY
	}
	fd, err := os.OpenFile(fn, flag, 0664)
	if err != nil {
		return nil, nil, err
	}
	return fd, info, nil
}

func (this *Service) doOpenRead(rd *Reader, ver clusterbase.OpVer) (bool, error) {
	// Seek log file
	fd, info, err := this.doSelectOpen(ver, os.O_RDONLY)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}
	rd.initReader(fd, info, info.beginVer)
	this.readers[rd] = true
	return true, nil
}

func (this *Service) doClose() {
	tmp := make([]*Reader, 0)
	for rd, _ := range this.readers {
		this.doCloseReader(rd)
		tmp = append(tmp, rd)
	}
	for _, rd := range tmp {
		if rd.puller != nil {
			rd.puller.WaitClose()
			rd.puller = nil
		}
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
	if rd.puller != nil {
		rd.pushData(0, nil, true)
		rd.puller.Close()
	}
}

func (this *Service) doWrite(ver clusterbase.OpVer, bs []byte) (bool, error) {
	if this.wfd == nil {
		return false, logger.Error(tag, "'%s' binlog closed when write", this.name)
	}
	if ver != this.lastver+1 {
		return false, nil
	}
	l := uint32(len(bs))
	total := int64(l + 8 + 4)
	if this.current.fileSize+total > this.config.FileMaxSize {
		err := this.doOpenWrite(true)
		if err != nil {
			return false, err
		}
	}
	old := this.lastver
	binary.BigEndian.PutUint32(this.wbuffer, l)
	_, err := this.wfd.Write(this.wbuffer[0:4])
	if err != nil {
		return false, logger.Error(tag, "'%s' write1 '%s' fail - %s", this.name, this.wfileName, err)
	}
	_, err = this.wfd.Write(bs)
	if err != nil {
		return false, logger.Error(tag, "'%s' write2 '%s' fail - %s", this.name, this.wfileName, err)
	}
	binary.BigEndian.PutUint64(this.wbuffer, uint64(ver))
	_, err = this.wfd.Write(this.wbuffer)
	if err != nil {
		return false, logger.Error(tag, "'%s' write3 '%s' fail - %s", this.name, this.wfileName, err)
	}
	this.lastver = ver
	this.current.fileSize = this.current.fileSize + total
	if this.config.SyncOpNum > 0 {
		this.opc++
		if this.opc >= this.config.SyncOpNum {
			this.opc = 0
			this.wfd.Sync()
		}
	}

	// push to waiting reader
	for rd, _ := range this.readers {
		// fmt.Println(rd.seeking, rd.listener != nil, rd.lastver, old)
		if !rd.seeking && rd.listener != nil && rd.lastver == old {
			rd.pushData(ver, bs, false)
		}
	}
	return true, nil
}

func (this *Service) doReadLength(fd *os.File, rbuf []byte) (uint32, error) {
	if fd == nil {
		return 0, errors.New("closed")
	}
	_, err := fd.Read(rbuf[0:4])
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(rbuf)
	return v, nil
}

func (this *Service) doReadVer(fd *os.File, rbuf []byte) (clusterbase.OpVer, error) {
	if fd == nil {
		return 0, errors.New("closed")
	}
	_, err := fd.Read(rbuf[0:8])
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(rbuf)
	return clusterbase.OpVer(v), nil
}

func (this *Service) doSeekProcess(fd *os.File, info *binlogInfo, rbuf []byte, fver clusterbase.OpVer, num int) (bool, error) {
	if fver == info.beginVer {
		return true, nil
	}
	if fver > info.lastVer {
		return false, nil
	}
	for i := 0; i < num; i++ {
		l, err := this.doReadLength(fd, rbuf)
		if err != nil {
			return false, err
		}
		_, err = fd.Seek(int64(l), os.SEEK_CUR)
		if err != nil {
			return false, err
		}
		ver, err2 := this.doReadVer(fd, rbuf)
		if err2 != nil {
			return false, err2
		}
		if ver >= fver {
			return true, nil
		}
	}
	return false, nil
}

func (this *Service) doTruncate(ver clusterbase.OpVer) (bool, error) {
	if ver > this.lastver {
		return false, nil
	}
	if ver == this.lastver {
		return true, nil
	}
	// close all
	this.doClose()

	fd, info, err := this.doSelectOpen(ver, os.O_RDWR)
	if err != nil {
		return false, err
	}
	defer fd.Close()
	if info == nil {
		return false, fmt.Errorf("can't truncate to %d", ver)
	}
	rbuf := make([]byte, 8)
	ok, err2 := this.doSeekProcess(fd, info, rbuf, ver, 0x7FFFFFFF)
	if err2 != nil {
		return false, err2
	}
	if !ok {
		return false, fmt.Errorf("can't seek %d when truncate", ver)
	}
	pos, err3 := fd.Seek(0, 1)
	if err3 != nil {
		return false, err3
	}
	logger.Debug(tag, "%s truncate '%s' %d", this, this.wfileName, pos)
	fd.Truncate(pos)
	fd.Close()

	this.lastver = ver
	info.lastVer = ver
	info.fileSize = pos
	idx := -1
	for i, bi := range this.infos {
		if bi.num > info.num {
			if idx == -1 {
				idx = i
			}
			fn := this.config.GetFileName(bi.num)
			logger.Debug(tag, "%s truncate remove '%s'", this, fn)
			err = os.Remove(fn)
			if err != nil {
				return false, err
			}
		}
	}
	if idx >= 0 {
		this.infos = this.infos[0:idx]
	}
	err = this.doOpenWrite(false)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *Service) NewWriter() (*Writer, error) {
	w := new(Writer)
	w.service = this
	return w, nil
}

func (this *Service) NewReader() (*Reader, error) {
	rd := new(Reader)
	rd.service = this
	rd.seeking = true
	return rd, nil
}
