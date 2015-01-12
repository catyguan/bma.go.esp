package vmmhttp

import (
	"bytes"
	"golua"
)

func RenderScriptPreprocess(content string) (string, error) {
	return golua.DoRenderScriptPreprocess(content, func(buf *bytes.Buffer) error {
		buf.WriteString("local out = httpserv.write\n")
		return nil
	}, nil)
}
