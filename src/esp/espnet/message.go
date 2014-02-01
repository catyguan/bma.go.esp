package espnet

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"esp/espnet/protpack"
	"fmt"
	"io"
	"logger"
)

type MessageKind byte

func (O MessageKind) String() string {
	switch O {
	case MK_RESPONSE:
		return "RESP"
	case MK_REQUEST:
		return "REQS"
	case MK_EVENT:
		return "EVENT"
	case MK_INFO:
		return "INFO"
	}
	return "UNKN"
}

// MessageValuesObj
var (
	notValueErr error = errors.New("not correct value")
)

type MessageValues struct {
	m     *Message
	coder mvCoder
}

func (this *MessageValues) Set(key string, value interface{}) {
	this.coder.Set(this.m.pack, key, value, nil)
}

func (this *MessageValues) Get(key string) (interface{}, error) {
	return this.coder.Get(this.m.pack, key, nil)
}

func (this *MessageValues) GetString(key string, defv string) (string, error) {
	v, err := this.Get(key)
	if err != nil {
		return "", err
	}
	return valutil.ToString(v, defv), nil
}

func (this *MessageValues) GetInt(key string, defv int64) (int64, error) {
	v, err := this.Get(key)
	if err != nil {
		return defv, err
	}
	return valutil.ToInt64(v, defv), nil
}

func (this *MessageValues) GetUint(key string, defv uint64) (uint64, error) {
	v, err := this.Get(key)
	if err != nil {
		return defv, err
	}
	return valutil.ToUint64(v, defv), nil
}

func (this *MessageValues) GetBool(key string) (bool, error) {
	v, err := this.Get(key)
	if err != nil {
		return false, err
	}
	r, ok := valutil.ToBoolNil(v)
	if ok {
		return r, nil
	}
	return false, errors.New("not bool")
}

func (this *MessageValues) Del(key string) {
	this.coder.Remove(this.m.pack, key)
}

func (this *MessageValues) List() []string {
	return this.coder.List(this.m.pack)
}

// Message
func NewMessage() *Message {
	r := new(Message)
	r.pack = protpack.NewPackage()
	return r
}

func NewRequestMessage() *Message {
	r := NewMessage()
	FrameCoders.MessageKind.Set(r.pack, MK_REQUEST)
	return r
}

func NewReplyMessage(msg *Message) *Message {
	r := NewMessage()
	p1 := msg.pack
	p2 := r.pack

	p2.PushBack(protpack.NewFrameV(MT_MESSAGE_KIND, MK_RESPONSE, FrameCoders.MessageKind))
	for e := p1.Front(); e != nil; e = e.Next() {
		switch e.MessageType() {
		case MT_SESSION_INFO:
			p2.PushBack(e.Clone(0, false))
		case MT_HEADER, MT_DATA, MT_PAYLOAD, MT_TRACE, MT_TRACE_RESP:
			continue
		case MT_MESSAGE_KIND:
			continue
		case MT_SOURCE_ADDRESS:
			p2.PushBack(e.Clone(MT_ADDRESS, false))
			continue
		case MT_MESSAGE_ID:
			p2.PushBack(e.Clone(MT_SOURCE_MESSAGE_ID, false))
			continue
		}
	}

	return r
}

func NewPackageMessage(pack *protpack.Package) *Message {
	r := new(Message)
	r.pack = pack
	return r
}

func NewBytesMessage(bs []byte) (*Message, error) {
	pr := protpack.NewPackageReader()
	pr.Append(bs)
	p, err := pr.ReadPackage(len(bs) + 1)
	if err != nil {
		return nil, err
	}
	return NewPackageMessage(p), nil
}

type Message struct {
	pack *protpack.Package
}

func (this *Message) Id() uint64 {
	return FrameCoders.MessageId.Get(this.pack)
}

func (this *Message) SureId() uint64 {
	return FrameCoders.MessageId.Sure(this.pack)
}

func (this *Message) SetId(v uint64) {
	FrameCoders.MessageId.Set(this.pack, v)
}

