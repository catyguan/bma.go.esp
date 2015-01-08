package goyacc

import (
	"bytes"
)

type Stream struct {
	bufpos    int
	bufsize   int
	buffer    []rune
	bufline   []int
	bufcolumn []int

	column int
	line   int

	prevCharIsCR bool
	prevCharIsLF bool
}

func newStream(str string) *Stream {
	r := new(Stream)
	r.bufpos = -1
	r.line = 1
	r.column = 1
	r.buffer = []rune(str)
	buffersize := len(r.buffer)
	r.bufsize = buffersize
	r.bufline = make([]int, buffersize)
	r.bufcolumn = make([]int, buffersize)

	return r
}

func (this *Stream) updateLineColumn(c rune) {
	if this.bufline[this.bufpos] != 0 {
		return
	}
	this.column++

	if this.prevCharIsLF {
		this.prevCharIsLF = false
		this.column = 1
		this.line += 1
	} else if this.prevCharIsCR {
		this.prevCharIsCR = false
		if c == '\n' {
			this.prevCharIsLF = true
		} else {
			this.column = 1
			this.line += 1
		}
	}

	switch c {
	case '\r':
		this.prevCharIsCR = true
		break
	case '\n':
		this.prevCharIsLF = true
		break
	case '\t':
		this.column--
		tabSize := 1
		this.column += (tabSize - (this.column % tabSize))
		break
	default:
		break
	}

	this.bufline[this.bufpos] = this.line
	this.bufcolumn[this.bufpos] = this.column
}

func (this *Stream) readChar() (rune, int) {
	if this.bufpos < this.bufsize {
		this.bufpos++
	}
	if this.bufpos >= this.bufsize {
		return 0, this.bufpos
	}
	c := this.buffer[this.bufpos]

	this.updateLineColumn(c)
	return c, this.bufpos
}

func (this *Stream) checkNext(ch rune) bool {
	c, _ := this.readChar()
	if c == ch {
		return true
	}
	this.backup(1)
	return false
}

func (this *Stream) skip1(c1 rune) {
	for {
		ch, _ := this.readChar()
		if ch == 0 {
			return
		}
		if ch == c1 {
			return
		}
	}
}

func (this *Stream) skip2(c1 rune, c2 rune) {
	for {
		ch, _ := this.readChar()
		if ch == 0 {
			return
		}
		if ch == c1 {
			c0, _ := this.readChar()
			if c0 == c2 {
				return
			}
			this.backup(1)
		}
	}
}

func (this *Stream) keepSkip2(c1 rune, c2 rune) string {
	buf := bytes.NewBuffer([]byte{})
	for {
		ch, _ := this.readChar()
		if ch == 0 {
			break
		}
		if ch == c1 {
			c0, _ := this.readChar()
			if c0 == c2 {
				break
			}
			this.backup(1)
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}

func (this *Stream) getColumn() int {
	if this.bufpos < this.bufsize {
		return this.bufcolumn[this.bufpos]
	}
	if this.bufpos > 0 && this.bufpos-1 < len(this.bufcolumn) {
		return this.bufcolumn[this.bufpos-1]
	}
	return 0

}

func (this *Stream) getLine() int {
	// print(this.bufpos, this.bufsize, len(this.bufline))
	if this.bufpos < this.bufsize {
		return this.bufline[this.bufpos]
	}
	if this.bufpos > 0 && this.bufpos-1 < len(this.bufline) {
		return this.bufline[this.bufpos-1]
	}
	return 0
}

func (this *Stream) backup(amount int) {
	// fmt.Println("backup", amount, " at ", this.bufpos)
	this.bufpos -= amount
}

func (this *Stream) getImage(p1, p2 int) string {
	if p1 > this.bufsize {
		p1 = this.bufsize
	}
	if p2 == -1 {
		p2 = this.bufpos + 1
	}
	if p2 > this.bufsize {
		p2 = this.bufsize
	}
	// fmt.Println(this.bufsize, this.bufpos, p1, p2)
	return string(this.buffer[p1:p2])
}
