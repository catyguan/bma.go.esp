package xmemservice

import (
	"bmautil/binlog"
	"bmautil/byteutil"
	xcoder "bmautil/coder"
	"esp/xmem/xmemprot"
	"fmt"
	"io/ioutil"
	"logger"
)

type XMemOP int

const (
	OP_NONE = iota
	OP_SET
	OP_DELETE
)

type XMemBinlog struct {
	Op       XMemOP
	Key      string
	Value    interface{}
	Size     int
	Version  xmemprot.MemVer
	IsAbsent bool
	// SET
	// group string, key MemKey, val interface{}, sz int, ver MemVer, isAbsent bool
	// DELETE
	// group string, key MemKey, ver MemVer
}

func (this *XMemBinlog) Encode(coder XMemCoder) ([]byte, error) {
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	w.WriteByte(byte(this.Op))
	xcoder.LenString.DoEncode(w, this.Key)
	if this.Op == OP_SET {
		flag, bs, err := coder.Encode(this.Value)
		if err != nil {
			return nil, err
		}
		xcoder.LenString.DoEncode(w, flag)
		xcoder.Int.DoEncode(w, len(bs))
		if bs != nil {
			w.Write(bs)
		}
	}
	if this.Version == xmemprot.VERSION_INVALID {
		w.WriteByte(0)
	} else {
		w.WriteByte(1)
		xcoder.Uint64.DoEncode(w, uint64(this.Version))
	}
	if this.Op == OP_SET {
		b := 0
		if this.IsAbsent {
			b = 1
		}
		w.WriteByte(byte(b))
	}
	return w.End().ToBytes(), nil
}

func DecodeBinlog(data []byte, coder XMemCoder, maxlen int) (*XMemBinlog, error) {
	this := new(XMemBinlog)
	buf := byteutil.NewBytesBufferB(data)
	r := buf.NewReader()
	v1, err1 := r.ReadByte()
	if err1 != nil {
		return nil, err1
	}
	this.Op = XMemOP(v1)
	v2, err2 := xcoder.LenString.DoDecode(r, 1024*100)
	if err2 != nil {
		return nil, err2
	}
	this.Key = v2
	if this.Op == OP_SET {
		v3, err3 := xcoder.LenString.DoDecode(r, maxlen)
		if err3 != nil {
			return nil, err3
		}
		v4, err4 := xcoder.Int.DoDecode(r)
		if err4 != nil {
			return nil, err4
		}
		bs := make([]byte, v4)
		_, err5 := r.Read(bs)
		if err5 != nil {
			return nil, err5
		}
		val, sz, err6 := coder.Decode(v3, bs)
		if err6 != nil {
			return nil, err6
		}
		this.Value = val
		this.Size = sz
	}
	v7, err7 := r.ReadByte()
	if err7 != nil {
		return nil, err7
	}
	if v7 == 0 {
		this.Version = xmemprot.VERSION_INVALID
	} else {
		v8, err8 := xcoder.Uint64.DoDecode(r)
		if err8 != nil {
			return nil, err8
		}
		this.Version = xmemprot.MemVer(v8)
	}
	if this.Op == OP_SET {
		v9, err9 := r.ReadByte()
		if err9 != nil {
			return nil, err9
		}
		if v9 == 1 {
			this.IsAbsent = true
		}
	}
	return this, nil
}

// Service
func (this *Service) doGetBinogVersion(name string) (master binlog.BinlogVer, slave binlog.BinlogVer, rerr error) {
	si, err := this.doGetGroup(name)
	if err != nil {
		return 0, 0, err
	}
	sv := si.group.blver
	if si.group.blwriter == nil {
		logger.Debug(tag, "'%s' binlog not start", name)
		return -1, sv, nil
	}
	mv, err1 := si.group.blwriter.GerVersion()
	if err1 != nil {
		logger.Debug(tag, "'%s' binlog get version fail - %s", name, err1)
		return -1, sv, nil
	}
	return mv, sv, nil
}

func (this *Service) doStartBinlog(name string, mg *localMemGroup, cfg *MemGroupConfig) error {
	if mg.blservice != nil {
		logger.Debug(tag, "'%s' already start binlog, skip", name)
		return nil
	}
	if !cfg.IsEnableBinlog() {
		return logger.Warn(tag, "'%s' binlog not enable", name)
	}
	s := binlog.NewBinLog(name, 16, cfg.BLConfig)
	if !s.Run() {
		return fmt.Errorf("'%s' binlog start fail", name)
	}
	logger.Debug(tag, "'%s' start binlog", name)
	mg.blservice = s
	if !cfg.BLConfig.Readonly {
		mg.blwriter, _ = s.NewWriter()
	}
	return nil
}

