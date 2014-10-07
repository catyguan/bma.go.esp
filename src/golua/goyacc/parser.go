package goyacc

import (
	"bytes"
	"fmt"
	"unicode"
)

type Parser struct {
	name         string
	stream       *Stream
	lastl, lastc int
	err          error
	chunk        Node
}

func NewParser(n string, str string) *Parser {
	r := new(Parser)
	r.name = n
	r.stream = newStream(str)
	return r
}

func (this *Parser) Error3(err interface{}, l, c int) {
	if this.err != nil {
		return
	}
	this.err = fmt.Errorf("%s(%d:%d)", err, l, c)
}

func (this *Parser) Error2(err error, tk *yyToken) {
	l := tk.line
	c := tk.column
	this.Error3(err, l, c)
}

func (this *Parser) Error(e string) {
	l := this.lastl
	c := this.lastc
	this.Error3(e, l, c)
}

func (this *Parser) Parse() (Node, error) {
	yyParse(this)
	if this.err != nil {
		return nil, this.err
	}
	return this.chunk, nil
}

//////////////////////////////////////////////////
// TOKEN /////////////////////////////////////////
// const keywordStart = FUNCTION
// var keywords []string = []string{
// 	"function", "if", "else",
// 	"while", "for", "IN",
// 	"break", "continue",
// 	"var", "new", "delete",
// 	"this",
// 	"true", "false", "null",
// 	"with", "return",
// }

const keywordStart = AND

var keywords []string = []string{
	"and", "break", "do",
	"else", "elseif", "end",
	"false", "for",
	"function", "goto",
	"if", "in", "local", "nil",
	"not", "or", "return",
	"repeat", "then", "true",
	"until", "while",
}

func (this *Parser) fillToken(lval *yySymType, k, p int) int {
	tk := &lval.token
	tk.kind = k
	tk.line = this.stream.getLine()
	tk.column = this.stream.getColumn()
	tk.image = this.stream.getImage(p, -1)
	this.lastl = tk.line
	this.lastc = tk.column
	lval.value = nil
	lval.op = OP_NONE
	// fmt.Println("fillToken", tk)
	return k
}

func (this *Parser) firstChar() (rune, int) {
	ch, sp := this.stream.readChar()
	if ch == 0 {
		return ch, sp
	}
	if unicode.IsSpace(ch) {
		for unicode.IsSpace(ch) {
			ch, sp = this.stream.readChar()
		}
		this.stream.backup(1)
		return 0, -1
	}
	// if ch == '/' {
	// 	c2, _ := this.stream.readChar()
	// 	if c2 == '/' {
	// 		this.stream.skip1('\n')
	// 		return 0, -1
	// 	}
	// 	if c2 == '*' {
	// 		this.stream.skip2('*', '/')
	// 		return 0, -1
	// 	}
	// }
	if ch == '-' {
		if this.stream.checkNext('-') {
			this.stream.skip1('\n')
			return 0, -1
		}
	}
	if ch == '[' {
		if this.stream.checkNext('[') {
			this.stream.skip2(']', ']')
			return 0, -1
		}
	}
	return ch, sp
}

func (this *Parser) Lex(lval *yySymType) int {
	if this.err != nil {
		return 0
	}
	r := this.lex(lval)
	if r != 0 {
		if yyDebug >= 3 {
			fmt.Println("lex => ", lval.token)
		}
	}
	return r
}

func isName(ch rune, fi bool) bool {
	if ch == '_' {
		return true
	}
	if !fi && unicode.IsDigit(ch) {
		return true
	}
	return unicode.IsLetter(ch)
}

