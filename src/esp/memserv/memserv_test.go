package memserv

import (
	"bmautil/memblock"
	"fmt"
	"os"
	"testing"
	"time"
)

func safeCall() {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func T2estMemGo(t *testing.T) {
	cfg := new(MemGoConfig)
	m := NewMemGo("tests", cfg)
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

func TestMemGoPro(t *testing.T) {
	safeCall()

	cfg := new(MemGoConfig)
	cfg.ClearStep = 5
	m := NewMemGo("tests", cfg)
	m.mem.Listener = func(k string, item *memblock.MapItem, rt memblock.REMOVE_TYPE) {
		v := item.Data
		fmt.Println("remove", k, v, rt)
	}
	err := m.Start()
	if err != nil {
		t.Error(err)
		return
	}
	m.Set("test", 1, 10)
	m.Set("test2", "abcdef", 10)
	m1 := make(map[string]interface{})
	m1["f1"] = true
	m1["f2"] = "hello kitty"
	m.Set("test3", m1, 200)
	m.Remove("test2")
	fmt.Println(m.Size())
	time.Sleep(100 * time.Millisecond)
	m.goo.StopAndWait()
	time.Sleep(100 * time.Millisecond)
}
