package xmemclient

import (
	"esp/espnet"
	"esp/xmem/xmemprot"
	"fmt"
	"logger"
	"time"
)

const (
	tag = "xmemclient"
)

type Xmem4Client struct {
	c       *espnet.ChannelClient
	a       espnet.Address
	g       string
	Timeout time.Duration
}

func NewClient(c *espnet.ChannelClient, addr espnet.Address, group string) *Xmem4Client {
	this := new(Xmem4Client)
	this.c = c
	this.a = addr
	this.g = group
	this.Timeout = time.Duration(5) * time.Second
	return this
}

func (this *Xmem4Client) ToXMem() xmemprot.XMem {
	return this
}

func (this *Xmem4Client) String() string {
	return fmt.Sprintf("XMem4Client[%v,%v,%s]", this.c, this.a, this.g)
}

func (this *Xmem4Client) call(msg *espnet.Message) (*espnet.Message, error) {
	var tm *time.Timer
	if this.Timeout > 0 {
		tm = time.NewTimer(this.Timeout)
	}
	return this.c.Call(msg, tm)
}

func (this *Xmem4Client) Get(key xmemprot.MemKey) (interface{}, xmemprot.MemVer, bool, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestGet)
	req.Init(this.g, key)
	req.Write(msg)

	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s Get(%v) fail - %s", this, key, err2)
		return nil, xmemprot.VERSION_INVALID, false, err2
	}

	o := new(xmemprot.SHResponseGet)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s Get(%v) decode fail - %s", this, key, err3)
		return nil, xmemprot.VERSION_INVALID, false, err3
	}
	logger.Debug(tag, "%s Get(%v) -> %v", this, key, o)
	return o.Value, o.Version, !o.Miss, nil
}

func (this *Xmem4Client) GetAndListen(key xmemprot.MemKey, id string, lis xmemprot.XMemListener) (interface{}, xmemprot.MemVer, bool, error) {
	return nil, xmemprot.VERSION_INVALID, false, fmt.Errorf("not impl")
}

func (this *Xmem4Client) List(key xmemprot.MemKey) ([]string, bool, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestList)
	req.Init(this.g, key)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s List(%v) fail - %s", this, key, err2)
		return nil, false, err2
	}

	o := new(xmemprot.SHResponseList)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s List(%v) decode fail - %s", this, key, err3)
		return nil, false, err3
	}
	logger.Debug(tag, "%s List(%v) -> %v", this, key, o)
	return o.Names, !o.Miss, nil
}

func (this *Xmem4Client) ListAndListen(key xmemprot.MemKey, id string, lis xmemprot.XMemListener) ([]string, bool, error) {
	return nil, false, fmt.Errorf("not impl")
}

func (this *Xmem4Client) AddListener(key xmemprot.MemKey, id string, lis xmemprot.XMemListener) error {
	return fmt.Errorf("not impl")
}

func (this *Xmem4Client) RemoveListener(key xmemprot.MemKey, id string) error {
	return fmt.Errorf("not impl")
}

func (this *Xmem4Client) Set(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestSet)
	req.InitSet(this.g, key, val, sz)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s Set(%v) fail - %s", this, key, err2)
		return xmemprot.VERSION_INVALID, err2
	}

	o := new(xmemprot.SHResponseSet)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s Set(%v) decode fail - %s", this, key, err3)
		return xmemprot.VERSION_INVALID, err3
	}
	logger.Debug(tag, "%s Set(%v) -> %v", this, key, o)
	return o.Version, nil
}

func (this *Xmem4Client) CompareAndSet(key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer) (xmemprot.MemVer, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestSet)
	req.InitCompareAndSet(this.g, key, val, sz, ver)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s CompareAndSet(%v) fail - %s", this, key, err2)
		return xmemprot.VERSION_INVALID, err2
	}

	o := new(xmemprot.SHResponseSet)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s CompareAndSet(%v) decode fail - %s", this, key, err3)
		return xmemprot.VERSION_INVALID, err3
	}
	logger.Debug(tag, "%s CompareAndSet(%v) -> %v", this, key, o)
	return o.Version, nil
}

func (this *Xmem4Client) SetIfAbsent(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestSet)
	req.InitSetIfAbsent(this.g, key, val, sz)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s SetIfAbsent(%v) fail - %s", this, key, err2)
		return xmemprot.VERSION_INVALID, err2
	}

	o := new(xmemprot.SHResponseSet)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s SetIfAbsent(%v) decode fail - %s", this, key, err3)
		return xmemprot.VERSION_INVALID, err3
	}
	logger.Debug(tag, "%s SetIfAbsent(%v) -> %v", this, key, o)
	return o.Version, nil
}

func (this *Xmem4Client) Delete(key xmemprot.MemKey) (bool, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestDelete)
	req.Init(this.g, key)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s Delete(%v) fail - %s", this, key, err2)
		return false, err2
	}

	o := new(xmemprot.SHResponseDelete)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s Delete(%v) decode fail - %s", this, key, err3)
		return false, err3
	}
	logger.Debug(tag, "%s Delete(%v) -> %v", this, key, o)
	return o.Done, nil
}

func (this *Xmem4Client) CompareAndDelete(key xmemprot.MemKey, ver xmemprot.MemVer) (bool, error) {
	msg := espnet.NewMessage()
	msg.SetAddress(this.a)
	req := new(xmemprot.SHRequestDelete)
	req.InitCompareAndDelete(this.g, key, ver)
	req.Write(msg)
	rmsg, err2 := this.call(msg)
	if err2 != nil {
		logger.Debug(tag, "%s CompareAndDelete(%v) fail - %s", this, key, err2)
		return false, err2
	}

	o := new(xmemprot.SHResponseDelete)
	err3 := o.Read(rmsg)
	if err3 != nil {
		logger.Debug(tag, "%s CompareAndDelete(%v) decode fail - %s", this, key, err3)
		return false, err3
	}
	logger.Debug(tag, "%s CompareAndDelete(%v) -> %v", this, key, o)
	return o.Done, nil
}
