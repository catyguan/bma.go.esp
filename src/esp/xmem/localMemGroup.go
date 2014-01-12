package xmem

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

// localMemLis
type localMemLis struct {
	items     map[string]*localMemLis
	listeners map[string]XMemListener
}

func (this *localMemLis) allInvokeListener(elist []*XMemEvent) {
	for _, lis := range this.listeners {
		safeInvokeListener(lis, elist)
	}
	for _, item := range this.items {
		item.allInvokeListener(elist)
	}
}

// localActionContext
type localActionContext struct {
	listeners map[*localMemLis][]*XMemEvent
}

func (this *localActionContext) Add(item *localMemLis, ev *XMemEvent) {
	if this.listeners == nil {
		this.listeners = make(map[*localMemLis][]*XMemEvent)
	}
	elist, ok := this.listeners[item]
	if !ok {
		elist = []*XMemEvent{ev}
	} else {
		elist = append(elist, ev)
	}
	this.listeners[item] = elist
}

func (this *localActionContext) Invoke() {
	for item, elist := range this.listeners {
		for _, lis := range item.listeners {
			safeInvokeListener(lis, elist)
		}
	}
}

// localMemGroup
type localMemGroup struct {
	name    string
	root    *localMemItem
	lisRoot *localMemLis
	count   int64
	size    int64
	version MemVer
}

func newLocalMemGroup(n string) *localMemGroup {
	this := new(localMemGroup)
	this.name = n
	this.root = new(localMemItem)
	this.lisRoot = new(localMemLis)
	this.count = 1
	this.size = local_SIZE_ITEM
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

func (this *localMemGroup) Walk(key MemKey, w XMemWalker) bool {
	item, ok := this.Get(key)
	if ok {
		item.Walk(key, w)
		return true
	}
	return false
}

func (this *localMemGroup) Snapshot(coder XMemCoder) ([]*XMemSnapshot, error) {
	var rerr error
	r := make([]*XMemSnapshot, 0, this.count)
	this.root.Walk(MemKey{}, func(key MemKey, val interface{}, ver MemVer) WalkStep {
		k, bs, err := coder.Encode(val)
		if err != nil {
			rerr = err
			return WALK_STEP_END
		}
		ss := new(XMemSnapshot)
		ss.Data = bs
		ss.Key = key.ToString()
		ss.Kind = k
		ss.Version = ver
		r = append(r, ss)
		return WALK_STEP_NEXT
	})
	return r, rerr
}

func (this *localMemGroup) BuildFromSnapshot(coder XMemCoder, slist []*XMemSnapshot) error {
	this.root.Clear()

	ev := new(XMemEvent)
	ev.Action = ACTION_CLEAR
	ev.GroupName = this.name
	ev.Key = MemKey{}
	ev.Value = nil
	ev.Version = MemVer(0)
	this.lisRoot.allInvokeListener([]*XMemEvent{ev})

	this.root = new(localMemItem)
	this.count = 1
	this.size = local_SIZE_ITEM

	var ctx localActionContext
	for _, ss := range slist {
		val, sz, err := coder.Decode(ss.Kind, ss.Data)
		if err != nil {
			return err
		}
		this.InitSet(&ctx, MemKeyFromString(ss.Key), val, sz, ss.Version)
	}
	this.version = this.root.version

	ctx.Invoke()

	return nil
}

func (this *localMemGroup) NextVersion() MemVer {
	this.version++
	if this.version == VERSION_INVALID {
		this.version++
	}
	return this.version
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

func (this *localMemGroup) Query(key MemKey, ctx *localActionContext) ([]*localMemItem, bool, MemVer) {
	r := []*localMemItem{this.root}
	cur := this.root
	ver := VERSION_INVALID
	for i, k := range key {
		if cur.items == nil {
			if ctx == nil {
				return r, false, ver
			}
			cur.items = make(map[string]*localMemItem)
			this.size += local_SIZE_ITEM_MAP
		}
		item, ok := cur.items[k]
		if !ok {
			if ctx == nil {
				return r, false, ver
			}
			if ver == VERSION_INVALID {
				ver = this.NextVersion()
			}
			item = new(localMemItem)
			item.version = ver

			cur.items[k] = item
			this.size += local_SIZE_ITEM
			this.count++

			this.tryInvokeListener(ctx, ACTION_NEW, key[:i+1], item)
		}
		r = append(r, item)
		cur = item
	}
	return r, true, ver
}

func (this *localMemGroup) LisQuery(key MemKey, createIfNotExists bool) ([]*localMemLis, bool) {
	r := []*localMemLis{this.lisRoot}
	cur := this.lisRoot
	for _, k := range key {
		if cur.items == nil {
			if !createIfNotExists {
				return r, false
			}
			cur.items = make(map[string]*localMemLis)
		}
		item, ok := cur.items[k]
		if !ok {
			if !createIfNotExists {
				return r, false
			}
			item = new(localMemLis)

			cur.items[k] = item
		}
		r = append(r, item)
		cur = item
	}
	return r, true
}

func (this *localMemGroup) Get(key MemKey) (*localMemItem, bool) {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return nil, false
	}
	return itemAt(list, -1), true
}

func (this *localMemGroup) CompareAndSet(key MemKey, val interface{}, sz int, ver MemVer) MemVer {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return VERSION_INVALID
	}
	item := itemAt(list, -1)
	if item.version != ver {
		return VERSION_INVALID
	}

	nver := this.NextVersion()
	for _, pi := range list {
		pi.version = nver
	}
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.invokeListener(ACTION_UPDATE, key, item)
	return nver
}

