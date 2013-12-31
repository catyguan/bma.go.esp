package byteutil

import (
	"bytes"
	"fmt"
	"io"
)

type BytesBuffer struct {
	DataList [][]byte
}

func NewBytesBuffer() *BytesBuffer {
	this := new(BytesBuffer)
	this.DataList = make([][]byte, 0)
	return this
}

func NewBytesBufferB(b []byte) *BytesBuffer {
	this := new(BytesBuffer)
	this.DataList = [][]byte{b}
	return this
}

func NewBytesBufferA(a [][]byte) *BytesBuffer {
	this := new(BytesBuffer)
	this.DataList = a
	return this
}

func (this *BytesBuffer) Add(d []byte) *BytesBuffer {
	if d != nil {
		this.DataList = append(this.DataList, d)
	}
	return this
}

func (this *BytesBuffer) AddAll(list [][]byte) *BytesBuffer {
	if list != nil {
		for _, d := range list {
			this.DataList = append(this.DataList, d)
		}
	}
	return this
}

func (this *BytesBuffer) CheckAndPop(buf []byte, size int) bool {
	ds := this.DataSize()
	if ds < size {
		return false
	}
	bp := 0
	var rm int
	for i, d := range this.DataList {
		if len(d) >= size {
			bp += copy(buf[bp:], d[:size])
			if len(d) > size {
				this.DataList[i] = d[size:]
			} else {
				rm = i + 1
			}
			this.DataList = this.DataList[rm:]
			return true
		} else {
			bp += copy(buf[bp:], d)
			size -= len(d)
			rm = i + 1
		}
	}
	panic("BUG!!")
}

func (this *BytesBuffer) PopTo(buf *BytesBuffer, size int) (int, bool) {
	var rm int
	var rd int
	done := false
	for i, d := range this.DataList {
		l := len(d)
		if l >= size {
			buf.Add(d[:size])
			rd += size
			if l == size {
				rm = i + 1
			} else {
				this.DataList[i] = d[size:]
			}
			done = true
			break
		} else {
			buf.Add(d)
			size -= l
			rd += l
			rm = i + 1
		}
	}
	if rm > 0 {
		this.DataList = this.DataList[rm:]
	}
	return rd, done
}

func (this *BytesBuffer) ToBytes() []byte {
	r := make([]byte, this.DataSize())
	if this.DataList != nil {
		bp := 0
		for _, d := range this.DataList {
			if d != nil {
				bp += copy(r[bp:], d)
			}
		}
	}
	return r
}

func (this *BytesBuffer) Len() int {
	if this.DataList == nil {
		return 0
	}
	return len(this.DataList)
}

func (this *BytesBuffer) DataSize() int {
	sz := 0
	if this.DataList != nil {
		for _, d := range this.DataList {
			if d != nil {
				sz += len(d)
			}
		}
	}
	return sz
}

func (this *BytesBuffer) String() string {
	return fmt.Sprintf("BytesBuffer(%d)", this.DataSize())
}

