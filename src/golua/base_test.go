package golua

import (
	"fmt"
	"testing"
)

func TestFuncs(t *testing.T) {
	safeCall()

	if false {
		s := "abc\t中耐高温alo#"
		s2 := addEscapes(s)
		fmt.Println("addEscapes:")
		fmt.Println(s)
		fmt.Println(s2)
	}
	if false {
		s := LexicalError(false, 12, 123, 234, "afterString", 64)
		fmt.Println("LexicalError: ", s)
	}
	if false {
		fmt.Println("SimpleCharStream")
		s := NewSimpleCharStream1("hello world\na = 1")
		for {
			c := s.readChar()
			fmt.Println(c, s.line, s.column)
			if c == 0 {
				break
			}
		}
	}
	if false {
		fmt.Println("Consts")
		fmt.Println(jjstrLiteralImages)
	}
	if false {
		s := NewSimpleCharStream1("ab = 1")
		p := newLuaParserTokenManager1(s)
		for {
			tk, err := p.getNextToken()
			if tk.Kind == 0 {
				break
			}
			fmt.Println("Kind:", tk.Kind, "Image:", tk.Image, "Err:", err)
			// break
		}
	}
}
