package shm4glua

import (
	"fmt"
	"testing"
	"time"
)

func TestLinkMap(t *testing.T) {
	m := newLinkMap()
	m.Put("test", 123, 3, 100)
	fmt.Printf("map = %s\n", m)

	m.Put("test", 234, 4, 100)
	fmt.Printf("map = %s\n", m)

	m.Put("test2", 345, 5, 200)
	for i := 0; i < 15; i++ {
		m.Put(fmt.Sprintf("test3-%d", i), 345, 6, 300)
	}
	m.Put("test4", 345, 5, 200)
	fmt.Printf("map = %s\n", m)

	val, ok := m.Get("test", time.Now())
	fmt.Printf("Get -> %v, %v\n", ok, val)
	fmt.Printf("MGet -> %v\n", m.MGet([]string{"test", "test2", "test4"}, time.Now()))

	time.Sleep(101 * time.Millisecond)
	_, ok2 := m.Get("test", time.Now())
	fmt.Printf("Get -> %v --> %s\n", ok2, m)
	fmt.Printf("MGet -> %v\n", m.MGet([]string{"test", "test2", "test4"}, time.Now()))

	c1 := m.Clear(10)
	fmt.Printf("Clear -> %s, %d\n", m, c1)
	time.Sleep(101 * time.Millisecond)
	c2 := m.Clear(10)
	fmt.Printf("Clear -> %s, %d\n", m, c2)
}
