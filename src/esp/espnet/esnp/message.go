package esnp

import (
	"bmautil/valutil"
	"bytes"
	"errors"
	"fmt"
)

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
	r.pack = NewPackage()
	return r
}

func NewRequestMessage() *Message {
	r := NewMessage()
	FrameCoders.Flag.Set(r.pack, FLAG_REQUEST)
	return r
}

func NewReplyMessage(msg *Message) *Message {
	r := NewMessage()
	p1 := msg.pack
	p2 := r.pack

	p2.PushFront(NewFrameV(MT_FLAG, FLAG_RESP, FrameCoders.Flag))
	for e := p1.Front(); e != nil; e = e.Next() {
		switch e.MessageType() {
		case MT_SESSION_INFO:
			p2.PushBack(e.Clone(0))
		case MT_HEADER, MT_DATA, MT_PAYLOAD, MT_TRACE, MT_TRACE_RESP:
			continue
		case MT_FLAG:
			o, err := e.Value(FrameCoders.Flag)
			if err == nil {
				if fo, ok := o.(MTFlag); ok {
					switch fo {
					case FLAG_REQUEST, FLAG_INFO:
						p2.PushBack(e.Clone(0))
					}
				}
			}
			continue
		case MT_SOURCE_ADDRESS:
			p2.PushBack(e.Clone(MT_ADDRESS))
			continue
		case MT_MESSAGE_ID:
			p2.PushBack(e.Clone(MT_SOURCE_MESSAGE_ID))
			continue
		}
	}

	return r
}

func NewPackageMessage(pack *Package) *Message {
	r := new(Message)
	r.pack = pack
	return r
}

func NewBytesMessage(bs []byte) (*Message, error) {
	pr := NewPackageReader()
	pr.Append(bs)
	p, err := pr.ReadPackage(len(bs) + 1)
	if err != nil {
		return nil, err
	}
	return NewPackageMessage(p), nil
}

type Message struct {
	pack *Package
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

func (this *Message) GetVersion() *MTVersion {
	return FrameCoders.Version.Get(this.pack)
}

func (this *Message) SetVersion(val *MTVersion) {
	FrameCoders.Version.Set(this.pack, val)
}

func (this *Message) IsRequest() bool {
	if FrameCoders.Flag.Has(this.pack, FLAG_REQUEST) {
		return !FrameCoders.Flag.Has(this.pack, FLAG_RESP)
	}
	return false
}

func (this *Message) SureRequest() {
	FrameCoders.Flag.Set(this.pack, FLAG_REQUEST)
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
func (this *Message) GetPayload() []byte {
	r := FrameCoders.Payload.Get(this.pack)
	return r
}
func (this *Message) SetPayload(data []byte) {
	FrameCoders.Payload.Remove(this.pack)
	this.pack.PushBack(NewFrameV(MT_PAYLOAD, data, FrameCoders.Payload))
}

func (this *Message) Clone() *Message {
	r := NewMessage()
	p1 := this.pack
	p2 := r.pack
	for e := p1.Front(); e != nil; e = e.Next() {
		p2.PushBack(e.Clone(0))
	}
	return r
}

func (this *Message) ToPackage() *Package {
	return this.pack
}

// helper
func (this *Message) ToError() error {
	ok, v := FrameCoders.Error.Get(this.pack)
	if !ok {
		return nil
	}
	return errors.New(v)
}

func (this *Message) BeError(err error) {
	FrameCoders.Error.Set(this.pack, err.Error())
}

func (this *Message) BeErrorS(err string) {
	FrameCoders.Error.Set(this.pack, err)
}

func (this *Message) ReplyMessage() *Message {
	return NewReplyMessage(this)
}
