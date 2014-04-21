package esnp

import (
	"bytes"
	"fmt"
)

// Package
type Package struct {
	head *Frame
	tail *Frame
	size int
}

func NewPackage() *Package {
	this := new(Package)
	return this
}

func (this *Package) Init() *Package {
	this.head = nil
	this.tail = nil
	this.size = 0
	return this
}

func (this *Package) FrameSize() int {
	return this.size
}

func (this *Package) Front() *Frame {
	if this.size == 0 {
		return nil
	}
	return this.head
}

func (this *Package) Back() *Frame {
	if this.size == 0 {
		return nil
	}
	return this.tail
}

func (this *Package) insert(e, at *Frame) *Frame {
	var n *Frame
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
	e.pack = this
	return e
}

func (this *Package) remove(e *Frame) *Frame {
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
	e.pack = nil
	this.size--
	return e
}

func (this *Package) Remove(e *Frame) {
	if e.pack == this {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		this.remove(e)
	}
}

func (this *Package) PushFront(f *Frame) *Frame {
	return this.insert(f, this.head)
}

func (this *Package) PushBack(f *Frame) *Frame {
	return this.insert(f, this.tail)
}

func (this *Package) InsertBefore(f, mark *Frame) *Frame {
	if mark.pack != this {
		return nil
	}
	return this.insert(f, mark.prev)
}

func (this *Package) InsertAfter(f, mark *Frame) *Frame {
	if mark.pack != this {
		return nil
	}
	return this.insert(f, mark)
}

func (this *Package) PushBackList(other *Package) {
	for e := other.Front(); e != nil; e.Next() {
		this.insert(e.Clone(0), this.tail)
	}
}

func (this *Package) PushFrontList(other *Package) {
	for e := other.Back(); e != nil; e = e.Prev() {
		this.insert(e.Clone(0), nil)
	}
}

func (this *Package) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("Package")
	buf.WriteString(fmt.Sprintf("[%d]:", this.size))
	for e := this.Front(); e != nil; e = e.Next() {
		if e != this.Front() {
			buf.WriteString(",")
		}
		buf.WriteString(e.String())
	}
	return buf.String()
}

func (this *Package) ToBytes() ([]byte, error) {
	var w BytesEncodeWriter
	for e := this.Front(); e != nil; e = e.Next() {
		err := e.Encode(&w)
		if err != nil {
			return nil, err
		}
	}
	headerWrite(&w, MT_END, 0)
	return w.ToBytes(), nil
}

// PackageReader
type PackageReader struct {
	buffer []byte
	wpos   int
	rpos   int
}

func NewPackageReader() *PackageReader {
	this := new(PackageReader)
	this.buffer = make([]byte, 1024)
	return this
}

func (this *PackageReader) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("[")
	buf.WriteString(fmt.Sprintf("R:%d/W:%d/C:%d", this.rpos, this.wpos, cap(this.buffer)))
	buf.WriteString("]")
	return buf.String()
}

func (this *PackageReader) Append(b []byte) {
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

func (this *PackageReader) ReadPackage(mp int) (*Package, error) {
	for {
		if this.rpos+size_FHEADER > this.wpos {
			// invalid header size
			return nil, nil
		}
		mt, sz := headerRead(this.buffer, this.rpos)
		if mp > 0 && this.rpos+size_FHEADER+sz > mp {
			return nil, fmt.Errorf("maxPackageSize %d/%d", this.rpos+size_FHEADER+sz, mp)
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

	// read frame body
	l := this.rpos
	data := make([]byte, l)
	copy(data, this.buffer[:l])
	r := NewPackage()
	rp := 0
	for rp < l {
		mt, sz := headerRead(data, rp)
		if mt == MT_END {
			break
		}
		s := rp + size_FHEADER
		f := NewFrame(mt, data[s:s+sz])
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

func (this *Package) FrameByType(mt byte) *Frame {
	for e := this.Front(); e != nil; e = e.Next() {
		if e.MessageType() == mt {
			return e
		}
	}
	return nil
}

func (this *Package) RemoveFrame(f func(e *Frame) (bool, bool)) {
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
