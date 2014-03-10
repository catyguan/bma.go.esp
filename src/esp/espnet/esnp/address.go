package esnp

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	ADDRESS_GROUP   = int(50)
	ADDRESS_HOST    = int(40)
	ADDRESS_SERVICE = int(30)
	ADDRESS_OP      = int(20)
	ADDRESS_OBJECT  = int(10)
)

type Address struct {
	pack       *Package
	coder      addrCoder
	annotation map[int]string
}

func NewAddress() *Address {
	this := new(Address)
	return this
}

func NewAddressP(pack *Package, mt byte) *Address {
	this := new(Address)
	this.pack = pack
	this.coder = addrCoder(mt)
	return this
}

func (this *Address) Annotations() []int {
	if this.pack != nil {
		return this.coder.List(this.pack)
	}
	r := make([]int, 0, len(this.annotation))
	if this.annotation != nil {
		for v, _ := range this.annotation {
			r = append(r, v)
		}
	}
	return r
}

func (this *Address) Get(ann int) string {
	if this.pack != nil {
		v, err := this.coder.Get(this.pack, ann)
		if err == nil {
			return v
		}
	} else {
		if this.annotation != nil {
			v, ok := this.annotation[ann]
			if ok {
				return v
			}
		}
	}
	return ""
}

func (this *Address) Set(ann int, val string) {
	if this.pack != nil {
		this.coder.Set(this.pack, ann, val)
	} else {
		if this.annotation != nil {
			this.annotation = make(map[int]string)
		}
		this.annotation[ann] = val
	}
}

func (this *Address) SetGroup(val string) {
	this.Set(ADDRESS_GROUP, val)
}

func (this *Address) GetGroup() string {
	return this.Get(ADDRESS_GROUP)
}

func (this *Address) SetHost(val string) {
	this.Set(ADDRESS_HOST, val)
}

func (this *Address) GetHost() string {
	return this.Get(ADDRESS_HOST)
}

func (this *Address) SetService(val string) {
	this.Set(ADDRESS_SERVICE, val)
}

func (this *Address) GetService() string {
	return this.Get(ADDRESS_SERVICE)
}

func (this *Address) SetOp(val string) {
	this.Set(ADDRESS_OP, val)
}

func (this *Address) GetOp() string {
	return this.Get(ADDRESS_OP)
}

func (this *Address) SetObject(val string) {
	this.Set(ADDRESS_OBJECT, val)
}

func (this *Address) GetObject() string {
	return this.Get(ADDRESS_OBJECT)
}

func (this *Address) SetCall(s string, op string) {
	this.SetService(s)
	this.SetOp(op)
}

func (this *Address) Remove(ann int) {
	if this.pack != nil {
		this.coder.Remove(this.pack, ann)
	} else {
		if this.annotation != nil {
			delete(this.annotation, ann)
		}
	}
}

func (this *Address) Bind(pack *Package, mt byte) {
	coder := addrCoder(mt)
	coder.Clear(this.pack)
	for ann, v := range this.annotation {
		coder.Set(this.pack, ann, v)
	}
}

func (this *Address) String() string {
	anns := this.Annotations()
	sort.Reverse(sort.IntSlice(anns))

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("Address[")
	for i, ann := range anns {
		v := this.Get(ann)
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%d=%s", ann, v))
	}
	buf.WriteString("]")
	return buf.String()
}
