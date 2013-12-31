package lru

import (
	"bytes"
	"errors"
	"fmt"
	// "hash/adler32"
	"sync"
	"time"
)

const (
	ITEM_SIZE = 20
	empty     = ""
)

type cacheItem struct {
	valid bool
	prev  int32
	next  int32

	key   string
	value interface{}
	hash  uint32

	accessTime int64
}

type EvictListener func(key string, val interface{})

// Cache is an LRU cache.
type Cache struct {
	entries  []cacheItem
	capacity int32
	maxsize  int32
	size     int32
	head     int32
	tail     int32

	Listener  EvictListener
	Locker    sync.Locker
	ValidTime int64

	// stats
	MaxCollide int
	MaxRefill  int
}

func HashCode(str string) uint32 {
	var val uint32 = 1
	sz := len(str)
	for i := 0; i < sz; i++ {
		val += (val * 37) + uint32(str[i])
	}
	return val
	// return adler32.Checksum([]byte(str))
}

func NewCache(sz int32) *Cache {
	this := new(Cache)
	cp := sz
	if cp < 8 {
		cp = 8
	}
	this.capacity = cp * 2
	this.entries = make([]cacheItem, this.capacity)
	this.maxsize = sz
	this.size = 0
	this.head = -1
	this.tail = -1
	return this
}

func (this *Cache) MaxSize() int32 {
	return this.maxsize
}

func (this *Cache) Size() int32 {
	return this.size
}

func (this *Cache) lock() {
	if this.Locker != nil {
		this.Locker.Lock()
	}
}

func (this *Cache) unlock() {
	if this.Locker != nil {
		this.Locker.Unlock()
	}
}

func (this *Cache) item(pos int32) *cacheItem {
	if pos < 0 {
		return nil
	}
	return &this.entries[pos]
}

func (this *Cache) Clear() {
	this.lock()
	defer this.unlock()
	var i int32
	for i = 0; i < this.capacity; i++ {
		e := &this.entries[i]
		e.valid = false
		e.key = empty
		e.value = nil
	}
	this.size = 0
	this.head = -1
	this.tail = -1
}

func (this *Cache) mask(hash uint32) int32 {
	return int32(hash % uint32(this.capacity))
}

func (this *Cache) incr(pos int32) int32 {
	pos++
	if pos < this.capacity {
		return pos
	}
	return pos % this.capacity
}

func (this *Cache) isEvictable(item *cacheItem) bool {
	if item.accessTime > 0 {
		now := time.Now()
		return item.accessTime < now.Unix()
	}
	return false
}

func (this *Cache) get(key string, touch bool) (interface{}, bool) {
	hash := HashCode(key)
	pos := this.mask(hash)
	count := this.size + 1
	outdate := false

	this.lock()
	defer this.unlock()

	ct := 0
	for ; count > 0; count-- {
		item := &this.entries[pos]

		if !item.valid {
			return nil, false
		}

		if item.key == key {
			if !touch && this.isEvictable(item) {
				outdate = true
				break
			}
			if touch {
				this.updateLru(item, pos)
			}
			return item.value, true
		}
		ct++
		if ct > this.MaxCollide {
			this.MaxCollide = ct
		}
		pos = this.incr(pos)
	}

	if outdate {
		this._remove(key, false)
	}

	return nil, false
}

func (this *Cache) Get(key string) (interface{}, bool) {
	return this.get(key, false)
}

func (this *Cache) Touch(key string) (interface{}, bool) {
	return this.get(key, true)
}

func (this *Cache) IsOverload() bool {
	return this.size >= this.maxsize
}

func (this *Cache) put(key string, v interface{}, putIfAbsent bool) (interface{}, bool) {
	hash := HashCode(key)
	pos := this.mask(hash)
	count := this.size + 1

	this.lock()
	defer this.unlock()

	for this.IsOverload() {
		item := &this.entries[this.tail]
		_, ok := this._remove(item.key, false)
		if !ok {
			panic("BUG!!")
		}
	}

	ct := 0
	for ; count > 0; count-- {
		item := &this.entries[pos]

		if !item.valid {
			item.key = key
			item.value = v
			item.hash = hash
			item.valid = true
			item.prev = -1
			if this.ValidTime > 0 {
				item.accessTime = time.Now().Unix() + this.ValidTime
			}

			this.size++
			item.next = this.head
			if this.head != -1 {
				this.entries[this.head].prev = pos
			} else {
				this.tail = pos
			}
			this.head = pos
			return nil, true
		}

		// matching item gets replaced
		if item.key == key {

			if putIfAbsent {
				return item.value, false
			}

			this.updateLru(item, pos)

			oldValue := item.value
			item.value = v
			return oldValue, true
		}

		ct++
		if ct > this.MaxCollide {
			this.MaxCollide = ct
		}
		pos = this.incr(pos)
	}

	return nil, false
}

func (this *Cache) Put(key string, v interface{}) interface{} {
	r, _ := this.put(key, v, false)
	return r
}

func (this *Cache) PutIfAbsent(key string, v interface{}) (interface{}, bool) {
	return this.put(key, v, true)
}

func (this *Cache) updateLru(item *cacheItem, pos int32) {
	prevPos := item.prev
	nextPos := item.next
	prev := this.item(prevPos)
	next := this.item(nextPos)

	if prev != nil && prev.valid {
		prev.next = nextPos

		item.prev = -1
		item.next = this.head
		this.entries[this.head].prev = pos
		this.head = pos

		if next != nil && next.valid {
			next.prev = prevPos
		} else {
			this.tail = prevPos
		}
	}

	if this.ValidTime > 0 {
		item.accessTime = time.Now().Unix() + this.ValidTime
	}
}

