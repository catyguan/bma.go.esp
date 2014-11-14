package memserv

import (
	"bmautil/memblock"
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estMemGo(t *testing.T) {
	cfg := new(MemGoConfig)
	m := NewMemGo(cfg)
	m.mem.Listener = func(k string, item *memblock.MapItem, rt memblock.REMOVE_TYPE) {
		fmt.Println("remove", k, item.Data, rt)
	}
	err := m.Start()
	if err != nil {
		t.Error(err)
		return
	}
	m.goo.DoSync(func() {
		m.mem.Put("test", 1, 4, 10)
	})
	time.Sleep(100 * time.Millisecond)
	m.goo.StopAndWait()
	time.Sleep(100 * time.Millisecond)
}
