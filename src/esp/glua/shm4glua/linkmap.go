package shm4glua

import (
	"fmt"
	"sync"
	"time"
)

type mapItem struct {
	key         string
	data        interface{}
	size        int32
	expiredTime time.Time
	prev        *mapItem
	next        *mapItem
}

type LinkMap struct {
	items map[string]*mapItem
	head  *mapItem
	tail  *mapItem
	pos   *mapItem
	size  int32
	mutex sync.RWMutex
}

func (this *LinkMap) String() string {
	return fmt.Sprintf("LinkMap(%d/%d)", len(this.items), this.size)
}

func newLinkMap() *LinkMap {
	r := new(LinkMap)
	r.items = make(map[string]*mapItem)
	return r
}

func (this *LinkMap) Get(key string, tm time.Time) (interface{}, bool) {
	var val interface{}
	this.mutex.RLock()
	item, ok := this.items[key]
	if ok {
		val = item.data
		out := item.expiredTime.Before(tm)
		this.mutex.RUnlock()
		if out {
			this.mutex.Lock()
			this._remove(item)
			this.mutex.Unlock()
			return nil, false
		}
		return val, true
	}
	this.mutex.RUnlock()
	return nil, false
}

func (this *LinkMap) MGet(keys []string, tm time.Time) map[string]interface{} {
	r := make(map[string]interface{})
	this.mutex.RLock()
	for _, key := range keys {
		item, ok := this.items[key]
		if ok {
			out := item.expiredTime.Before(tm)
			if out {
				this.mutex.RUnlock()
				this.mutex.Lock()
				this._remove(item)
				this.mutex.Unlock()
				this.mutex.RLock()
			} else {
				r[key] = item.data
			}
		}
	}
	this.mutex.RUnlock()
	return r
}

func (this *LinkMap) Put(key string, val interface{}, size int32, timeoutMS int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item, ok := this.items[key]
	if ok {
		this.size -= item.size
	} else {
		item = new(mapItem)
		item.key = key
		this.items[key] = item
		if this.head == nil {
			this.head = item
		}
		if this.tail == nil {
			this.tail = item
		}
		if this.tail != nil {
			i1 := this.tail
			this.tail = item
			i1.next = item
			item.prev = i1
		}
	}
	item.data = val
	item.size = size
	this.size += size
	item.expiredTime = time.Now().Add(time.Millisecond * time.Duration(timeoutMS))
}

func (this *LinkMap) _remove(item *mapItem) bool {
	old, ok := this.items[item.key]
	if !ok || old != item {
		return false
	}
	this.size -= item.size
	delete(this.items, item.key)
	i1 := item.prev
	i2 := item.next
	if i1 != nil {
		i1.next = i2
	}
	if i2 != nil {
		i2.prev = i1
	}
	if this.head == item {
		this.head = i2
	}
	if this.tail == item {
		this.tail = i2
	}
	if this.pos == item {
		this.pos = i2
	}
	return true
}

func (this *LinkMap) Remove(key string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if item, ok := this.items[key]; ok {
		return this._remove(item)
	}
	return false
}

func (this *LinkMap) MRemove(keys []string) int {
	c := 0
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, key := range keys {
		if item, ok := this.items[key]; ok {
			if this._remove(item) {
				c = c + 1
			}
		}
	}
	return c
}

func (this *LinkMap) Clear(maxStep int) int {
	c := 0
	tm := time.Now()
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item := this.pos
	for i := 0; i < maxStep; i++ {
		if item == nil {
			item = this.head
		} else {
			item = item.next
		}
		if item == nil {
			break
		}
		if item.expiredTime.Before(tm) {
			this._remove(item)
			c = c + 1
		}
	}
	this.pos = item
	return c
}
