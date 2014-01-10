package clumem

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

type localMemGroup struct {
	name  string
	root  *localMemItem
	count int64
	size  int64
}

func newLocalMemGroup(n string) *localMemGroup {
	this := new(localMemGroup)
	this.name = n
	this.root = new(localMemItem)
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

func itemAt(list []*localMemItem, idx int) *localMemItem {
	if list == nil {
		return nil
	}
	if idx < 0 {
		idx = len(list) + idx
	}
	if idx >= 0 && idx < len(list) {
		return list[idx]
	}
	return nil
}

func keyAt(list MemKey, idx int) (string, bool) {
	if list == nil {
		return "", false
	}
	if idx < 0 {
		idx = len(list) - idx
	}
	if idx >= 0 && idx < len(list) {
		return list[idx], true
	}
	return "", false
}

func (this *localMemGroup) Query(key MemKey, createIfNotExists bool) ([]*localMemItem, bool) {
	r := []*localMemItem{this.root}
	cur := this.root
	for i, k := range key {
		if cur.items == nil {
			if !createIfNotExists {
				return r, false
			}
			cur.items = make(map[string]*localMemItem)
			this.size += local_SIZE_ITEM_MAP
		}
		item, ok := cur.items[k]
		if !ok {
			if !createIfNotExists {
				return r, false
			}
			item = new(localMemItem)
			cur.items[k] = item
			cur.NextVersion()
			this.size += local_SIZE_ITEM
			this.count++

			this.invokeListener(key[:i+1], r, ACTION_NEW)
		}
		r = append(r, item)
		cur = item
	}
	return r, true
}

func (this *localMemGroup) Get(key MemKey) (*localMemItem, bool) {
	list, ok := this.Query(key, false)
	if !ok {
		return nil, false
	}
	return itemAt(list, -1), true
}

func (this *localMemGroup) CompareAndSet(key MemKey, val interface{}, sz int, ver MemVer) bool {
	list, ok := this.Query(key, false)
	if !ok {
		return false
	}
	item := itemAt(list, -1)
	if item.version != ver {
		return false
	}

	item.NextVersion()
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.invokeListener(key, list, ACTION_UPDATE)
	return true
}

func (this *localMemGroup) Set(key MemKey, val interface{}, sz int) {
	list, _ := this.Query(key, true)
	item := itemAt(list, -1)
	item.NextVersion()
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.invokeListener(key, list, ACTION_UPDATE)
}

func (this *localMemGroup) doDelete(items []*localMemItem, cur *localMemItem, key MemKey) {
	for k, item := range cur.items {
		nkey := append(key, k)
		nitems := append(items, item)

		this.doDelete(nitems, item, nkey)
		delete(cur.items, k)
		this.invokeListener(nkey, nitems, ACTION_DELETE)
	}
	cur.items = nil
	this.size -= local_SIZE_ITEM_MAP

	if cur.listeners != nil {
		cur.listeners = nil
		this.size -= local_SIZE_ITEM_MAP
	}

	cur.value = nil
	this.size -= int64(cur.size)

	this.size -= local_SIZE_ITEM
	this.count--
}

func (this *localMemGroup) Delete(key MemKey) {
	list, ok := this.Query(key, false)
	if !ok {
		return
	}
	if len(list) == 1 {
		// don't delete root
		return
	}
	item := itemAt(list, -1)
	p := itemAt(list, -2)
	skey, _ := keyAt(key, -1)

	this.doDelete(list, item, key)
	delete(p.items, skey)
	p.NextVersion()

	this.invokeListener(key, list, ACTION_DELETE)
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

func safeInvokeListener(lis IMemListener, action Action, groupName string, key MemKey, val interface{}, ver MemVer) {
	defer func() {
		recover()
	}()
	if lis != nil {
		lis(action, groupName, key, val, ver)
	}
}

func (this *localMemGroup) invokeListener(key MemKey, items []*localMemItem, action Action) {
	for _, item := range items {
		for _, lis := range item.listeners {
			safeInvokeListener(lis, action, this.name, key, item.value, item.version)
		}
	}
}
