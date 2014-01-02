package tbus

import (
	"bytes"
	"thrift"
)

func SerializeException(name string, seqId int32, errmsg string) *bytes.Buffer {
	buft := thrift.NewTMemoryBuffer()
	framet := thrift.NewTFramedTransport(buft, nil)
	oprot := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(framet)

	ex := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, errmsg)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	ex.Write(oprot)
	oprot.WriteMessageEnd()

	return buft.Buffer
}
