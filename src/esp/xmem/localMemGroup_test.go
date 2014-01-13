package xmem

import (
	"boot"
	"bytes"
	"encoding/binary"
	"esp/sqlite"
	"fmt"
	"io/ioutil"
	"testing"
)

func listener(elist []*XMemEvent) {
	fmt.Println("EVENT ------")
	for _, ev := range elist {
		fmt.Println("\t", ev.Action, ev.GroupName, ev.Key, ev.Value, ev.Version)
	}
}

func TestLocalMemGroup(t *testing.T) {
	mg := newLocalMemGroup("test")
	mg.AddListener(MemKey{}, "test", listener)

	key := MemKey{"a", "b", "c"}
	fmt.Println(mg.Get(key))
	mg.Set(MemKey{"a"}, nil, 0)
	mg.Set(key, 123, 4)
	mg.Set(MemKey{"a", "b", "d"}, 234, 4)
	mg.Set(MemKey{"a", "e"}, 345, 4)

	fmt.Println(mg.Get(key))

	fmt.Println("----AfterSet----")
	fmt.Print(mg.Dump())

	mg.Delete(MemKey{"a", "b"})
	mg.Delete(MemKey{"a", "f"})
	fmt.Println("----AfterDelete----")
	fmt.Print(mg.Dump())
}

func TestWalk(t *testing.T) {
	mg := newLocalMemGroup("test")

	mg.Set(MemKey{"a"}, nil, 0)
	mg.Set(MemKey{"a", "b", "c"}, 123, 4)
	mg.Set(MemKey{"a", "b", "d"}, 234, 4)
	mg.Set(MemKey{"a", "e"}, 345, 4)

	fmt.Println("----Dump----")
	fmt.Print(mg.Dump())

	mg.Walk(MemKey{}, func(key MemKey, val interface{}, ver MemVer) WalkStep {
		fmt.Println(key, val, ver)
		k, _ := key.At(-1)
		if k == "b" {
			return WALK_STEP_OVER
		}
		if k == "c" {
			return WALK_STEP_OUT
		}
		return WALK_STEP_NEXT
	})
}

type testcoder int

func (O testcoder) Encode(val interface{}) (string, []byte, error) {
	if val == nil {
		return "", nil, nil
	}
	v := val.(int)
	buf := make([]byte, binary.MaxVarintLen64)
	l := binary.PutUvarint(buf, uint64(v))
	return "int", buf[:l], nil
}

func (O testcoder) Decode(flag string, data []byte) (interface{}, int, error) {
	if flag == "" {
		return nil, 0, nil
	}
	v, err := binary.ReadUvarint(bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, err
	}
	return int(v), 4, nil
}

func TestSnapshot(t *testing.T) {
	var c XMemCoder
	c = testcoder(0)
	c = SimpleCoder(0)

	mg := newLocalMemGroup("test")

	mg.Set(MemKey{"a"}, nil, 0)
	mg.Set(MemKey{"a", "b", "c"}, 123, 4)
	mg.Set(MemKey{"a", "b", "d"}, 234, 4)
	mg.Set(MemKey{"a", "e"}, 345, 4)

	fmt.Println("----Dump----")
	fmt.Print(mg.Dump())

	gss, _ := mg.Snapshot(c)
	fmt.Println("----Snapshot----")
	for _, ss := range gss.Snapshots {
		fmt.Println(ss)
	}

	mg2 := newLocalMemGroup("test2")
	mg2.AddListener(MemKey{}, "test", listener)

	mg2.BuildFromSnapshot(c, gss)
	fmt.Println("----Dump2----")
	fmt.Print(mg2.Dump())
}

func TestSaveLoad(t *testing.T) {

	cfile := "../../../bin/config/xmem-config.json"

	// sqliteServer
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	// TBusServer
	xmemService := NewService("xmemService", sqliteServer)
	boot.QuickDefine(xmemService, "", true)

	f1 := func() {
		mg := newLocalMemGroup("test")

		mg.Set(MemKey{"a"}, nil, 0)
		mg.Set(MemKey{"a", "b", "c"}, 123, 4)
		mg.Set(MemKey{"a", "b", "d"}, 234, 4)
		mg.Set(MemKey{"a", "e"}, 345, 4)

		fmt.Println("----Dump----")
		fmt.Print(mg.Dump())

		bs, _ := xmemService.doExecMemEncode("test", mg, SimpleCoder(0))
		err := ioutil.WriteFile("test.dat", bs, 0644)
		// err := xmemService.doSnapshotSave("test", bs)
		if err != nil {
			t.Error(err)
		}
	}
	if f1 != nil {

	}
	f2 := func() {
		// mg := newLocalMemGroup("test")
		// err := xmemService.doExecMemLoad("test", mg, testcoder(0))
		// if err != nil {
		// 	t.Error(err)
		// 	return
		// }
		// fmt.Println("----Dump----")
		// fmt.Print(mg.Dump())
	}
	if f2 != nil {

	}

	funl := []func(){f1}

	boot.TestGo(cfile, 2, funl)
}
