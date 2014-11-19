package memblock

import (
	"bmautil/syncutil"
	"fmt"
	"time"
)

type REMOVE_TYPE int

const (
	RT_REMOVE       = 0
	RT_LOCAL_REMOVE = 1
	RT_CLEAR        = 2
	RT_OUT          = 3
	RT_CLOSE        = 4
)

func (o REMOVE_TYPE) String() string {
	switch o {
	case RT_CLEAR:
		return "CLEAR(2)"
	case RT_OUT:
		return "OUT(3)"
	case RT_REMOVE:
		return "REMOVE(0)"
	case RT_LOCAL_REMOVE:
		return "LOCAL_REMOVE(1)"
	case RT_CLOSE:
		return "CLOSE(4)"
	default:
		return fmt.Sprintf("UNKNOW(%d)", o)
	}
}

type RemoveListener func(key string, item *MapItem, rt REMOVE_TYPE)

type MapItem struct {
	Key         string
	Data        interface{}
	Size        int32
	ExpiredTime time.Time
	prev        *MapItem
	next        *MapItem
}

func (this *MapItem) Next() *MapItem {
	return this.next
}
func (this *MapItem) Prev() *MapItem {
	return this.prev
}

const (
	ItemSize = 6 * 8
)

type MemBlock struct {
	items    map[string]*MapItem
	head     *MapItem
	tail     *MapItem
	clearPos *MapItem
	size     int32
	mutex    syncutil.PRWMutex
	MaxCount int
	Listener RemoveListener
}

func (this *MemBlock) String() string {
	return fmt.Sprintf("MemBlock(%d/%d)", len(this.items), this.size)
}

func (this *MemBlock) Size() (int, int32) {
	return len(this.items), this.size
}

func New() *MemBlock {
	r := new(MemBlock)
	r.items = make(map[string]*MapItem)
	return r
}

func (this *MemBlock) Head() *MapItem {
	return this.head
}

func (this *MemBlock) EnableMutex() {
	this.mutex.Enable()
}

func (this *MemBlock) Get(key string) (interface{}, bool) {
	return this._get(key, nil)
}

func (this *MemBlock) GetWithTimeout(key string, tm time.Time) (interface{}, bool) {
	return this._get(key, &tm)
}

func (this *MemBlock) _get(key string, tm *time.Time) (interface{}, bool) {
	var val interface{}
	this.mutex.RLock()
	item, ok := this.items[key]
	if !ok {
		this.mutex.RUnlock()
		return nil, false
	}
	val = item.Data
	out := false
	if tm != nil && item.ExpiredTime.Unix() != 0 {
		out = item.ExpiredTime.Before(*tm)
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

func (this *MemBlock) MGet(keys []string) map[string]interface{} {
	return this._mget(keys, nil)
}

func (this *MemBlock) MGetWithTimeout(keys []string, tm time.Time) map[string]interface{} {
	return this._mget(keys, &tm)
}

func (this *MemBlock) _mget(keys []string, tm *time.Time) map[string]interface{} {
	var tmp []*MapItem
	r := make(map[string]interface{})
	this.mutex.RLock()
	for _, key := range keys {
		item, ok := this.items[key]
		if !ok {
			continue
		}
		out := false
		if tm != nil && item.ExpiredTime.Unix() != 0 {
			out = item.ExpiredTime.Before(*tm)
		}
		if out {
			tmp = append(tmp, item)
		} else {
			r[key] = item.Data
		}
	}
	this.mutex.RUnlock()

	if tmp != nil {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		for _, item := range tmp {
			this._remove(item, RT_CLEAR)
		}
	}

	return r
}

func (this *MemBlock) _lruMove(item *MapItem) {
	if this.MaxCount > 0 {
		// LRU
		if item.next != nil {
			if item.prev != nil {
				item.prev.next = item.next
			} else {
				this.head = item.next
			}
			item.next.prev = item.prev
			i1 := this.tail
			this.tail = item
			i1.next = item
			item.prev = i1
			item.next = nil
		}
	}
}

func (this *MemBlock) _utime(item *MapItem, timeoutMS int) {
	if timeoutMS > 0 {
		item.ExpiredTime = time.Now().Add(time.Millisecond * time.Duration(timeoutMS))
	} else {
		item.ExpiredTime = time.Unix(0, 0)
	}
}

func (this *MemBlock) Touch(key string, timeoutMS int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item, ok := this.items[key]
	if !ok {
		return
	}
	this._lruMove(item)
	this._utime(item, timeoutMS)
}

func (this *MemBlock) Put(key string, val interface{}, size int32, timeoutMS int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item, ok := this.items[key]
	if ok {
		this.size -= item.Size
		this._lruMove(item)
	} else {
		item = new(MapItem)
		item.Key = key
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
	item.Data = val
	item.Size = size
	this.size += size
	this._utime(item, timeoutMS)

	if this.MaxCount > 0 {
		if len(this.items) > this.MaxCount {
			this._remove(this.head, RT_OUT)
		}
	}
}

func (this *MemBlock) _remove(item *MapItem, rt REMOVE_TYPE) bool {
	old, ok := this.items[item.Key]
	if !ok || old != item {
		return false
	}
	this.size -= item.Size
	delete(this.items, item.Key)
	// fmt.Printf("before remove %p, %p, %p, %p, %p\n", this.head, this.tail, item, item.prev, item.next)
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
	if this.clearPos == item {
		this.clearPos = i2
	}
	if this.Listener != nil {
		this.Listener(item.Key, item, rt)
	}
	item.Data = nil
	// fmt.Printf("after remove %p, %p, %p, %p, %p\n", this.head, this.tail, item, item.prev, item.next)
	return true
}

func (this *MemBlock) Remove(key string, local bool) bool {
	rt := RT_REMOVE
	if local {
		rt = RT_LOCAL_REMOVE
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if item, ok := this.items[key]; ok {
		return this._remove(item, REMOVE_TYPE(rt))
	}
	return false
}

func (this *MemBlock) MRemove(keys []string, local bool) int {
	rt := RT_REMOVE
	if local {
		rt = RT_LOCAL_REMOVE
	}
	c := 0
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, key := range keys {
		if item, ok := this.items[key]; ok {
			if this._remove(item, REMOVE_TYPE(rt)) {
				c = c + 1
			}
		}
	}
	return c
}

func (this *MemBlock) Clear(maxStep int) int {
	c := 0
	tm := time.Now()
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item := this.clearPos
	for i := 0; i < maxStep; i++ {
		if item == nil {
			item = this.head
		} else {
			item = item.next
		}
		if item == nil {
			break
		}
		if item.ExpiredTime.Unix() != 0 && item.ExpiredTime.Before(tm) {
			// fmt.Println("asjdlkadjlkasjdlasd, ", item.ExpiredTime, item.ExpiredTime.Unix())
			this._remove(item, RT_CLEAR)
			c = c + 1
		}
	}
	this.clearPos = item
	return c
}

func (this *MemBlock) Count() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return len(this.items)
}

func (this *MemBlock) CloseClear(notice bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	item := this.head
	for item != nil {
		delete(this.items, item.Key)
		if notice && this.Listener != nil {
			this.Listener(item.Key, item, RT_CLOSE)
		}
		item.Data = nil
		item = item.next
	}
}
