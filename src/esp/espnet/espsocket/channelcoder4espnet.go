package espsocket

import (
	"bmautil/valutil"
	"esp/espnet/esnp"
	"logger"
)

type ChannelCoder4Espnet struct {
	maxpackage int
	reader     *esnp.MessageReader
}

func (this *ChannelCoder4Espnet) Init() {
	this.maxpackage = 0xFFFFFF
	this.reader = esnp.NewMessageReader()
}

func (this *ChannelCoder4Espnet) EncodeMessage(ch *SocketChannel, ev interface{}, next func(ev interface{}) error) error {
	if ev != nil {
		var p *esnp.Message
		if m, ok := ev.(*esnp.Message); ok {
			p = m
		}
		if p != nil {
			if logger.EnableDebug(tag) {
				logger.Debug(tag, "wpackage -> %s", p)
			}
			b, err := p.ToBytes()
			if err != nil {
				return err
			}
			return next(b)
		}
	}
	return next(ev)
}

func (this *ChannelCoder4Espnet) Decode(ch *SocketChannel, b []byte, next func(ev interface{}) error) error {
	return this.doDecode(ch, b, next)
}

func (this *ChannelCoder4Espnet) DecodeMessage(ch *SocketChannel, b []byte, next func(ev interface{}) error) error {
	return this.doDecode(ch, b, next)
}

func (this *ChannelCoder4Espnet) doDecode(who interface{}, b []byte, next func(ev interface{}) error) error {
	pr := this.reader
	pr.Append(b)
	for {
		mp := this.maxpackage
		p, err := pr.ReadMessage(mp)
		if err != nil {
			logger.Error(tag, "rpackage -> fail, %s, %s", who, err)
			return err
		}
		if p == nil {
			return nil
		}
		if logger.EnableDebug(tag) {
			if p != nil {
				logger.Debug(tag, "rpackage -> %s", p)
			}
		}
		next(p)
	}
}

func (this *ChannelCoder4Espnet) SetProperty(name string, val interface{}) bool {
	switch name {
	case PROP_ESPNET_MAXPACKAGE:
		this.maxpackage = valutil.ToInt(val, 0)
		return true
	}
	return false
}

func (this *ChannelCoder4Espnet) GetProperty(name string) (interface{}, bool) {
	switch name {
	case PROP_ESPNET_MAXPACKAGE:
		return this.maxpackage, true
	}
	return nil, false
}

func init() {
	RegSocketChannelCoder(SOCKET_CHANNEL_CODER_ESPNET, func() SocketChannelCoder {
		r := new(ChannelCoder4Espnet)
		r.Init()
		return r
	})
}
