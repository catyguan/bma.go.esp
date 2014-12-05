package servproxy

import (
	"bmautil/valutil"
	"fmt"
	"golua"
	"net"
)

type ThriftProxyReq struct {
	s         *Service
	conn      net.Conn
	size      int
	readed    int
	hsize     int
	name      string
	typeId    int32
	seqId     int32
	responsed bool
	oneway    bool
}

func (this *ThriftProxyReq) IsOneway() bool {
	if this.typeId == 4 {
		return true
	}
	return this.oneway
}

func (this *ThriftProxyReq) Remain() int {
	return this.size - this.readed
}

func (this *ThriftProxyReq) String() string {
	return fmt.Sprintf("[%s, %d, %d](%d:%d)", this.name, this.typeId, this.seqId, this.size, this.readed)
}

type gooThriftProxyReq int

func (gooThriftProxyReq) Get(vm *golua.VM, o interface{}, key string) (interface{}, error) {
	if obj, ok := o.(*ThriftProxyReq); ok {
		switch key {
		case "Name":
			return obj.name, nil
		case "TypeId":
			return obj.typeId, nil
		case "SeqId":
			return obj.seqId, nil
		case "Size":
			return obj.size, nil
		case "Remain":
			return obj.Remain(), nil
		case "Oneway":
			return obj.IsOneway(), nil
		case "String":
			return obj.String(), nil
		}
	}
	return nil, nil
}

func (gooThriftProxyReq) Set(vm *golua.VM, o interface{}, key string, val interface{}) error {
	if obj, ok := o.(*ThriftProxyReq); ok {
		switch key {
		case "Name":
			obj.name = valutil.ToString(val, obj.name)
		case "TypeId":
			obj.typeId = valutil.ToInt32(val, obj.typeId)
		case "SeqId":
			obj.seqId = valutil.ToInt32(val, obj.seqId)
		case "Oneway":
			obj.oneway = valutil.ToBool(val, obj.oneway)
		case "Size":
		case "Remain":
		default:
			return fmt.Errorf("unknow set(%s)", key)
		}
	}
	return nil
}

func (gooThriftProxyReq) ToMap(o interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	if obj, ok := o.(*ThriftProxyReq); ok {
		r["Name"] = obj.name
		r["TypeId"] = obj.typeId
		r["SeqId"] = obj.seqId
		r["Size"] = obj.size
		r["Remain"] = obj.Remain()
		r["Oneway"] = obj.IsOneway()
	}
	return r
}

func (gooThriftProxyReq) CanClose() bool {
	return false
}

func (gooThriftProxyReq) Close(o interface{}) {
}
