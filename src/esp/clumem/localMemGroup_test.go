package clumem

import (
	"fmt"
	"testing"
)

func listener(action Action, groupName string, key MemKey, val interface{}, ver MemVer) {
	fmt.Println("EVENT", action, groupName, key, val, ver)
}

func TestLocalMemGroup(t *testing.T) {
	mg := newLocalMemGroup("test")

	key := MemKey{"a", "b", "c"}
	fmt.Println(mg.Get(key))
	mg.Set(MemKey{"a"}, nil, 0)
	mg.AddListener(MemKey{"a"}, "test", listener)

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
