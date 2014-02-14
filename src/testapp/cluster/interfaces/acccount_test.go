package interfaces

import (
	"fmt"
	"testing"

	"bmautil/byteutil"
	"code.google.com/p/goprotobuf/proto"
)

func TestAccountCoder(t *testing.T) {
	obj := new(Req4AccountModify)
	obj.Name = proto.String("hello")
	obj.Value = proto.Int64(1234)

	coder := OpCoder4Account(0)
	buf := byteutil.NewBytesBuffer()
	w := buf.NewWriter()
	err := coder.Encode(w, obj)
	w.End()
	if err != nil {
		t.Error("encode error: ", err)
		return
	}

	fmt.Printf("%X\n", buf.ToBytes())

	r := buf.NewReader()
	newObj, err2 := coder.Decode(r)
	if err2 != nil {
		t.Error("decode error: ", err2)
		return
	}

	fmt.Println("newObj", newObj)

}
