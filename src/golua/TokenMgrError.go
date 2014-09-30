package golua

import (
	"bytes"
	"fmt"
)

const (
	/**
	* Lexical error occurred.
	 */
	LEXICAL_ERROR int = 0

	/**
	 * An attempt was made to create a second instance of a static token manager.
	 */
	STATIC_LEXER_ERROR int = 1

	/**
	 * Tried to change to an invalid lexical state.
	 */
	INVALID_LEXICAL_STATE int = 2

	/**
	 * Detected (and bailed out of) an infinite loop in the token manager.
	 */
	LOOP_DETECTED int = 3
)

func addEscapes(str string) string {
	retval := bytes.NewBuffer([]byte{})
	for _, c := range []rune(str) {
		switch c {
		case 0:
			continue
		case '\b':
			retval.WriteString("\\b")
			continue
		case '\t':
			retval.WriteString("\\t")
			continue
		case '\n':
			retval.WriteString("\\n")
			continue
		case '\f':
			retval.WriteString("\\f")
			continue
		case '\r':
			retval.WriteString("\\r")
			continue
		case '"':
			retval.WriteString("\\\"")
			continue
		case '\'':
			retval.WriteString("\\'")
			continue
		case '\\':
			retval.WriteString("\\\\")
			continue
		default:
			retval.WriteRune(c)
			continue
		}
	}
	return retval.String()
}

func LexicalError(EOFSeen bool, lexState int, errorLine int, errorColumn int, errorAfter string, curChar rune) string {
	eof := "<EOF>"
	if !EOFSeen {
		eof = fmt.Sprintf("\"%s\"", addEscapes(string([]rune{curChar})))
	}
	s1 := fmt.Sprintf("%s (%d)", eof, curChar)
	return fmt.Sprintf("Lexical error at line %d, column %d.  Encountered: %s after : %s", errorLine, errorColumn, s1, errorAfter)
}

type TokenMgrError struct {
	s         string
	errorCode int
}

func newTokenMgrError(message string, reason int) *TokenMgrError {
	r := new(TokenMgrError)
	r.s = message
	r.errorCode = reason
	return r
}

func newTokenMgrErrorAll(EOFSeen bool, lexState int, errorLine int, errorColumn int, errorAfter string, curChar rune, reason int) *TokenMgrError {
	s := LexicalError(EOFSeen, lexState, errorLine, errorColumn, errorAfter, curChar)
	return newTokenMgrError(s, reason)
}

func (this *TokenMgrError) String() string {
	return this.s
}

func (this *TokenMgrError) Error() string {
	return this.s
}
