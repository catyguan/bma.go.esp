package vmmhttp

import (
	"bytes"
	"strings"
)

func renderMatch(cl []rune, idx int, mcl []rune) (bool, int) {
	for i := 0; i < len(mcl); i++ {
		if i+idx < len(cl) {
			c1 := cl[i+idx]
			c2 := mcl[i]
			if c1 != c2 {
				return false, 0
			}
		} else {
			return false, 0
		}
	}
	return true, idx + len(mcl)
}

var (
	k1   = []rune("<?=")
	k2   = []rune("<?golua")
	kEnd = []rune("?>")
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

func RenderScriptPreprocess(content string) (string, error) {
	clist := []rune(content)
	l := len(clist)
	buf := bytes.NewBuffer(make([]byte, 0, 1024*4))
	word := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteString("local out = httpserv.write\n")

	status := 0
	for i := 0; i < l; i++ {
		ch := clist[i]
		switch status {
		case 0:
			if ch == '<' {
				if ok, idx := renderMatch(clist, i, k1); ok {
					status = 1
					if word.Len() > 0 {
						buf.WriteString("out(\"")
						buf.WriteString(addEscapes(word.String()))
						buf.WriteString("\")\n")
					}
					word.Reset()
					buf.WriteString("out(")
					i = idx - 1
					continue
				}
				if ok, idx := renderMatch(clist, i, k2); ok {
					status = 2
					if word.Len() > 0 {
						buf.WriteString("out(\"")
						buf.WriteString(addEscapes(word.String()))
						buf.WriteString("\")\n")
					}
					word.Reset()
					i = idx - 1
					continue
				}
			}
			word.WriteRune(ch)
		case 1: // <?=
			if ch == '?' {
				if ok, idx := renderMatch(clist, i, kEnd); ok {
					status = 0
					if word.Len() > 0 {
						str := strings.TrimSpace(word.String())
						buf.WriteString(str)
						buf.WriteString(")\n")
					}
					word.Reset()
					i = idx - 1
					continue
				}
			}
			word.WriteRune(ch)
		case 2: // <?=
			if ch == '?' {
				if ok, idx := renderMatch(clist, i, kEnd); ok {
					status = 0
					if word.Len() > 0 {
						str := strings.TrimSpace(word.String())
						buf.WriteString(str)
						buf.WriteString("\n")
					}
					word.Reset()
					i = idx - 1
					continue
				}
			}
			word.WriteRune(ch)
		}
	}
	if word.Len() > 0 {
		buf.WriteString("out(\"")
		buf.WriteString(addEscapes(word.String()))
		buf.WriteString("\")\n")
	}

	return buf.String(), nil
}
