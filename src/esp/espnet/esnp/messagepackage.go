package esnp

import (
	"bytes"
	"fmt"
)

// Message
type Message struct {
	head *MessageLine
	tail *MessageLine
	size int
	hbuf []byte
}

func NewMessage() *Message {
	this := new(Message)
	this.Init()
	return this
}

func (this *Message) Init() *Message {
	this.head = nil
	this.tail = nil
	this.size = 0
	return this
}

func (this *Message) PackageSize() int {
	return this.size
}

func (this *Message) Front() *MessageLine {
	if this.size == 0 {
		return nil
	}
	return this.head
}

func (this *Message) Back() *MessageLine {
	if this.size == 0 {
		return nil
	}
	return this.tail
}

func (this *Message) insert(e, at *MessageLine) *MessageLine {
	var n *MessageLine
	if at != nil {
		n = at.next
		at.next = e
	} else {
		n = nil
		this.head = e
	}
	e.prev = at
	e.next = n
	if n != nil {
		n.prev = e
	} else {
		this.tail = e
	}
	this.size++
	e.message = this
	return e
}

func (this *Message) remove(e *MessageLine) *MessageLine {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		this.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		this.tail = e.prev
	}
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.message = nil
	this.size--
	return e
}

func (this *Message) Remove(e *MessageLine) {
	if e.message == this {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		this.remove(e)
	}
}

func (this *Message) PushFront(f *MessageLine) *MessageLine {
	return this.insert(f, this.head)
}

func (this *Message) PushBack(f *MessageLine) *MessageLine {
	return this.insert(f, this.tail)
}

func (this *Message) InsertBefore(f, mark *MessageLine) *MessageLine {
	if mark.message != this {
		return nil
	}
	return this.insert(f, mark.prev)
}

func (this *Message) InsertAfter(f, mark *MessageLine) *MessageLine {
	if mark.message != this {
		return nil
	}
	return this.insert(f, mark)
}

func (this *Message) PushBackList(other *Message) {
	for e := other.Front(); e != nil; e.Next() {
		this.insert(e.Clone(0), this.tail)
	}
}

func (this *Message) PushFrontList(other *Message) {
	for e := other.Back(); e != nil; e = e.Prev() {
		this.insert(e.Clone(0), nil)
	}
}

func (this *Message) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("Message")
	buf.WriteString(fmt.Sprintf("[%d]:", this.size))
	for e := this.Front(); e != nil; e = e.Next() {
		if e != this.Front() {
			buf.WriteString(",")
		}
		buf.WriteString(e.String())
	}
	return buf.String()
}

func (this *Message) Write(w EncodeWriter) error {
	for e := this.Front(); e != nil; e = e.Next() {
		err := e.Encode(w)
		if err != nil {
			return err
		}
	}
	return MessageLineHeaderWrite(w, MLT_END, 0)
}

func (this *Message) ReadLineHeader(r DecodeReader) (mt byte, sz int, err error) {
	if this.hbuf == nil {
		this.hbuf = make([]byte, 4)
	}
	pos := 0
	for {
		n, err := r.Read(this.hbuf[pos:])
		if err != nil {
			return 0, 0, err
		}
		if n+pos < 4 {
			pos += n
		}
	}
	this.size += 4
	mt, sz = MessageLineHeaderRead(this.hbuf, 0)
	return
}

func (this *Message) ReadLine(r DecodeReader, mt byte, sz int) (*MessageLine, error) {
	b := make([]byte, sz, 0)
	pos := 0
	for {
		n, err := r.Read(b[pos:])
		if err != nil {
			return nil, err
		}
		if n+pos < sz {
			pos += n
		}
	}
	this.size += sz
	if mt == 0 {
		return nil, nil
	}
	return NewMessageLine(mt, b), nil
}

func (this *Message) ReadAll(r DecodeReader) error {
	for {
		mt, sz, err0 := this.ReadLineHeader(r)
		if err0 != nil {
			return err0
		}
		ml, err1 := this.ReadLine(r, mt, sz)
		if err1 != nil {
			return err1
		}
		if ml == nil {
			break
		}
		this.PushBack(ml)
	}
	return nil
}

func (this *Message) ToBytes() ([]byte, error) {
	var w BytesEncodeWriter
	err := this.Write(&w)
	if err != nil {
		return nil, err
	}
	return w.ToBytes(), nil
}

func (this *Message) MessageLineByType(mt byte) *MessageLine {
	for e := this.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			return e
		}
	}
	return nil
}

func (this *Message) RemoveMessageLine(f func(e *MessageLine) (bool, bool)) {
	e := this.Front()
	for {
		if e == nil {
			break
		}
		ne := e.Next()
		del, stop := f(e)
		if del {
			// fmt.Println("DELETE", e.MessageType())
			this.Remove(e)
		}
		if stop {
			return
		}
		e = ne
	}
}

// MessageReader
type MessageReader struct {
	buffer []byte
	wpos   int
	rpos   int
}

func NewMessageReader() *MessageReader {
	this := new(MessageReader)
	this.buffer = make([]byte, 1024)
	return this
}

func (this *MessageReader) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("[")
	buf.WriteString(fmt.Sprintf("R:%d/W:%d/C:%d", this.rpos, this.wpos, cap(this.buffer)))
	buf.WriteString("]")
	return buf.String()
}

func (this *MessageReader) Append(b []byte) {
	l := len(b)
	if l+this.wpos > cap(this.buffer) {
		gl := ((l+this.wpos)/1024 + 1) * 1024
		buf := make([]byte, gl)
		copy(buf, this.buffer[:this.wpos])
		this.buffer = buf
	}
	copy(this.buffer[this.wpos:], b)
	this.wpos = this.wpos + l
}

func (this *MessageReader) ReadMessage(mp int) (*Message, error) {
	for {
		if this.rpos+size_FHEADER > this.wpos {
			// invalid header size
			return nil, nil
		}
		mt, sz := MessageLineHeaderRead(this.buffer, this.rpos)
		if mp > 0 && this.rpos+size_FHEADER+sz > mp {
			return nil, fmt.Errorf("maxMessageSize %d/%d", this.rpos+size_FHEADER+sz, mp)
		}
		rp := this.rpos + size_FHEADER
		if rp+sz > this.wpos {
			return nil, nil
		}
		if mt != 0 {
			this.rpos = rp + sz
			continue
		}
		break
	}

	// read body
	r := NewMessage()
	l := this.rpos
	data := make([]byte, l)
	copy(data, this.buffer[:l])
	rp := 0
	for rp < l {
		mt, sz := MessageLineHeaderRead(data, rp)
		if mt == MLT_END {
			break
		}
		s := rp + size_FHEADER
		f := NewMessageLine(mt, data[s:s+sz])
		rp = s + sz
		r.PushBack(f)
	}
	rp = rp + size_FHEADER

	// move buffer data
	this.rpos = 0
	copy(this.buffer, this.buffer[rp:this.wpos])
	this.wpos = this.wpos - rp
	return r, nil
}
