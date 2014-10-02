package goluaparser

import (
	"bytes"
	"fmt"
)

type ParseException struct {
	s string

	/**
	 * This is the last token that has been consumed successfully.  If
	 * this object has been created due to a parse error, the token
	 * followng this token will (therefore) be the first error token.
	 */
	currentToken *Token

	/**
	 * Each entry in this array is an array of integers.  Each array
	 * of integers represents a sequence of tokens (by their ordinal
	 * values) that is expected at this point of the parse.
	 */
	expectedTokenSequences [][]int

	/**
	 * This is a reference to the "tokenImage" array of the generated
	 * parser within which the parse error occurred.  This array is
	 * defined in the generated ...Constants interface.
	 */
	tokenImage []string
}

func newParseException(msg string) *ParseException {
	r := new(ParseException)
	r.s = msg
	return r
}

func newParseExceptionAll(currentToken *Token, expectedTokenSequences [][]int, tokenImage []string) *ParseException {
	eol := "\n"
	expected := bytes.NewBuffer([]byte{})
	maxSize := 0
	for i := 0; i < len(expectedTokenSequences); i++ {
		ets := expectedTokenSequences[i]
		if maxSize < len(expectedTokenSequences) {
			maxSize = len(ets)
		}
		for j := 0; j < len(ets); j++ {
			expected.WriteString(tokenImage[ets[j]])
			expected.WriteByte(' ')
		}
		if ets[len(ets)-1] != 0 {
			expected.WriteString("...")
		}
		expected.WriteString(eol)
		expected.WriteString("    ")
	}
	retval := bytes.NewBuffer([]byte{})
	retval.WriteString("Encountered \"")
	tok := currentToken.Next
	for i := 0; i < maxSize; i++ {
		if i != 0 {
			retval.WriteString(" ")
		}
		if tok.Kind == 0 {
			retval.WriteString(tokenImage[0])
			break
		}
		retval.WriteString(" ")
		retval.WriteString(tokenImage[tok.Kind])
		retval.WriteString(" \"")
		retval.WriteString(addEscapes(tok.Image))
		retval.WriteString(" \"")
		tok = tok.Next
	}
	retval.WriteString(fmt.Sprintf("\" at line %d, column %d", currentToken.Next.BeginLine, currentToken.Next.BeginColumn))
	retval.WriteString(".")
	retval.WriteString(eol)
	if len(expectedTokenSequences) == 1 {
		retval.WriteString("Was expecting:")
	} else {
		retval.WriteString("Was expecting one of:")
	}
	retval.WriteString(eol)
	retval.WriteString("    ")
	retval.WriteString(expected.String())

	r := newParseException(retval.String())
	r.currentToken = currentToken
	r.expectedTokenSequences = expectedTokenSequences
	r.tokenImage = tokenImage
	return r
}

func (this *ParseException) String() string {
	return this.s
}

func (this *ParseException) Error() string {
	return this.s
}
