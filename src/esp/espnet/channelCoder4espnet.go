package espnet

import (
	"bmautil/valutil"
	"esp/espnet/protpack"
	"logger"
)

type ChannelCoder4Espnet struct {
	maxframe int
	reader   *protpack.PackageReader
}

func (this *ChannelCoder4Espnet) Init() {
	this.maxframe = 10 * 1024 * 1024
	this.reader = protpack.NewPackageReader()
}

func (this *ChannelCoder4Espnet) EncodeMessage(ch *SocketChannel, ev interface{}, next func(ev interface{}) error) error {
	if ev != nil {
		var p *protpack.Package
		if m, ok := ev.(*Message); ok {
			p = m.ToPackage()
		} else {
			p, _ = ev.(*protpack.Package)
		}
		if p != nil {
			b, err := p.ToBytesBuffer()
			if err != nil {
				return err
			}
			return next(b)
		}
	}
	return next(ev)
}

func (this *ChannelCoder4Espnet) Decode(ch *SocketChannel, b []byte, next func(ev interface{}) error) error {
	return this.doDecode(ch, b, false, next)
}

func (this *ChannelCoder4Espnet) DecodeMessage(ch *SocketChannel, b []byte, next func(ev interface{}) error) error {
	return this.doDecode(ch, b, true, next)
}

func (this *ChannelCoder4Espnet) doDecode(who interface{}, b []byte, msg bool, next func(ev interface{}) error) error {
	pr := this.reader
	pr.Append(b)
	for {
		mp := this.maxframe
		p, err := pr.ReadPackage(mp)
		if err != nil {
			logger.Error(tag, "%s read package fail - %s", who, err)
			return err
		}
		if p == nil {
			// if logger.EnableDebug(tag) {
			// 	logger.Debug(tag, "reading package ## %s", pr)
			// }
			return nil
		}
		if logger.EnableDebug(tag) {
			logger.Debug(tag, "read package -> %s ## %s", p, pr)
		}
		if msg {
			next(NewPackageMessage(p))
		} else {
			next(p)
		}
	}
}

func (this *ChannelCoder4Espnet) SetProperty(name string, val interface{}) bool {
	switch name {
	case PROP_ESPNET_MAXFRAME:
		this.maxframe = valutil.ToInt(val, 0)
		return true
	}
	return false
}

func (this *ChannelCoder4Espnet) GetProperty(name string) (interface{}, bool) {
	switch name {
	case PROP_ESPNET_MAXFRAME:
		return this.maxframe, true
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
