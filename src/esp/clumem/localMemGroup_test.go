package clumem

import (
	"fmt"
	"testing"
)

func listener(action Action, groupName string, key MemKey, val interface{}) {
	fmt.Println("EVENT", action, groupName, key, val)
}

func TestLocalMemGroup(t *testing.T) {
	mg := newLocalMemGroup("test")

	key := MemKey{"a", "b", "c"}
	fmt.Println(mg.Get(key))
	mg.Set(MemKey{"a"}, nil, MemVer(1), 0)
	mg.AddListener(MemKey{"a"}, "test", listener)

	mg.Set(key, 123, MemVer(1), 4)
	mg.Set(MemKey{"a", "b", "d"}, 234, MemVer(2), 4)
	mg.Set(MemKey{"a", "e"}, 345, MemVer(3), 4)

	fmt.Println(mg.Get(key))

	fmt.Println("----AfterSet----")
	fmt.Print(mg.Dump())

	mg.Delete(MemKey{"a", "b"}, MemVer(4))
	mg.Delete(MemKey{"a", "f"}, MemVer(5))
	fmt.Println("----AfterDelete----")
	fmt.Print(mg.Dump())
}
