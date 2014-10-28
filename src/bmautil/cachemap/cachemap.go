package linkmap

import (
	"fmt"
	"sync"
	"time"
)

type REMOVE_TYPE int

const (
	RT_REMOVE = 0
	RT_CLEAR  = 1
	RT_OUT    = 2
)

func (o REMOVE_TYPE) String() string {
	switch o {
	case RT_CLEAR:
		return "CLEAR(1)"
	case RT_OUT:
		return "OUT(2)"
	case RT_REMOVE:
		return "REMOVE(0)"
	default:
		return fmt.Sprintf("UNKNOW(%d)", o)
	}
}

type RemoveListener func(key string, clear REMOVE_TYPE)

type mapItem struct {
	key         string
	data        interface{}
	size        int32
	expiredTime time.Time
	prev        *mapItem
	next        *mapItem
}

type CacheMap struct {
	items    map[string]*mapItem
	head     *mapItem
	tail     *mapItem
	pos      *mapItem
	size     int32
	mutex    sync.RWMutex
	MaxCount int
	Listener RemoveListener
}

func (this *CacheMap) String() string {
	return fmt.Sprintf("CacheMap(%d/%d)", len(this.items), this.size)
}

func New() *CacheMap {
	r := new(CacheMap)
	r.items = make(map[string]*mapItem)
	return r
}

func (this *CacheMap) Get(key string) (interface{}, bool) {
	return this._get(key, nil)
}

func (this *CacheMap) GetWithTimeout(key string, tm time.Time) (interface{}, bool) {
	return this._get(key, &tm)
}

func (this *CacheMap) _get(key string, tm *time.Time) (interface{}, bool) {
	var val interface{}
	this.mutex.RLock()
	item, ok := this.items[key]
	if !ok {
		this.mutex.RUnlock()
		return nil, false
	}
	val = item.data
	out := false
	if tm != nil {
		out = item.expiredTime.Before(*tm)
	}
	this.mutex.RUnlock()
	if out {
		this.mutex.Lock()
		this._remove(item, RT_CLEAR)
		this.mutex.Unlock()
		return nil, false
	}
	return val, true
}

func (this *CacheMap) MGet(keys []string) map[string]interface{} {
	return this._mget(keys, nil)
}

func (this *CacheMap) MGetWithTimeout(keys []string, tm time.Time) map[string]interface{} {
	return this._mget(keys, &tm)
}

func (this *CacheMap) _mget(keys []string, tm *time.Time) map[string]interface{} {
	r := make(map[string]interface{})
	this.mutex.RLock()
	for _, key := range keys {
		item, ok := this.items[key]
		if !ok {
			continue
		}
		out := false
		if tm != nil {
			out = item.expiredTime.Before(*tm)
		}
		if out {
			this.mutex.RUnlock()
			this.mutex.Lock()
			this._remove(item, RT_CLEAR)
			this.mutex.Unlock()
			this.mutex.RLock()
		} else {
			r[key] = item.data
		}
	}
	this.mutex.RUnlock()
	return r
}

func (this *CacheMap) Put(key string, val interface{}, size int32, timeoutMS int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item, ok := this.items[key]
	if ok {
		this.size -= item.size
		if this.MaxCount > 0 {
			// LRU
			if item.prev != nil {
				item.prev.next = item.next
			} else {
				this.head = item.next
			}
			if item.next != nil {
				item.next.prev = item.prev
				i1 := this.tail
				this.tail = item
				i1.next = item
				item.prev = i1
				item.next = nil
			}
		}
	} else {
		item = new(mapItem)
		item.key = key
		this.items[key] = item
		if this.tail == nil {
			this.head = item
			this.tail = item
		} else {
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

	if this.MaxCount > 0 {
		if len(this.items) > this.MaxCount {
			this._remove(this.head, RT_OUT)
		}
	}
}

func (this *CacheMap) _remove(item *mapItem, clear REMOVE_TYPE) bool {
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
	if this.Listener != nil {
		this.Listener(item.key, clear)
	}
	return true
}

func (this *CacheMap) Remove(key string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if item, ok := this.items[key]; ok {
		return this._remove(item, RT_REMOVE)
	}
	return false
}

func (this *CacheMap) MRemove(keys []string) int {
	c := 0
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, key := range keys {
		if item, ok := this.items[key]; ok {
			if this._remove(item, RT_REMOVE) {
				c = c + 1
			}
		}
	}
	return c
}

func (this *CacheMap) Clear(maxStep int) int {
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
			this._remove(item, RT_CLEAR)
			c = c + 1
		}
	}
	this.pos = item
	return c
}

func (this *CacheMap) Count() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return len(this.items)
}
