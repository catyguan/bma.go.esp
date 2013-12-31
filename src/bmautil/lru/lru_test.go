package lru

import (
	"fmt"
	"math/rand"
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
	data := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
	}
	lru := NewCache(5)
	lru.Listener = func(k string, v interface{}) {
		t.Logf("%s = %v evit", k, v)
	}
	for k, v := range data {
		lru.Put(k, v)
	}
	lru.Touch("key2")
	t.Log("init:", lru.ValidDump())

	for k, v := range data {
		pv, ok := lru.Get(k)
		if !ok || pv != v {
			t.Errorf("Get %s fail, expect %s but is %s", k, v, pv)
		}
	}

	lru.Put("key6", "val6")
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
	lru.Touch("key3")
	time.Sleep(2 * time.Second)
	if _, ok := lru.Get("key3"); ok {
		t.Errorf("timeout not work")
	}
	t.Log("timeout:", lru.ValidDump())

	lru.Put("key18", "val18")
	lru.Put("key28", "key28")
	lru.Put("key27", "val27")
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
	data := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
		"key5": "val5",
	}
	lru := NewCache(10)
	for k, v := range data {
		lru.Put(k, v)
	}
	walker := func(k string, v interface{}) bool {
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

func BenchmarkOverPut(b *testing.B) {

	// move out
	var lru = NewCache(10000)

	for i := 0; i < b.N; i++ {
		v := rand.Intn(99999-10000) + 10000
		key := fmt.Sprintf("%d", v)
		lru.Put(key, v)
	}
	fmt.Println("end", lru.MaxCollide, lru.MaxRefill)
}

func BenchmarkOverGet(b *testing.B) {

	// move out
	var lru = NewCache(100000)
	for i := 10000; i <= 99999; i++ {
		key := fmt.Sprintf("%d", i)
		lru.Put(key, i)
	}

	for i := 0; i < b.N; i++ {
		v := rand.Intn(99999-10000) + 10000
		key := fmt.Sprintf("%d", v)
		lru.Get(key)
	}

}

func BenchmarkMapGet(b *testing.B) {

	// move out
	var lru = make(map[string]int)
	for i := 10000; i <= 99999; i++ {
		key := fmt.Sprintf("%d", i)
		lru[key] = i
	}

	for i := 0; i < b.N; i++ {
		v := rand.Intn(99999-10000) + 10000
		key := fmt.Sprintf("%d", v)
		_, ok := lru[key]
		if ok {

		}
	}
}
