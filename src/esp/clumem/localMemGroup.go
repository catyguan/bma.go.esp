package clumem

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

type localMemGroup struct {
	name    string
	root    *localMemItem
	version MemVer
	count   int64
	size    int64
}

func newLocalMemGroup(n string) *localMemGroup {
	this := new(localMemGroup)
	this.name = n
	this.root = new(localMemItem)
	this.version = MemVer(0)
	this.count = 1
	this.size += local_SIZE_ITEM
	return this
}

func (this *localMemGroup) String() string {
	return fmt.Sprintf("%s(%d:%s)", this.name, this.count, valutil.MakeSizeString(uint64(this.size)))
}

func (this *localMemGroup) Dump() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(this.String())
	buf.WriteString("\n")
	this.root.Dump("<root>", buf, 0)
	return buf.String()
}

func (this *localMemGroup) NextVersion() MemVer {
	this.version++
	if this.version == 0 {
		this.version++
	}
	return this.version
}

func (this *localMemGroup) Get(key MemKey) (*localMemItem, bool) {
	cur := this.root
	if key != nil {
		for _, k := range key {
			if cur.items != nil {
				item, ok := cur.items[k]
				if ok {
					cur = item
					continue
				}
			}
			return nil, false
		}
	}
	return cur, true
}

func (this *localMemGroup) Set(key MemKey, val interface{}, ver MemVer, sz int) {
	cur := this.root
	if key != nil {
		for _, k := range key {
			cur.version = ver
			if cur.items == nil {
				cur.items = make(map[string]*localMemItem)
				this.size += local_SIZE_ITEM_MAP
			}
			item, ok := cur.items[k]
			if !ok {
				item = new(localMemItem)
				cur.items[k] = item
				this.size += local_SIZE_ITEM
				this.count++
			}
			cur = item
		}
	}
	cur.value = val
	cur.version = ver
	cur.size = sz
	this.size += int64(sz)
}

func (this *localMemGroup) doDelete(cur *localMemItem) {
	cur.value = nil
	this.size -= int64(cur.size)
	if cur.items != nil {
		for k, item := range cur.items {
			this.doDelete(item)
			delete(cur.items, k)
		}
		cur.items = nil
		this.size -= local_SIZE_ITEM_MAP
	}
	this.size -= local_SIZE_ITEM
	this.count--
}

func (this *localMemGroup) Delete(key MemKey, ver MemVer) {
	var p *localMemItem
	cur := this.root
	curk := ""
	if key != nil {
		for _, k := range key {
			cur.version = ver
			if cur.items != nil {
				item, ok := cur.items[k]
				if ok {
					p = cur
					cur = item
					curk = k
					continue
				}
			}
			return
		}
	}
	if cur != this.root {
		this.doDelete(cur)
		delete(p.items, curk)
	}
}
