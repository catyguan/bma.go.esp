package gom

import (
	"bytes"
	"fmt"
)

// MStructField
type MStructField struct {
	annos   *MAnnotations
	name    string
	valtype *MValType
}

func (this *MStructField) String() string {
	return fmt.Sprintf("%s:%s", this.name, this.valtype)
}

func (this *MStructField) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(this.String())
}

// MStruct
type MStruct struct {
	annos  *MAnnotations
	name   string
	fields []*MStructField
}

func (this *MStruct) String() string {
	return fmt.Sprintf("struct(%s)", this.name)
}

func (this *MStruct) Dump(buf *bytes.Buffer, prex string) {
	if this.annos != nil {
		this.annos.Dump(buf, prex)
	}
	buf.WriteString(prex)
	buf.WriteString(fmt.Sprintf("struct %s {", this.name))
	for i, f := range this.fields {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
		p2 := prex + "\t"
		f.Dump(buf, p2)
	}
	buf.WriteString("\n}\n")
}