func (this *Service) doStopBinlog(name string, mg *localMemGroup) error {
	if mg.blservice != nil {
		logger.Debug(tag, "'%s' stop binlog", name)
		mg.blservice.Stop()
		mg.blservice = nil
		mg.blwriter = nil
	}
	return nil
}

func (this *Service) doWriteBinlog(group string, si *serviceItem, bl *XMemBinlog) {
	if si.group.blwriter == nil {
		logger.Warn(tag, "'%s' binlog not start, lost op=%d", group, bl.Op)
		return
	}
	if si.profile == nil {
		logger.Warn(tag, "'%s' profile invalid, lost op=%d", group, bl.Op)
		return
	}
	bs, err := bl.Encode(si.profile.Coder)
	if err != nil {
		logger.Warn(tag, "'%s' binlog encode fail - %s", group, err)
		return
	}
	logger.Debug(tag, "'%s' binlog op=%d", group, bl.Op)
	if !si.group.blwriter.Write(bs) {
		logger.Warn(tag, "'%s' binlog write fail", group)
	}
}

func (this *Service) doSaveBinlogSnapshot(name string, fileName string) error {
	si, err := this.doGetGroup(name)
	if err != nil {
		return err
	}
	if si.profile == nil {
		return fmt.Errorf("'%s' no profile", name)
	}

	mver, _, err1 := this.doGetBinogVersion(name)
	if err1 != nil {
		return err1
	}
	if mver <= 0 {
		return fmt.Errorf("'%s' no master binlog version", name)
	}

	logger.Debug(tag, "doBinlogSave(%s,%s)", name, fileName)
	gss, err2 := si.group.Snapshot(si.profile.Coder)
	if err2 != nil {
		return err2
	}
	gss.BLVer = mver
	bs, err3 := gss.Encode()
	if err3 != nil {
		return err3
	}
	return ioutil.WriteFile(fileName, bs, 0664)
}

func (this *Service) doRunBinlog(name string, fileName string) error {
	si, err := this.doGetGroup(name)
	if err != nil {
		return err
	}
	if si.profile == nil {
		return fmt.Errorf("'%s' no profile", name)
	}

	cfg := new(binlog.BinlogConfig)
	cfg.FileName = fileName
	cfg.Readonly = true

	bls := binlog.NewBinLog(name+"_reader", 16, cfg)
	if !bls.Run() {
		return logger.Warn("'%s' binlog reader start fail - %s", name)
	}
	defer bls.Stop()

	rd, err := bls.NewReader()
	if err != nil {
		logger.Warn(tag, "'%s' binlog reader create fail - %s", name, err)
		return err
	}

	logger.Info(tag, "'%s' process binlog file '%s' start", name, fileName)
	defer func() {
		logger.Info(tag, "'%s' process binlog file '%s' end", name, fileName)
	}()

	_, err = rd.Seek(si.group.blver)
	if err != nil {
		logger.Warn(tag, "'%s' binlog reader seek fail - %s", name, err)
		return err
	}

	c := 0
	for {
		blver, bs, err1 := rd.Read()
		if err1 != nil {
			logger.Warn(tag, "'%s' process %d binlog and read fail -%d", name, c, err1)
			return err1
		}
		if bs == nil {
			break
		}
		err2 := this.doProcessBinog(name, blver, bs)
		if err2 != nil {
			logger.Warn(tag, "'%s' process %d binlog and fail -%d", name, c, err2)
			return err2
		}
		c++
	}
	logger.Info(tag, "'%s' process %d binlog", name, c)

	return nil
}

func (this *Service) doProcessBinog(group string, blver binlog.BinlogVer, bs []byte) error {
	si, err := this.doGetGroup(group)
	if err != nil {
		return err
	}
	if si.profile == nil {
		return fmt.Errorf("'%s' profile nil", group)
	}
	bl, err1 := DecodeBinlog(bs, si.profile.Coder, 1024*1024)
	if err1 != nil {
		return err1
	}
	key := xmemprot.MemKeyFromString(bl.Key)
	switch bl.Op {
	case OP_SET:
		_, err2 := this.doSetOp(group, key, bl.Value, bl.Size, bl.Version, bl.IsAbsent)
		if err2 != nil {
			return err2
		}
	case OP_DELETE:
		_, err2 := this.doDeleteOp(group, key, bl.Version)
		if err2 != nil {
			return err2
		}
	default:
		return fmt.Errorf("'%s' unknow op %d", group, bl.Op)
	}
	si.group.blver = blver
	return nil
}
