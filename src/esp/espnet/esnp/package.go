package esnp

import (
	"bmautil/byteutil"
	"bytes"
	"fmt"
)

const (
	size_READBUF = 4
	MT_END       = byte(0x00)
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

func (this *Package) PushBackList(other *Package, cloneData bool) {
	for e := other.Front(); e != nil; e.Next() {
		this.insert(e.Clone(0, cloneData), this.tail)
	}
}

func (this *Package) PushFrontList(other *Package, cloneData bool) {
	for e := other.Back(); e != nil; e = e.Prev() {
		this.insert(e.Clone(0, cloneData), nil)
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

func (this *Package) ToBytesBuffer() (*byteutil.BytesBuffer, error) {
	l := 1
	for e := this.Front(); e != nil; e = e.Next() {
		l += 1
		if e.data != nil {
			l += e.data.Len()
		}
	}
	r := byteutil.NewBytesBuffer()
	fh := FHeader{}
	for e := this.Front(); e != nil; e = e.Next() {
		dl, err := e.Data()
		if err != nil {
			return nil, err
		}
		if dl == nil {
			continue
		}
		fh.MessageType = e.MessageType()
		fh.Size = uint32(dl.DataSize())
		r.Add(fh.ToBytes())
		r.AddAll(dl.DataList)
	}
	ph := FHeader{MT_END, 0}
	r.Add(ph.ToBytes())
	return r, nil
}

// PackageReader
type PackageReader struct {
	buffer *byteutil.BytesBuffer

	hbyte []byte

	pack    *Package
	fheader FHeader
	frame   *Frame
}

func NewPackageReader() *PackageReader {
	this := new(PackageReader)
	this.hbyte = make([]byte, size_READBUF)
	this.buffer = byteutil.NewBytesBuffer()
	return this
}

func (this *PackageReader) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("[")
	if this.pack != nil {
		buf.WriteString("P")
	}
	if this.frame != nil {
		buf.WriteString("F")
	}
	buf.WriteString(fmt.Sprintf("%d", this.buffer.DataSize()))
	buf.WriteString("]")
	return buf.String()
}

func (this *PackageReader) Append(b []byte) {
	this.buffer.Add(b)
}

func (this *PackageReader) ReadPackage(mp int) (*Package, error) {
	if this.pack == nil {
		this.pack = NewPackage()
	}
	for {
		if this.frame == nil {
			// read frame header
			if !this.buffer.CheckAndPop(this.hbyte, size_FHEADER) {
				return nil, nil
			}
			this.fheader.Read(this.hbyte, 0)
			if mp > 0 && this.fheader.Size > uint32(mp) {
				return nil, fmt.Errorf("maxframe %d/%d", this.fheader.Size, mp)
			}
			this.frame = newFrameH(this.fheader)
		}

		// read frame body
		done := false
		remain := int(this.fheader.Size) - this.frame.data.DataSize()
		if remain == 0 {
			done = true
		} else {
			_, done = this.buffer.PopTo(this.frame.data, remain)
		}
		if done {
			if this.frame.MessageType() == MT_END {
				this.frame = nil
				p := this.pack
				this.pack = nil
				return p, nil
			}
			this.pack.PushBack(this.frame)
			this.frame = nil
		} else {
			return nil, nil
		}
	}
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
			fmt.Println("DELETE", e.MessageType())
			this.Remove(e)
		}
		if stop {
			return
		}
		e = ne
	}
}
