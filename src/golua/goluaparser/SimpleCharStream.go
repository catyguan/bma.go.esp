package goluaparser

type SimpleCharStream struct {
	bufsize    int
	available  int
	tokenBegin int
	/** Position in buffer. */
	bufpos    int
	bufline   []int
	bufcolumn []int

	column int
	line   int

	prevCharIsCR bool
	prevCharIsLF bool

	buffer          []rune
	maxNextCharInd  int
	inBuf           int
	tabSize         int
	trackLineColumn bool
}

func newSimpleCharStream() *SimpleCharStream {
	r := new(SimpleCharStream)
	r.bufpos = -1
	r.line = 1
	r.tabSize = 8
	r.trackLineColumn = true
	return r
}

func NewSimpleCharStream3(str string, startline int, startcolumn int) *SimpleCharStream {
	buffersize := len(str)

	r := newSimpleCharStream()
	r.line = startline
	r.column = startcolumn - 1
	r.available = buffersize
	r.bufsize = buffersize
	r.buffer = []rune(str)
	r.bufline = make([]int, buffersize)
	r.bufcolumn = make([]int, buffersize)
	r.maxNextCharInd = buffersize
	return r
}

func NewSimpleCharStream1(str string) *SimpleCharStream {
	return NewSimpleCharStream3(str, 1, 1)
}

func (this *SimpleCharStream) setTabSize(i int) {
	this.tabSize = i
}

func (this *SimpleCharStream) getTabSize() int {
	return this.tabSize
}

/** Start. */
func (this *SimpleCharStream) BeginToken() rune {
	this.tokenBegin = -1
	c := this.readChar()
	this.tokenBegin = this.bufpos
	return c
}

func (this *SimpleCharStream) UpdateLineColumn(c rune) {
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
		this.column += (this.tabSize - (this.column % this.tabSize))
		break
	default:
		break
	}

	this.bufline[this.bufpos] = this.line
	this.bufcolumn[this.bufpos] = this.column
}

/** Read a character. */
func (this *SimpleCharStream) readChar() rune {
	if this.inBuf > 0 {
		this.inBuf--

		this.bufpos++
		if this.bufpos == this.bufsize {
			this.bufpos = 0
		}

		return this.buffer[this.bufpos]
	}

	if this.bufpos < this.bufsize {
		this.bufpos++
	}
	if this.bufpos >= this.bufsize {
		return 0
	}
	c := this.buffer[this.bufpos]

	this.UpdateLineColumn(c)
	return c
}

/** Get token end column number. */
func (this *SimpleCharStream) getEndColumn() int {
	if this.bufpos < this.bufsize {
		return this.bufcolumn[this.bufpos]
	}
	return this.bufcolumn[this.bufpos-1]
}

/** Get token end line number. */
func (this *SimpleCharStream) getEndLine() int {
	if this.bufpos < this.bufsize {
		return this.bufline[this.bufpos]
	}
	return this.bufline[this.bufpos-1]
}

/** Get token beginning column number. */
func (this *SimpleCharStream) getBeginColumn() int {
	return this.bufcolumn[this.tokenBegin]
}

/** Get token beginning line number. */
func (this *SimpleCharStream) getBeginLine() int {
	return this.bufline[this.tokenBegin]
}

/** Backup a number of characters. */
func (this *SimpleCharStream) backup(amount int) {
	this.inBuf += amount
	this.bufpos -= amount
}

/** Get token literal value. */
func (this *SimpleCharStream) GetImage() string {
	var s string
	if this.bufpos < this.bufsize {
		s = string(this.buffer[this.tokenBegin : this.bufpos+1])
	} else {
		s = string(this.buffer[this.tokenBegin:])
	}
	return s
}

/** Get the suffix. */
func (this *SimpleCharStream) GetSuffix(l int) []rune {
	ret := make([]rune, l)
	copy(ret, this.buffer[this.bufpos-l+1:this.bufpos+1])
	return ret
}

/** Reset buffer when finished. */
func (this *SimpleCharStream) Close() {
	this.buffer = nil
	this.bufline = nil
	this.bufcolumn = nil
}

/**
 * Method to adjust line and column numbers for the start of a token.
 */
func (this *SimpleCharStream) adjustBeginLineColumn(newLine int, newCol int) {
	start := this.tokenBegin
	l := 0

	if this.bufpos >= this.tokenBegin {
		l = this.bufpos - this.tokenBegin + this.inBuf + 1
	} else {
		l = this.bufsize - this.tokenBegin + this.bufpos + 1 + this.inBuf
	}

	i := 0
	j := 0
	k := 0
	nextColDiff := 0
	columnDiff := 0

	for {
		// while (i < len && bufline[j = start % bufsize] == bufline[k = ++start % bufsize])
		j = start % this.bufsize
		start += 1
		k = start % this.bufsize
		if !(i < l && this.bufline[j] == this.bufline[k]) {
			break
		}
		this.bufline[j] = newLine
		nextColDiff = columnDiff + this.bufcolumn[k] - this.bufcolumn[j]
		this.bufcolumn[j] = newCol + columnDiff
		columnDiff = nextColDiff
		i += 1
	}

	if i < l {
		this.bufline[j] = newLine
		newLine += 1
		this.bufcolumn[j] = newCol + columnDiff

		for {
			// while (i++ < len)
			i += 1
			if !(i < l) {
				break
			}
			j = start % this.bufsize
			start += 1
			k = start % this.bufsize
			if this.bufline[j] != this.bufline[k] {
				newLine += 1
				this.bufline[j] = newLine
			} else {
				this.bufline[j] = newLine
			}
		}
	}

	this.line = this.bufline[j]
	this.column = this.bufcolumn[j]
}

func (this *SimpleCharStream) getTrackLineColumn() bool {
	return this.trackLineColumn
}

func (this *SimpleCharStream) setTrackLineColumn(tlc bool) {
	this.trackLineColumn = tlc
}