func (this *BytesBuffer) TraceString(l int) string {
	if this.DataList == nil {
		return ""
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	last := 8
	last2 := 0
	for _, d := range this.DataList {
		if d == nil {
			continue
		}
		if l <= 0 {
			continue
		}
		if len(d) > l {
			d = d[:l]
			l = 0
		} else {
			l -= len(d)
		}
		first := 8 - last
		d, last2 = formatTraceData(buf, d, first)
		if d == nil {
			last += last2
			if last > 8 {
				last = 8
			}
			continue
		}
		for {
			if buf.Len() > 0 {
				buf.WriteString("-")
			}
			d, last = formatTraceData(buf, d, 8)
			if d == nil {
				break
			}
		}
	}
	return buf.String()
}

func formatTraceData(buf *bytes.Buffer, data []byte, l int) ([]byte, int) {
	if data == nil || len(data) == 0 {
		return nil, 0
	}
	var d, r1 []byte
	var r2 int
	if len(data) > l {
		d = data[:l]
		r1 = data[l:]
		r2 = l
	} else {
		d = data
		r1 = nil
		r2 = len(data)

	}
	buf.WriteString(fmt.Sprintf("%X", d))
	return r1, r2
}

// READER
type BytesBufferReader struct {
	buffer     *BytesBuffer
	lpos, dpos int
}

func (this *BytesBuffer) NewReader() *BytesBufferReader {
	r := new(BytesBufferReader)
	r.buffer = this
	r.lpos = 0
	r.dpos = 0
	return r
}

func (this *BytesBufferReader) ReadByte() (byte, error) {
	for {
		if this.lpos >= len(this.buffer.DataList) {
			return 0, io.EOF
		}
		data := this.buffer.DataList[this.lpos]
		if this.dpos < len(data) {
			r := data[this.dpos]
			this.dpos++
			return r, nil
		} else {
			// move next
			this.lpos++
			this.dpos = 0
		}
	}
}

func (this *BytesBufferReader) Read(p []byte) (int, error) {
	n := 0
	l := len(p)
	for {
		if l <= 0 {
			return n, nil
		}
		if this.lpos >= len(this.buffer.DataList) {
			return n, io.EOF
		}
		data := this.buffer.DataList[this.lpos]
		sz := len(data)
		if this.dpos < sz {
			if this.dpos+l <= sz {
				// read all
				copy(p[n:], data[this.dpos:this.dpos+l])
				this.dpos += l
				n += l
				return n, nil
			} else {
				// read part
				rsz := sz - this.dpos
				copy(p[n:], data[this.dpos:])
				this.lpos++
				this.dpos = 0
				n += rsz
				l -= rsz
			}
		} else {
			// move next
			this.lpos++
			this.dpos = 0
		}
	}
}

func (this *BytesBufferReader) ReadAll() []byte {
	r := bytes.NewBuffer(make([]byte, 0))
	for {
		if this.lpos >= len(this.buffer.DataList) {
			break
		}
		data := this.buffer.DataList[this.lpos]
		if this.dpos > 0 {
			r.Write(data[this.dpos:])
		} else {
			r.Write(data)
		}
		this.lpos++
		this.dpos = 0
	}
	return r.Bytes()
}

func (this *BytesBufferReader) Seek(offset int64, whence int) (int64, error) {
	sz := int64(this.buffer.DataSize())
	var npos, pos int64
	switch whence {
	case 1:
		for i := 0; i < this.lpos; i++ {
			npos += int64(len(this.buffer.DataList[i]))
		}
		npos += int64(this.dpos) + offset
	case 2:
		npos = sz - offset
	default:
		npos = offset
	}
	this.lpos = 0
	this.dpos = 0
	if npos < 0 {
		return 0, io.EOF
	}

	for {
		if npos <= 0 {
			return pos, nil
		}
		if this.lpos >= len(this.buffer.DataList) {
			this.lpos = len(this.buffer.DataList)
			this.dpos = 0
			return int64(sz), io.EOF
		}
		data := this.buffer.DataList[this.lpos]
		l := int64(len(data))
		if pos+l < npos {
			// move forward
			pos += l
			this.lpos++
			this.dpos = 0
		} else {
			// end
			this.dpos = int(npos - pos)
			return npos, nil
		}
	}
}

// Writer
type BytesBufferWriter struct {
	buffer *BytesBuffer
	temp   *bytes.Buffer
}

func (this *BytesBuffer) NewWriter() *BytesBufferWriter {
	r := new(BytesBufferWriter)
	r.buffer = this
	r.temp = bytes.NewBuffer(make([]byte, 0))
	return r
}

func (this *BytesBufferWriter) Append(p []byte) {
	if this.temp.Len() > 0 {
		this.buffer.Add(this.temp.Bytes())
		this.temp = bytes.NewBuffer(make([]byte, 0))
	}
	this.buffer.Add(p)
}

func (this *BytesBufferWriter) Write(p []byte) (n int, err error) {
	return this.temp.Write(p)
}

func (this *BytesBufferWriter) WriteString(s string) (n int, err error) {
	return this.temp.WriteString(s)
}

func (this *BytesBufferWriter) WriteByte(c byte) error {
	return this.temp.WriteByte(c)
}

func (this *BytesBufferWriter) End() *BytesBuffer {
	if this.temp.Len() > 0 {
		this.buffer.Add(this.temp.Bytes())
	}
	return this.buffer
}
