package tbus

import (
	"bytes"
	"esp/espnet"
	"fmt"
	"thrift"
)

func SerializeReplyException(msg *espnet.Message, errmsg string) []byte {
	hs := msg.Headers()
	name, e1 := hs.GetString(THRIFT_TMESSAGE_NAME, "")
	if e1 != nil {
		return nil
	}
	seqId, e2 := hs.GetInt(THRIFT_TMESSAGE_SEQ, 0)
	if e2 != nil {
		return nil
	}
	bs := SerializeException(name, int32(seqId), errmsg)
	fmt.Println(bs.Bytes())
	return bs.Bytes()
}

func SerializeException(name string, seqId int32, errmsg string) *bytes.Buffer {
	buft := thrift.NewTMemoryBuffer()
	framet := thrift.NewTFramedTransport(buft, nil)
	oprot := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(framet)

	ex := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, errmsg)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	ex.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()

	return buft.Buffer
}
