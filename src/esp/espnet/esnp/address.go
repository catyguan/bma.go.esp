package esnp

import (
	"bytes"
	"fmt"
	"net"
	"sort"
)

const (
	ADDRESS_GROUP   = int(50)
	ADDRESS_HOST    = int(40)
	ADDRESS_NODE    = int(35)
	ADDRESS_SERVICE = int(30)
	ADDRESS_OP      = int(20)
	ADDRESS_OBJECT  = int(10)
)

type Address struct {
	message    *Message
	coder      mlt_address
	annotation map[int]string
}

func NewAddress() *Address {
	this := new(Address)
	return this
}

func NewAddressP(message *Message, mt byte) *Address {
	this := new(Address)
	this.message = message
	this.coder = mlt_address(mt)
	return this
}

func (this *Address) Annotations() []int {
	if this.message != nil {
		return this.coder.List(this.message)
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
	if this.message != nil {
		v, err := this.coder.Get(this.message, ann)
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
	if this.message != nil {
		this.coder.Set(this.message, ann, val)
	} else {
		if this.annotation == nil {
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

func (this *Address) CheckHost(localName string, set bool) string {
	s := this.GetHost()
	if s == "" {
		return ""
	}
	h1, p1, _ := net.SplitHostPort(s)
	h2, _, _ := net.SplitHostPort(localName)
	if h1 == "" || h1 == "HOST" {
		nh := net.JoinHostPort(h2, p1)
		if set {
			this.SetHost(nh)
		}
		return nh
	}
	return s
}

func (this *Address) GetHost() string {
	return this.Get(ADDRESS_HOST)
}

func (this *Address) SetNode(val string) {
	this.Set(ADDRESS_NODE, val)
}

func (this *Address) GetNode() string {
	return this.Get(ADDRESS_NODE)
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
	if this.message != nil {
		this.coder.Remove(this.message, ann)
	} else {
		if this.annotation != nil {
			delete(this.annotation, ann)
		}
	}
}

func (this *Address) Bind(message *Message, mt byte) {
	coder := mlt_address(mt)
	coder.Clear(this.message)
	for ann, v := range this.annotation {
		coder.Set(this.message, ann, v)
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

func (this *Address) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	var str string
	str = this.GetHost()
	if str != "" {
		r["Host"] = str
	}
	str = this.GetGroup()
	if str != "" {
		r["Group"] = str
	}
	str = this.GetService()
	if str != "" {
		r["Service"] = str
	}
	str = this.GetOp()
	if str != "" {
		r["Op"] = str
	}
	str = this.GetObject()
	if str != "" {
		r["Object"] = str
	}
	return r
}

func (this *Address) BindMap(m map[string]interface{}) {
	var v interface{}
	var ok bool
	v, ok = m["Host"]
	if ok {
		if str, ok2 := v.(string); ok2 {
			this.SetHost(str)
		}
	}
	v, ok = m["Group"]
	if ok {
		if str, ok2 := v.(string); ok2 {
			this.SetGroup(str)
		}
	}
	v, ok = m["Service"]
	if ok {
		if str, ok2 := v.(string); ok2 {
			this.SetService(str)
		}
	}
	v, ok = m["Op"]
	if ok {
		if str, ok2 := v.(string); ok2 {
			this.SetOp(str)
		}
	}
	v, ok = m["Object"]
	if ok {
		if str, ok2 := v.(string); ok2 {
			this.SetObject(str)
		}
	}
}
