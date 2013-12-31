package tbus

import (
	"bmautil/byteutil"
	"bmautil/valutil"
	"esp/espnet"
	"fmt"
	"testing"
)

func tobs(s string) []byte {
	return valutil.ToBytes(s)
}

func sample(frame bool) *byteutil.BytesBuffer {
	buf := byteutil.NewBytesBuffer()
	if frame {
		buf.Add(tobs("0000005A"))
	}
	buf.Add(tobs("80"))
	buf.Add(tobs("0100010000000766756E43616C6C000000010A000100000143281B64320B0002000000"))
	buf.Add(tobs("056C6F67696E0D00030B0B00000002000000046E616D650000000671696E65726700000006706173"))
	buf.Add(tobs("7377640000000631323334353600"))
	return buf
}

func TestTMessage(t *testing.T) {
	var m TMessage
	buf := sample(false)
	fmt.Println(m.Read(buf.NewReader()))
	fmt.Printf("%v\n", m)
}

func TestChannelCoder(t *testing.T) {
	coder := NewChannelCoder(1024 * 1024)
	buf := sample(true)
	idx := 1
	for _, b := range buf.DataList {
		fmt.Println(idx, "TURN")
		idx++
		coder.DecodeMessage(nil, b, func(ev interface{}) error {
			msg := ev.(*espnet.Message)
			fmt.Printf("MESSAGE = %v\n", msg.Dump())
			return nil
		})
		fmt.Println("BUFFER", coder.buffer.TraceString(128))
		fmt.Println("FrameBody", coder.frameBody)
	}
}