func (this *Cache) RemoveTail() (interface{}, bool) {
	last := this.item(this.tail)
	if last == nil || !last.valid {
		return nil, false
	}
	return this.remove(last.key, true)
}

func (this *Cache) Remove(key string) (interface{}, bool) {
	return this.remove(key, true)
}

func (this *Cache) remove(key string, quite bool) (interface{}, bool) {
	this.lock()
	defer this.unlock()
	return this._remove(key, quite)
}

func evicCall(c EvictListener, k string, v interface{}) {
	defer func() {
		recover()
	}()
	c(k, v)
}

func (this *Cache) _remove(key string, quite bool) (interface{}, bool) {
	hash := HashCode(key)
	pos := this.mask(hash)
	count := this.size + 1

	for ; count > 0; count-- {
		item := &this.entries[pos]

		if !item.valid {
			return nil, false
		}

		if item.key == key {

			if !quite && this.Listener != nil {
				evicCall(this.Listener, item.key, item.value)
			}

			item.valid = false
			r := item.value
			item.key = empty
			item.value = nil
			this.size--

			prev := this.item(item.prev)
			next := this.item(item.next)

			if prev != nil && prev.valid {
				prev.next = item.next
			} else {
				this.head = item.next
			}

			if next != nil && next.valid {
				next.prev = item.prev
			} else {
				this.tail = item.prev
			}

			// Shift colliding entries down
			// fmt.Println("shift", pos, item.key, count)
			var i int32
			cpos := pos
			for i = 0; i < count; i++ {
				cpos = this.incr(cpos)
				nextItem := &this.entries[cpos]
				if !nextItem.valid {
					break
				}
				// fmt.Println("refill", nextPos, nextItem.key)
				if this.refillEntry(cpos, nextItem) {
					nextItem.valid = false
					nextItem.key = empty
					nextItem.value = nil
				}
				if i+1 > int32(this.MaxRefill) {
					this.MaxRefill = int(i + 1)
				}
			}
			// fmt.Println("shift", pos, item.key, "end")

			return r, true
		}
		pos = this.incr(pos)
	}

	if count < 0 {
		panic(errors.New("internal cache error"))
	}
	return nil, false
}

func (this *Cache) refillEntry(cpos int32, item *cacheItem) bool {
	pos := this.mask(item.hash)

	var count int32
	for count = 0; count < this.size+1; count++ {
		if pos == cpos {
			return false
		}
		nitem := &this.entries[pos]
		if !nitem.valid {
			nitem.hash = item.hash
			nitem.valid = true
			nitem.key = item.key
			nitem.next = item.next
			nitem.prev = item.prev
			nitem.accessTime = item.accessTime
			nitem.value = item.value

			prev := this.item(item.prev)
			if prev != nil && prev.valid {
				prev.next = pos
			} else {
				this.head = pos
			}
			next := this.item(item.next)
			if next != nil && next.valid {
				next.prev = pos
			} else {
				this.tail = pos
			}
			return true
		}
		pos = this.incr(pos)
	}
	return false
}

type Walker func(key string, val interface{}) bool

func (this *Cache) CalPosition(key string) int32 {
	return this.mask(HashCode(key))
}

func (this *Cache) Walk(walker Walker, step int32) (int32, bool) {
	return this.WalkAt(0, walker, step)
}

func (this *Cache) WalkAt(pos int32, walker Walker, step int32) (int32, bool) {
	var i int32
	for i = 0; i < step; i++ {
		npos := pos + i
		if npos < this.capacity {
			item := &this.entries[npos]
			if item.valid && !walker(item.key, item.value) {
				return npos, true
			}
		} else {
			return npos, false
		}
	}
	return pos + i, true
}

func (this *Cache) ValidDump() string {
	err := this.ValidLink()
	if err != nil {
		return "ValidFail:" + err.Error() + "\n" + this.Dump()
	}
	return this.Dump()
}

func (this *Cache) Dump() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(fmt.Sprintf("head=%d, tail=%d, size=%d, MaxCollide=%d, MaxRefill=%d\n", this.head, this.tail, this.size, this.MaxCollide, this.MaxRefill))
	for i, item := range this.entries {
		if item.valid {
			buf.WriteString(fmt.Sprintf("%d: %v", i, item))
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (this *Cache) ValidLink() error {
	vfun := func(pos int32, item *cacheItem, valid bool) error {
		if !item.valid {
			if valid {
				return errors.New(fmt.Sprintf("%d item not valid", pos))
			}
			return nil
		}
		if item.prev == -1 {
			if this.head != pos {
				return errors.New(fmt.Sprintf("head wrong %d -> %d", this.head, pos))
			}
		} else {
			p := this.item(item.prev)
			if p.next != pos {
				return errors.New(fmt.Sprintf("%d item prev wrong %d:%d -> %d:%d", pos, pos, item.prev, item.prev, p.next))
			}
		}
		if item.next == -1 {
			if this.tail != pos {
				return errors.New(fmt.Sprintf("tail wrong %d -> %d", this.tail, pos))
			}
		} else {
			p := this.item(item.next)
			if p.prev != pos {
				return errors.New(fmt.Sprintf("%d item next wrong %d:%d -> %d:%d", pos, pos, item.next, item.next, p.prev))
			}
		}
		return nil
	}
	i1 := this.item(this.head)
	if i1 != nil {
		err := vfun(this.head, i1, true)
		if err != nil {
			return err
		}
	}
	i2 := this.item(this.tail)
	if i2 != nil {
		err := vfun(this.tail, i2, true)
		if err != nil {
			return err
		}
	}
	for i, item := range this.entries {
		err := vfun(int32(i), &item, false)
		if err != nil {
			return err
		}
	}
	return nil
}
