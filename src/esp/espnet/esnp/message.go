package esnp

import (
	"bytes"
	"fmt"
)

func NewBytesMessage(bs []byte) (*Message, error) {
	pr := NewMessageReaderB(bs)
	p := NewMessage()
	p, err := pr.ReadMessage(0)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("unknow message format")
	}
	return p, nil
}

func (this *Message) Id() uint64 {
	return MessageLineCoders.MessageId.Get(this)
}

func (this *Message) SureId() uint64 {
	return MessageLineCoders.MessageId.Sure(this)
}

func (this *Message) SetId(v uint64) {
	MessageLineCoders.MessageId.Set(this, v)
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
	return this.String()
}

func (this *Message) GetAddress() *Address {
	return NewAddressP(this, byte(MessageLineCoders.Address))
}

func (this *Message) SetAddress(addr *Address) {
	if addr.message != nil && addr.message == this {
		return
	}
	addr.Bind(this, byte(MessageLineCoders.Address))
}

func (this *Message) GetSourceAddress() *Address {
	return NewAddressP(this, byte(MessageLineCoders.SourceAddress))
}

func (this *Message) SetSourceAddress(addr *Address) {
	if addr.message != nil && addr.message == this {
		return
	}
	addr.Bind(this, byte(MessageLineCoders.SourceAddress))
}

func (this *Message) GetVersion() *Version {
	return MessageLineCoders.Version.Get(this)
}

func (this *Message) SetVersion(val *Version) {
	MessageLineCoders.Version.Set(this, val)
}

func (this *Message) IsRequest() bool {
	if MessageLineCoders.Flag.Has(this, FLAG_REQUEST) {
		return !MessageLineCoders.Flag.Has(this, FLAG_RESP)
	}
	return false
}

func (this *Message) SureRequest() {
	MessageLineCoders.Flag.Set(this, FLAG_REQUEST)
}

func (this *Message) Headers() *MessageValues {
	return &MessageValues{this, MessageLineCoders.Header}
}
func (this *Message) Datas() *MessageValues {
	return &MessageValues{this, MessageLineCoders.Data}
}
func (this *Message) XDatas() *MessageXData {
	return &MessageXData{this, MessageLineCoders.XData}
}
func (this *Message) XDataIterator() *XDataIterator {
	it := MessageLineCoders.XData.Iterator(this)
	it.moveFirst()
	return it
}
func (this *Message) GetPayload() []byte {
	r := MessageLineCoders.Payload.Get(this)
	return r
}
func (this *Message) SetPayload(data []byte) {
	MessageLineCoders.Payload.Remove(this)
	this.PushBack(NewMessageLineV(MLT_PAYLOAD, data, MessageLineCoders.Payload))
}

func (this *Message) Clone() *Message {
	r := NewMessage()
	p1 := this
	p2 := r
	for e := p1.Front(); e != nil; e = e.Next() {
		p2.PushBack(e.Clone(0))
	}
	return r
}

// helper
type remoteErr struct {
	err string
}

func (this *remoteErr) Error() string {
	return this.err
}

func (this *remoteErr) String() string {
	return this.Error()
}

func NewRemoteErr(s string) error {
	return &remoteErr{s}
}

func IsRemoteErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*remoteErr)
	return ok
}

func (this *Message) ToError() error {
	ok, v := MessageLineCoders.Error.Get(this)
	if !ok {
		return nil
	}
	return NewRemoteErr(v)
}

func (this *Message) BeError(err error) {
	MessageLineCoders.Error.Set(this, err.Error())
}

func (this *Message) BeErrorS(err string) {
	MessageLineCoders.Error.Set(this, err)
}

func (this *Message) ReplyMessage() *Message {
	return NewReplyMessage(this)
}

// Message
func NewRequestMessage() *Message {
	r := NewMessage()
	MessageLineCoders.Flag.Set(r, FLAG_REQUEST)
	return r
}

func NewRequestMessageWithId() *Message {
	r := NewRequestMessage()
	r.SureId()
	return r
}

func NewReplyMessage(msg *Message) *Message {
	r := NewMessage()
	p1 := msg
	p2 := r

	p2.PushFront(NewMessageLineV(MLT_FLAG, FLAG_RESP, MessageLineCoders.Flag))
	for e := p1.Front(); e != nil; e = e.Next() {
		switch e.MessageType() {
		case MLT_SESSION_INFO:
			p2.PushBack(e.Clone(0))
		case MLT_HEADER, MLT_DATA, MLT_PAYLOAD, MLT_TRACE, MLT_TRACE_RESP:
			continue
		case MLT_FLAG:
			o, err := e.Value(MessageLineCoders.Flag)
			if err == nil {
				if fo, ok := o.(Flag); ok {
					switch fo {
					case FLAG_REQUEST, FLAG_INFO:
						p2.PushBack(e.Clone(0))
					}
				}
			}
			continue
		case MLT_SOURCE_ADDRESS:
			p2.PushBack(e.Clone(MLT_ADDRESS))
			continue
		case MLT_MESSAGE_ID:
			p2.PushBack(e.Clone(MLT_SOURCE_MESSAGE_ID))
			continue
		}
	}

	return r
}
