package gom

import "bytes"

type MValType struct {
	name       string
	innerType1 *MValType
	innerType2 *MValType
}

func (this *MValType) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	buf.WriteString(this.name)
	if this.innerType1 != nil {
		buf.WriteByte('<')
		buf.WriteString(this.innerType1.String())
		if this.innerType2 != nil {
			buf.WriteByte(',')
			buf.WriteString(this.innerType2.String())
		}
		buf.WriteByte('>')
	}
	return buf.String()
}

type MValue struct {
	annos *MAnnotations
	value interface{}
}
