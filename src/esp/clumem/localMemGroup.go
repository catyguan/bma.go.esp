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
	this.root.Dump("<root>", buf, 0, true)
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
	var p *localMemItem
	cur := this.root
	if key != nil {
		for i, k := range key {
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

				this.invokeListener(key[:i+1], cur, ACTION_NEW)
			}
			p = cur
			cur = item
		}
	}
	cur.value = val
	cur.version = ver
	cur.size = sz
	this.size += int64(sz)

	this.invokeListener(key, cur, ACTION_UPDATE)
	if p != nil {
		this.invokeListener(key, p, ACTION_UPDATE)
	}
}

func (this *localMemGroup) doDelete(cur *localMemItem, key MemKey) {
	cur.value = nil
	this.size -= int64(cur.size)
	if cur.items != nil {
		for k, item := range cur.items {
			nkey := append(key, k)
			this.invokeListener(nkey, item, ACTION_DELETE)
			this.invokeListener(nkey, cur, ACTION_DELETE)
			this.doDelete(item, nkey)
			delete(cur.items, k)
		}
		cur.items = nil
		this.size -= local_SIZE_ITEM_MAP
	}
	if cur.listeners != nil {
		cur.listeners = nil
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
		this.invokeListener(key, cur, ACTION_DELETE)
		if p != nil {
			this.invokeListener(key, p, ACTION_DELETE)
		}
		this.doDelete(cur, key)
		delete(p.items, curk)
	}
}

func (this *localMemGroup) AddListener(key MemKey, id string, lis IMemListener) bool {
	item, ok := this.Get(key)
	if !ok {
		return false
	}
	if item.listeners == nil {
		item.listeners = make(map[string]IMemListener)
		this.size += local_SIZE_ITEM_MAP
	}
	item.listeners[id] = lis
	return true
}

func (this *localMemGroup) RemoveListener(key MemKey, id string) {
	item, ok := this.Get(key)
	if !ok {
		return
	}
	if item.listeners != nil {
		delete(item.listeners, id)
	}
}

func safeInvokeListener(lis IMemListener, action Action, groupName string, key MemKey, val interface{}) {
	defer func() {
		recover()
	}()
	if lis != nil {
		lis(action, groupName, key, val)
	}
}

func (this *localMemGroup) invokeListener(key MemKey, item *localMemItem, action Action) {
	if item.listeners != nil {
		for _, lis := range item.listeners {
			safeInvokeListener(lis, action, this.name, key, item.value)
		}
	}
}
