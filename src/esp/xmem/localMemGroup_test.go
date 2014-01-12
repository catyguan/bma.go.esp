package xmem

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

type coder int

func (O coder) Encode(val interface{}) (string, []byte, error) {
	if val == nil {
		return "", nil, nil
	}
	v := val.(int)
	buf := make([]byte, binary.MaxVarintLen64)
	l := binary.PutUvarint(buf, uint64(v))
	return "int", buf[:l], nil
}

func (O coder) Decode(flag string, data []byte) (interface{}, int, error) {
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
	mg := newLocalMemGroup("test")

	mg.Set(MemKey{"a"}, nil, 0)
	mg.Set(MemKey{"a", "b", "c"}, 123, 4)
	mg.Set(MemKey{"a", "b", "d"}, 234, 4)
	mg.Set(MemKey{"a", "e"}, 345, 4)

	fmt.Println("----Dump----")
	fmt.Print(mg.Dump())

	slist, _ := mg.Snapshot(coder(0))
	fmt.Println("----Snapshot----")
	for _, ss := range slist {
		fmt.Println(ss)
	}

	mg2 := newLocalMemGroup("test2")
	mg2.AddListener(MemKey{}, "test", listener)

	mg2.BuildFromSnapshot(coder(0), slist)
	fmt.Println("----Dump2----")
	fmt.Print(mg2.Dump())
}
