package tbus

import (
	"bmautil/byteutil"
	"bmautil/valutil"
	"encoding/binary"
	"esp/espnet"
	"logger"
)

type ChannelCoder4TBus struct {
	maxframe int

	frameBody int
	seqno     int
	buffer    *byteutil.BytesBuffer
}

func NewChannelCoder(maxframe int) *ChannelCoder4TBus {
	this := new(ChannelCoder4TBus)
	this.maxframe = maxframe
	this.buffer = byteutil.NewBytesBuffer()
	this.frameBody = -1
	return this
}

func (this *ChannelCoder4TBus) EncodeMessage(ch *espnet.SocketChannel, ev interface{}, next func(ev interface{}) error) error {
	if ev != nil {
		var b []byte
		if m, ok := ev.(*espnet.Message); ok {
			err := m.ToError()
			if err != nil {
				ch.AskClose()
				return nil
			}
			b = m.GetPayloadB()
			return next(b)
		}
	}
	return next(ev)
}

func (this *ChannelCoder4TBus) DecodeMessage(ch *espnet.SocketChannel, b []byte, next func(ev interface{}) error) error {
	this.buffer.Add(b)
	reader := this.buffer.NewReader()

	name := ""
	for {
		if this.frameBody < 0 {
			buf := []byte{0, 0, 0, 0}
			_, err := reader.Read(buf)
			if err != nil {
				return nil
			}
			frameSize := binary.BigEndian.Uint32(buf)
			if frameSize > uint32(this.maxframe) {
				return logger.Error(tag, "%s maxframe reach %d/%d", ch, frameSize, this.maxframe)
			}
			var message TMessage
			ok, err := message.Read(reader)
			if err != nil {
				logger.Error(tag, "decode TMessage fail - %s", err)
				return err
			}
			if !ok {
				return nil
			}
			name = message.name
			this.frameBody = int(frameSize) + 4
			this.seqno = espnet.FrameCoders.SeqNO.FirstSeqno()
		}

		// read frameBody and send it
		var buf []byte
		end := false
		if this.buffer.DataSize() >= this.frameBody {
			buf = make([]byte, this.frameBody)
			this.buffer.CheckAndPop(buf, this.frameBody)
			this.frameBody = 0
			end = true
		} else {
			buf = this.buffer.ToBytes()
			this.buffer.DataList = nil
			this.frameBody -= len(buf)
		}

		msg := espnet.NewRequestMessage()
		msg.SetAddress(espnet.NewAddress(name))
		maxseq := 0
		if end {
			maxseq = this.seqno
		}
		espnet.FrameCoders.SeqNO.Set(msg.ToPackage(), this.seqno, maxseq)

		msg.SetPayload(buf)
		this.seqno++

		err := next(msg)
		if err != nil {
			return err
		}

		if !end {
			return nil
		}
		this.frameBody = -1
	}
}

func (this *ChannelCoder4TBus) SetProperty(name string, val interface{}) bool {
	switch name {
	case espnet.PROP_ESPNET_MAXFRAME:
		this.maxframe = valutil.ToInt(val, 0)
		return true
	}
	return false
}

func (this *ChannelCoder4TBus) GetProperty(name string) (interface{}, bool) {
	switch name {
	case espnet.PROP_ESPNET_MAXFRAME:
		return this.maxframe, true
	}
	return nil, false
}
