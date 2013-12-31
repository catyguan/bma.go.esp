package tbus

import (
	"bmautil/byteutil"
	"encoding/binary"
	"fmt"
)

const (
	VERSION_MASK = 0xffff0000
	VERSION_1    = 0x80010000
)

type TMessage struct {
	name   string
	typeId byte
	seqid  int32

	buffer [4]byte
}

func (this *TMessage) readI32(reader *byteutil.BytesBufferReader) (int32, bool) {
	buf := this.buffer[:4]
	_, err := reader.Read(buf)
	if err != nil {
		return 0, false
	}
	value := int32(binary.BigEndian.Uint32(buf))
	return value, true
}

func (this *TMessage) readStringBody(reader *byteutil.BytesBufferReader, size int32) (string, bool) {
	if size < 0 {
		return "", false
	}
	buf := make([]byte, size)
	_, err := reader.Read(buf)
	if err != nil {
		return "", false
	}
	return string(buf), true
}

func (this *TMessage) Read(reader *byteutil.BytesBufferReader) (bool, error) {
	size, ok := this.readI32(reader)
	if !ok {
		return false, nil
	}
	if size < 0 {
		this.typeId = byte(size & 0x0ff)
		version := int64(int64(size) & VERSION_MASK)
		if version != VERSION_1 {
			return false, fmt.Errorf("Bad version(%d) in ReadMessageBegin", version)
		}
		sz, ok2 := this.readI32(reader)
		if !ok2 {
			return false, nil
		}
		name, ok3 := this.readStringBody(reader, sz)
		if !ok3 {
			return false, nil
		}
		this.name = name
		seqId, ok4 := this.readI32(reader)
		if !ok4 {
			return false, nil
		}
		this.seqid = seqId
		return true, nil
	}

	name, ok3 := this.readStringBody(reader, size)
	if !ok3 {
		return false, nil
	}
	this.name = name

	b, err4 := reader.ReadByte()
	if err4 != nil {
		return false, nil
	}
	this.typeId = b

	seqId, ok5 := this.readI32(reader)
	if !ok5 {
		return false, nil
	}
	this.seqid = seqId
	return true, nil
}
