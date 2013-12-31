package cacheserver

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	keys := []string{
		"",
		"a",
		"abc",
		strings.Repeat("a", 12812),
	}
	for _, k := range keys {
		t.Log("hash:", len(k), " -> ", HashCode(k))
	}

	set := make(map[int32]bool)
	for i := 1; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		v := int32(HashCode(key) % 16)
		if _, ok := set[v]; ok {
			t.Logf("%s => %d", key, v)
		}
		set[v] = true
	}
}

func TestPutGet(t *testing.T) {
	data := map[string][]byte{
		"key1": []byte("val1"),
		"key2": []byte("val2"),
		"key3": []byte("val3"),
		"key4": []byte("val4"),
		"key5": []byte("val5"),
	}
	lru := NewCache(5)
	lru.Listener = func(k string, v []byte) {
		t.Logf("%s = %v evit", k, v)
	}
	for k, v := range data {
		lru.Put(k, v, -1)
	}
	// lru.Touch("key2")
	t.Log("init:", lru.ValidDump())

	for k, v := range data {
		pv, ok := lru.Get(k)
		if !ok {
			t.Errorf("Get %s fail, expect %v but is %v", k, v, pv)
		}
	}

	lru.Put("key6", []byte("val6"), -1)
	if lru.Size() > 5 {
		t.Error("size limit fail")
	}
	t.Log("sizeLimit:", lru.ValidDump())

	lru.Remove("key6")
	if _, ok := lru.Get("key6"); ok {
		t.Error("remove fail")
	}
	t.Log("remove:", lru.ValidDump())

	lru.ValidTime = 1
	// lru.Touch("key3")
	time.Sleep(2 * time.Second)
	if _, ok := lru.Get("key3"); ok {
		t.Errorf("timeout not work")
	}
	t.Log("timeout:", lru.ValidDump())

	lru.Put("key18", []byte("val18"), -1)
	lru.Put("key28", []byte("key28"), -1)
	lru.Put("key27", []byte("val27"), -1)
	if _, ok := lru.Get("key27"); !ok {
		t.Error("put same hash fail")
	}
	t.Log("sameHash:", lru.ValidDump())

	lru.Remove("key18")
	if _, ok := lru.Get("key27"); !ok {
		t.Error("remove same hash fail")
	}
	t.Log("removeSameHash:", lru.ValidDump())

}

func TestWalk(t *testing.T) {
	data := map[string][]byte{
		"key1": []byte("val1"),
		"key2": []byte("val2"),
		"key3": []byte("val3"),
		"key4": []byte("val4"),
		"key5": []byte("val5"),
	}
	lru := NewCache(10)
	for k, v := range data {
		lru.Put(k, v, -1)
	}
	walker := func(k string, v []byte) bool {
		t.Log(k, "=", v)
		return true
	}
	var pos int32
	var ok bool
	for {
		if pos, ok = lru.WalkAt(pos, walker, 3); !ok {
			break
		}
	}
}

func TestUpdate(t *testing.T) {
	data := map[string][]byte{
		"key1": []byte("val1"),
		"key2": []byte("val2"),
		"key3": []byte("val3"),
		"key4": []byte("val4"),
		"key5": []byte("val5"),
	}
	lru := NewCache(10)
	for k, v := range data {
		lru.Put(k, v, -1)
	}
	t.Log(lru.Dump())
	updater := func(k string, utime int64) bool {
		t.Log(k, "=", utime)
		return true
	}
	for i := 0; i < 3; i++ {
		if updater != nil {

		}
		// t.Log(i, "scan", lru.ScanUpdate(updater, 3))
		// lru.Update("key4", []byte("val4"))
	}
}

func TestTotalUse(t *testing.T) {
	lru := NewCache(100000)
	t.Log(lru.Dump())
}

func TestRemove(t *testing.T) {
	lru := NewCache(100)
	lru.Put("test", []byte("hello"), -1)
	_, done := lru.Remove("test")
	if !done {
		t.Error("remove fail")
	}
}