func (this *localMemGroup) Set(key MemKey, val interface{}, sz int) MemVer {
	var ctx localActionContext
	list, _, ver := this.Query(key, &ctx)
	if ver == VERSION_INVALID {
		ver = this.NextVersion()
	}
	for _, pi := range list {
		pi.version = ver
	}

	item := itemAt(list, -1)
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.tryInvokeListener(&ctx, ACTION_UPDATE, key, item)
	ctx.Invoke()
	return ver
}

func (this *localMemGroup) InitSet(ctx *localActionContext, key MemKey, val interface{}, sz int, ver MemVer) {
	list, _, _ := this.Query(key, ctx)
	for _, pi := range list {
		pi.version = ver
	}

	item := itemAt(list, -1)
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.tryInvokeListener(ctx, ACTION_UPDATE, key, item)
}

func (this *localMemGroup) doDelete(ctx *localActionContext, key MemKey, cur *localMemItem) {
	for k, item := range cur.items {
		nkey := append(key, k)

		delete(cur.items, k)
		this.tryInvokeListener(ctx, ACTION_DELETE, nkey, item)
		this.doDelete(ctx, nkey, item)
	}
	cur.items = nil
	this.size -= local_SIZE_ITEM_MAP

	cur.value = nil
	this.size -= int64(cur.size)

	this.size -= local_SIZE_ITEM
	this.count--
}

func (this *localMemGroup) Delete(key MemKey) MemVer {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return VERSION_INVALID
	}
	if len(list) == 1 {
		// don't delete root
		return VERSION_INVALID
	}
	item := itemAt(list, -1)
	p := itemAt(list, -2)
	skey, _ := key.At(-1)

	ver := this.NextVersion()
	for _, pi := range list {
		pi.version = ver
	}

	var ctx localActionContext
	delete(p.items, skey)
	this.tryInvokeListener(&ctx, ACTION_DELETE, key, item)
	this.doDelete(&ctx, key, item)
	ctx.Invoke()
	return ver
}

func (this *localMemGroup) AddListener(key MemKey, id string, lis XMemListener) bool {
	items, _ := this.LisQuery(key, true)
	item := items[len(items)-1]
	if item.listeners == nil {
		item.listeners = make(map[string]XMemListener)
	}
	item.listeners[id] = lis
	return true
}

func (this *localMemGroup) RemoveListener(key MemKey, id string) {
	items, ok := this.LisQuery(key, false)
	if !ok {
		return
	}
	item := items[len(items)-1]
	if item.listeners != nil {
		delete(item.listeners, id)
	}
	for i := len(items) - 1; i > 0; i-- {
		it := items[i]
		if len(it.items) == 0 && len(it.listeners) == 0 {
			it.listeners = nil
			p := items[i-1]
			k := key[i-1]
			delete(p.items, k)
		}
	}
}

func safeInvokeListener(lis XMemListener, elist []*XMemEvent) {
	defer func() {
		recover()
	}()
	if lis != nil {
		lis(elist)
	}
}

func (this *localMemGroup) tryInvokeListener(ctx *localActionContext, action Action, key MemKey, eitem *localMemItem) {
	list, _ := this.LisQuery(key, false)
	var ev *XMemEvent
	for _, item := range list {
		if len(item.listeners) > 0 {
			if ev == nil {
				ev = new(XMemEvent)
				ev.Action = action
				ev.Key = key
				ev.GroupName = this.name
				ev.Value = eitem.value
				ev.Version = eitem.version
			}
			ctx.Add(item, ev)
		}
	}
}

func (this *localMemGroup) invokeListener(action Action, key MemKey, eitem *localMemItem) {
	list, _ := this.LisQuery(key, false)
	var elist []*XMemEvent
	for _, item := range list {
		for _, lis := range item.listeners {
			if elist == nil {
				ev := new(XMemEvent)
				ev.Action = action
				ev.Key = key
				ev.GroupName = this.name
				ev.Value = eitem.value
				ev.Version = eitem.version
				elist = []*XMemEvent{ev}
			}
			safeInvokeListener(lis, elist)
		}
	}
}
