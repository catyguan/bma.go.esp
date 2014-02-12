package protobuf

import (
	"fmt"
	"testing"

	"code.google.com/p/goprotobuf/proto"
)

func TestProto(t *testing.T) {
	obj := new(Person)
	obj.Id = proto.Int32(123)
	obj.Name = proto.String("guanzhong")
	obj.Email = proto.String("catyguan@163.com")

	data, err := proto.Marshal(obj)
	if err != nil {
		t.Error("marshaling error: ", err)
		return
	}

	fmt.Printf("%X\n", data)

	newObj := new(Person)
	err = proto.Unmarshal(data, newObj)
	if err != nil {
		t.Error("unmarshaling error: ", err)
		return
	}

	fmt.Println("newObj", newObj)

	// Now test and newTest contain the same data.
	if obj.GetEmail() != newObj.GetEmail() {
		t.Errorf("data mismatch %q != %q", obj.GetEmail(), newObj.GetEmail())
	}
}
