package xmemprot

import (
	"bmautil/binlog"
	"esp/espnet"
	"fmt"
	"testing"
)

func doTestReadWrite(n string, o SHObject, o2 SHObject) error {
	fmt.Println(n, "PROT", o)

	msg := espnet.NewMessage()
	o.Write(msg)
	b, err := msg.ToPackage().ToBytesBuffer()
	if err != nil {
		return err
	}
	bs := b.ToBytes()
	fmt.Println(n, "BYTES", bs)

	msg2, err2 := espnet.NewBytesMessage(bs)
	if err2 != nil {
		return err2
	}

	err3 := o2.Read(msg2)
	if err3 != nil {
		return err3
	}
	fmt.Println(n, "DATA", o2)
	return nil
}

func TestProtocolReadWrite(t *testing.T) {
	if true {
		o := new(SHRequestSlaveJoin)
		o.Group = "group"
		o.Version = binlog.BinlogVer(123)

		o2 := new(SHRequestSlaveJoin)
		err := doTestReadWrite("SHRequestSlaveJoin", o, o2)
		if err != nil {
			t.Error(err)
		}
	}

	if true {
		o := new(SHEventBinlog)
		o.Group = "group"
		o.Version = binlog.BinlogVer(123)
		o.Data = []byte{2, 3, 4}

		o2 := new(SHEventBinlog)
		err := doTestReadWrite("SHEventBinlog", o, o2)
		if err != nil {
			t.Error(err)
		}
	}

	if true {
		o := new(SHRequestGet)
		o.Init("group", MemKey{"a", "b"})

		o2 := new(SHRequestGet)
		err := doTestReadWrite("SHRequestGet", o, o2)
		if err != nil {
			t.Error(err)
		}
	}

	if true {
		o := new(SHResponseGet)
		o.Version = MemVer(123)
		o.Value = "hello"
		o.Miss = true

		o2 := new(SHResponseGet)
		err := doTestReadWrite("SHResponseGet", o, o2)
		if err != nil {
			t.Error(err)
		}
	}

	if true {
		o := new(SHRequestSet)
		o.InitCompareAndSet("group", MemKey{"a", "b"}, "hello", 5, MemVer(123))

		o2 := new(SHRequestSet)
		err := doTestReadWrite("SHRequestSet", o, o2)
		if err != nil {
			t.Error(err)
		}
	}

	if true {
		o := new(SHResponseSet)
		o.Version = MemVer(234)

		o2 := new(SHResponseSet)
		err := doTestReadWrite("SHResponseSet", o, o2)
		if err != nil {
			t.Error(err)
		}
	}
}
