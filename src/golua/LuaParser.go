package golua

import "fmt"

const (
	VAR  = 0
	CALL = 1
)

type lookaheadSuccess struct {
}

func (this *lookaheadSuccess) Error() string {
	return "lookaheadSuccess"
}

type LuaParser struct {
	jjtree *JJTLuaParserState

	token_source    *LuaParserTokenManager
	jj_input_stream *SimpleCharStream
	/** Current token. */
	token *Token
	/** Next token. */
	jj_nt      *Token
	jj_ntk     int
	jj_scanpos *Token
	jj_lastpos *Token
	jj_la      int

	// static private final class LookaheadSuccess extends java.lang.Error { }
	// final private LookaheadSuccess jj_ls = new LookaheadSuccess();
	jj_ls *lookaheadSuccess
}

func newLuaParser() *LuaParser {
	r := new(LuaParser)
	r.jjtree = newJJTLuaParserState()
	r.jj_ls = new(lookaheadSuccess)
	return r
}

func NewLuaParser1(str string) *LuaParser {
	this := newLuaParser()
	this.jj_input_stream = NewSimpleCharStream3(str, 1, 1)
	this.token_source = newLuaParserTokenManager1(this.jj_input_stream)
	this.token = newToken(0)
	this.jj_ntk = -1
	return this
}

/** Enable tracing. */
func (this *LuaParser) enable_tracing() {
}

/** Disable tracing. */
func (this *LuaParser) disable_tracing() {
}

/** Generate ParseException. */
func (this *LuaParser) generateParseException() error {
	errortok := this.token.Next
	line := errortok.BeginLine
	column := errortok.BeginColumn
	mess := ""
	if errortok.Kind == 0 {
		mess = tokenImage[0]
	} else {
		mess = errortok.Image
	}
	err := fmt.Sprintf("Parse error at line %d, column %d.  Encountered: %s", line, column, mess)
	return newParseException(err)
}

func (this *LuaParser) jj_ntk_f() (int, error) {
	this.jj_nt = this.token.Next
	if this.jj_nt == nil {
		tk, err := this.token_source.getNextToken()
		if err != nil {
			return 0, err
		}
		this.token.Next = tk
		this.jj_ntk = tk.Kind
	} else {
		this.jj_ntk = this.jj_nt.Kind
	}
	return this.jj_ntk, nil
}

/** Get the specific Token. */
func (this *LuaParser) getToken(index int) (*Token, error) {
	t := this.token
	for i := 0; i < index; i++ {
		if t.Next != nil {
			t = t.Next
		} else {
			t2, err := this.token_source.getNextToken()
			if err != nil {
				return nil, err
			}
			t.Next = t2
			t = t2
		}
	}
	return t, nil
}

/** Get the next Token. */
func (this *LuaParser) getNextToken() (*Token, error) {
	if this.token.Next != nil {
		this.token = this.token.Next
	} else {
		t2, err := this.token_source.getNextToken()
		if err != nil {
			return nil, err
		}
		this.token.Next = t2
		this.token = t2
	}
	this.jj_ntk = -1
	return this.token, nil
}

func (this *LuaParser) jj_scan_token(kind int) (bool, error) {
	if this.jj_scanpos == this.jj_lastpos {
		this.jj_la--
		if this.jj_scanpos.Next == nil {
			t, err := this.token_source.getNextToken()
			if err != nil {
				return false, err
			}
			this.jj_scanpos.Next = t
			this.jj_scanpos = t
			this.jj_lastpos = t
		} else {
			this.jj_scanpos = this.jj_scanpos.Next
			this.jj_lastpos = this.jj_scanpos
		}
	} else {
		this.jj_scanpos = this.jj_scanpos.Next
	}
	if this.jj_scanpos.Kind != kind {
		return true, nil
	}
	if this.jj_la == 0 && this.jj_scanpos == this.jj_lastpos {
		return false, this.jj_ls
	}
	return false, nil
}

func (this *LuaParser) jj_consume_token(kind int) (*Token, error) {
	oldToken := this.token
	if oldToken.Next != nil {
		this.token = this.token.Next
	} else {
		t, err := this.token_source.getNextToken()
		if err != nil {
			return nil, err
		}
		this.token.Next = t
		this.token = t
	}
	this.jj_ntk = -1
	if this.token.Kind == kind {
		this.pushTokenNode(kind)
		return this.token, nil
	}
	this.token = oldToken
	return nil, this.generateParseException()
}

func (this *LuaParser) pushTokenNode(kind int) {
	// fmt.Println("jj_consume_token", kind, NAME, this.token, this.token.Next)
	switch kind {
	case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, FALSE,
		NIL, NOT, TRUE, NAME, NUMBER, STRING, CHARSTRING,
		74, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94,
		AND, OR:
		n := NewSimpleNodeT(this.token)
		this.jjtree.pushNode(n)
	}
}
