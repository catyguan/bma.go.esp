package gom

import "bytes"

const (
	tag = "gom"
)

type SupportDump interface {
	Dump(buf *bytes.Buffer, prex string)
}