func (this *Parser) lex(lval *yySymType) int {
	ch, sp := this.firstChar()
	for sp < 0 {
		ch, sp = this.firstChar()
	}
	if unicode.IsDigit(ch) {
		dot := false
		for {
			c1, _ := this.stream.readChar()
			if c1 == '.' {
				if dot {
					break
				}
				dot = true
				continue
			}
			if !unicode.IsDigit(c1) {
				break
			}
		}
		this.stream.backup(1)
		return this.fillToken(lval, NUMBER, sp)
	}
	if ch == '"' || ch == '\'' {
		buf := bytes.NewBuffer(make([]byte, 0, 32))
		for {
			c1, p := this.stream.readChar()
			if c1 == 0 {
				return 0
			}
			if c1 == ch {
				this.fillToken(lval, STRING, p)
				lval.token.image = buf.String()
				return STRING
			}
			if c1 == '\\' {
				c2, _ := this.stream.readChar()
				switch c2 {
				case 'b':
					buf.WriteByte('\b')
					continue
				case 't':
					buf.WriteByte('\t')
					continue
				case 'n':
					buf.WriteByte('\n')
					continue
				case 'f':
					buf.WriteByte('\f')
					continue
				case 'r':
					buf.WriteByte('\r')
					continue
				case '"':
					buf.WriteByte('"')
					continue
				case '\'':
					buf.WriteByte('\'')
					continue
				case '\\':
					buf.WriteByte('\\')
					continue
				}
				this.stream.backup(1)
			}
			buf.WriteRune(c1)
		}
	}
	if ch >= 0x21 && ch <= 0x3f {
		switch ch {
		// case '+':
		// 	if this.stream.checkNext('=') {
		// 		return this.fillToken(lval, SADDASS, sp)
		// 	}
		// case '-':
		// 	if this.stream.checkNext('=') {
		// 		return this.fillToken(lval, SSUBASS, sp)
		// 	}
		// case '*':
		// 	if this.stream.checkNext('=') {
		// 		return this.fillToken(lval, SMULASS, sp)
		// 	}
		// case '/':
		// 	if this.stream.checkNext('=') {
		// 		return this.fillToken(lval, SDIVASS, sp)
		// 	}
		case '<':
			if this.stream.checkNext('=') {
				return this.fillToken(lval, SLTEQ, sp)
				// } else if this.stream.checkNext('<') {
				// 	return this.fillToken(lval, SLSHIFT, sp)
			}
		case '>':
			if this.stream.checkNext('=') {
				return this.fillToken(lval, SGTEQ, sp)
				// } else if this.stream.checkNext('>') {
				// 	return this.fillToken(lval, SRSHIFT, sp)
			}
		case '=':
			if this.stream.checkNext('=') {
				// if this.stream.checkNext('=') {
				// 	return this.fillToken(lval, SEQ3, sp)
				// }
				return this.fillToken(lval, SEQ, sp)
			}
		case '!':
			if this.stream.checkNext('=') {
				return this.fillToken(lval, SNOTEQ, sp)
			}
		case '.':
			if this.stream.checkNext('.') {
				if this.stream.checkNext('.') {
					return this.fillToken(lval, MORE, sp)
				}
				return this.fillToken(lval, STRADD, sp)
			}
			// case '|':
			// 	if this.stream.checkNext('|') {
			// 		return this.fillToken(lval, SOR, sp)
			// 	}
			// case '&':
			// 	if this.stream.checkNext('&') {
			// 		return this.fillToken(lval, SAND, sp)
			// 	}
		}
		return this.fillToken(lval, int(ch), sp)
	}
	if ch < 0xFF {
		for j, kw := range keywords {
			if kw != "" && kw[0] == byte(ch) {
				c1 := ch
				m := true
				for i, kch := range kw {
					if c1 == kch {
						c1, _ = this.stream.readChar()
					} else {
						m = false
						this.stream.backup(i)
						break
					}
				}
				if m {
					this.stream.backup(1)
					return this.fillToken(lval, j+keywordStart, sp)
				}
			}
		}
	}

	if isName(ch, true) {
		for {
			c1, _ := this.stream.readChar()
			if !isName(c1, false) {
				this.stream.backup(1)
				return this.fillToken(lval, NAME, sp)
			}
		}
	}
	return this.fillToken(lval, int(ch), sp)
}
