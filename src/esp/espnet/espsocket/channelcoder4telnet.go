package espsocket

import (
	"bmautil/valutil"
	"bytes"
	"esp/espnet/esnp"
	"fmt"
	"logger"
)

type ChannelCoder4Telnet struct {
	maxpackage    int
	buffer        []byte
	Error2String  func(err error) string
	Bytes2Message func(str []byte) *esnp.Message
}

func (this *ChannelCoder4Telnet) Init() {
	this.maxpackage = 1024 * 8
}

func (this *ChannelCoder4Telnet) EncodeMessage(ch *SocketChannel, ev interface{}, next func(ev interface{}) error) error {
	if ev != nil {
		var b []byte
		if m, ok := ev.(*esnp.Message); ok {
			err := m.ToError()
			if err != nil {
				if this.Error2String != nil {
					str := this.Error2String(err)
					b = []byte(str)
				} else {
					b = []byte(fmt.Sprintf("ERROR:%s\n", err))
				}
			} else {
				b = m.GetPayload()
			}
			return next(b)
		}
	}
	return next(ev)
}

func (this *ChannelCoder4Telnet) DecodeMessage(ch *SocketChannel, b []byte, next func(ev interface{}) error) error {
	var r *bytes.Buffer
	var l int
	if this.buffer == nil {
		l = 0
		r = bytes.NewBuffer(make([]byte, 0))
	} else {
		l = len(this.buffer)
		r = bytes.NewBuffer(this.buffer)
	}
	if l+len(b) > this.maxpackage {
		return logger.Error(tag, "%s maxpackage reach %d/%d", ch, l+len(b), this.maxpackage)
	}
	r.Write(b)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			this.buffer = line
			return nil
		}
		var msg *esnp.Message
		if this.Bytes2Message != nil {
			msg = this.Bytes2Message(line)
		} else {
			msg = esnp.NewRequestMessage()
			msg.SetPayload(line)
		}
		next(msg)
	}
}

func (this *ChannelCoder4Telnet) SetProperty(name string, val interface{}) bool {
	switch name {
	case PROP_ESPNET_MAXPACKAGE:
		this.maxpackage = valutil.ToInt(val, 0)
		return true
	}
	return false
}

func (this *ChannelCoder4Telnet) GetProperty(name string) (interface{}, bool) {
	switch name {
	case PROP_ESPNET_MAXPACKAGE:
		return this.maxpackage, true
	}
	return nil, false
}
