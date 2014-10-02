package goluaparser

type Token struct {
	Kind        int
	BeginLine   int
	BeginColumn int
	EndLine     int
	EndColumn   int

	Image        string
	Value        interface{}
	Next         *Token
	SpecialToken *Token
}

func newToken(kind int) *Token {
	return newToken2(kind, "")
}

func newToken2(kind int, image string) *Token {
	switch kind {
	default:
		r := new(Token)
		r.Kind = kind
		r.Image = image
		return r
	}

}

func (this *Token) String() string {
	s := tokenImage[this.Kind]
	if this.Kind == NAME {
		return s + "(" + this.Image + ")"
	}
	return s
}