func (this *Message) dumpValues(buf *bytes.Buffer, vs *MessageValues) {
	if vs == nil {
		return
	}
	ln := vs.List()
	for i, n := range ln {
		if i > 0 {
			buf.WriteString("; ")
		}
		v, err := vs.Get(n)
		if err != nil {
			v = "<ERR:" + err.Error() + ">"
		}
		buf.WriteString(fmt.Sprintf("%s=%v", n, v))
	}
}

func (this *Message) dumpXData(buf *bytes.Buffer, it *XDataIterator) {
	if it == nil {
		return
	}
	for i := 0; !it.IsEnd(); it.Next() {
		if i > 0 {
			buf.WriteString("; ")
		}
		i++
		buf.WriteString(fmt.Sprintf("%d", it.Xid()))
	}
}

func (this *Message) Dump() string {
	return this.pack.String()
}

func (this *Message) GetAddress() *Address {
	return NewAddressP(this.pack, byte(FrameCoders.Address))
}

func (this *Message) SetAddress(addr *Address) {
	if addr.pack != nil && addr.pack == this.pack {
		return
	}
	addr.Bind(this.pack, byte(FrameCoders.Address))
}

func (this *Message) GetSourceAddress() *Address {
	return NewAddressP(this.pack, byte(FrameCoders.SourceAddress))
}

func (this *Message) SetSourceAddress(addr *Address) {
	if addr.pack != nil && addr.pack == this.pack {
		return
	}
	addr.Bind(this.pack, byte(FrameCoders.SourceAddress))
}

func (this *Message) GetKind() MessageKind {
	return FrameCoders.MessageKind.Get(this.pack)
}

func (this *Message) SetKind(mt MessageKind) {
	FrameCoders.MessageKind.Set(this.pack, mt)
}

func (this *Message) Headers() *MessageValues {
	return &MessageValues{this, FrameCoders.Header}
}
func (this *Message) Datas() *MessageValues {
	return &MessageValues{this, FrameCoders.Data}
}
func (this *Message) XDatas() *MessageXData {
	return &MessageXData{this, FrameCoders.XData}
}
func (this *Message) XDataIterator() *XDataIterator {
	it := FrameCoders.XData.Iterator(this.ToPackage())
	it.moveFirst()
	return it
}
func (this *Message) GetPayload() (io.Reader, int) {
	e := this.pack.FrameByType(MT_PAYLOAD)
	if e != nil {
		v, err := e.Data()
		if err != nil {
			logger.Debug(tag, "get payload fail - %s", err)
			return nil, 0
		}
		return v.NewReader(), v.DataSize()
	}
	return nil, 0
}
func (this *Message) SetPayload(data []byte) {
	FrameCoders.Payload.Remove(this.pack)
	this.pack.PushBack(protpack.NewFrameV(MT_PAYLOAD, data, FrameCoders.Payload))
}
func (this *Message) GetPayloadB() []byte {
	r := FrameCoders.Payload.Get(this.pack)
	return r
}

func (this *Message) Clone() *Message {
	r := NewMessage()
	p1 := this.pack
	p2 := r.pack
	for e := p1.Front(); e != nil; e = e.Next() {
		p2.PushBack(e.Clone(0, false))
	}
	return r
}

func (this *Message) ToPackage() *protpack.Package {
	return this.pack
}

// helper
func (this *Message) ToError() error {
	v, err := FrameCoders.Header.Get(this.pack, "error", nil)
	if err != nil {
		return nil
	}
	if v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return errors.New(s)
}

func (this *Message) BeError(err error) {
	FrameCoders.Header.Set(this.pack, "error", err.Error(), nil)
}

func (this *Message) BeErrorS(err string) {
	FrameCoders.Header.Set(this.pack, "error", err, nil)
}

func (this *Message) ReplyMessage() *Message {
	return NewReplyMessage(this)
}

func (this *Message) TryRelyError(ch Channel, err error) {
	if this.GetKind() == MK_REQUEST {
		rmsg := this.ReplyMessage()
		rmsg.BeError(err)
		ch.SendMessage(rmsg)
	}
}
