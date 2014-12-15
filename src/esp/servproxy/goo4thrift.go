package servproxy

import (
	"bmautil/valutil"
	"fmt"
	"golua"
	"time"
)

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
		case "Write":
			return obj.write, nil
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
		case "Write":
			obj.write = valutil.ToBool(val, obj.write)
		case "Timeout":
			tm := valutil.ToInt(val, 0)
			obj.deadline = time.Now().Add(time.Duration(tm) * time.Millisecond)
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
		r["Write"] = obj.write
	}
	return r
}

func (gooThriftProxyReq) CanClose() bool {
	return false
}

func (gooThriftProxyReq) Close(o interface{}) {
}
