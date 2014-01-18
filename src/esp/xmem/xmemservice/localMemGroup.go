package xmemservice

import (
	"bmautil/binlog"
	"bmautil/valutil"
	"bytes"
	"esp/xmem/xmemprot"
	"fmt"
)

// localMemLis
type localMemLis struct {
	items     map[string]*localMemLis
	listeners map[string]xmemprot.XMemListener
}

func (this *localMemLis) allInvokeListener(elist []*xmemprot.XMemEvent) {
	for _, lis := range this.listeners {
		safeInvokeListener(lis, elist)
	}
	for _, item := range this.items {
		item.allInvokeListener(elist)
	}
}

// localActionContext
type localActionContext struct {
	listeners map[*localMemLis][]*xmemprot.XMemEvent
}

func (this *localActionContext) Add(item *localMemLis, ev *xmemprot.XMemEvent) {
	if this.listeners == nil {
		this.listeners = make(map[*localMemLis][]*xmemprot.XMemEvent)
	}
	elist, ok := this.listeners[item]
	if !ok {
		elist = []*xmemprot.XMemEvent{ev}
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
	version xmemprot.MemVer

	// for binlog & master/slave sync
	blver     binlog.BinlogVer
	blservice *binlog.Service
	blwriter  *binlog.Writer
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
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(this.name)
	buf.WriteString("[")
	buf.WriteString(fmt.Sprintf("count=%v; ", this.count))
	buf.WriteString(fmt.Sprintf("ver=%v; ", this.version))
	buf.WriteString(fmt.Sprintf("blver=%v; ", this.blver))
	buf.WriteString(fmt.Sprintf("size=%v; ", valutil.MakeSizeString(uint64(this.size))))
	buf.WriteString("]")
	return buf.String()
}

func (this *localMemGroup) Dump() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(this.String())
	buf.WriteString("\n")
	this.root.Dump("<root>", buf, 0, true)
	return buf.String()
}

func (this *localMemGroup) Walk(key xmemprot.MemKey, w XMemWalker) bool {
	item, ok := this.Get(key)
	if ok {
		item.Walk(key, w)
		return true
	}
	return false
}

func (this *localMemGroup) Snapshot(coder XMemCoder) (*XMemGroupSnapshot, error) {
	var rerr error
	r := make([]*XMemSnapshot, 0, this.count)
	this.root.Walk(xmemprot.MemKey{}, func(key xmemprot.MemKey, val interface{}, ver xmemprot.MemVer) WalkStep {
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

	obj := new(XMemGroupSnapshot)
	obj.BLVer = this.blver
	obj.Snapshots = r
	return obj, rerr
}

func (this *localMemGroup) Clear() {
	this.root.Clear()
	this.root = new(localMemItem)
	this.count = 1
	this.size = local_SIZE_ITEM
}

func (this *localMemGroup) BuildFromSnapshot(coder XMemCoder, gss *XMemGroupSnapshot) error {
	this.Clear()

	ev := new(xmemprot.XMemEvent)
	ev.Action = xmemprot.ACTION_CLEAR
	ev.GroupName = this.name
	ev.Key = xmemprot.MemKey{}
	ev.Value = nil
	ev.Version = xmemprot.MemVer(0)
	this.lisRoot.allInvokeListener([]*xmemprot.XMemEvent{ev})

	if gss.BLVer >= 0 {
		this.blver = gss.BLVer
	}
	slist := gss.Snapshots
	var ctx localActionContext
	for _, ss := range slist {
		val, sz, err := coder.Decode(ss.Kind, ss.Data)
		if err != nil {
			return err
		}
		this.InitSet(&ctx, xmemprot.MemKeyFromString(ss.Key), val, sz, ss.Version)
	}
	this.version = this.root.version

	ctx.Invoke()

	return nil
}

func (this *localMemGroup) NextVersion() xmemprot.MemVer {
	this.version++
	if this.version == xmemprot.VERSION_INVALID {
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

func (this *localMemGroup) Query(key xmemprot.MemKey, ctx *localActionContext) ([]*localMemItem, bool, xmemprot.MemVer) {
	r := []*localMemItem{this.root}
	cur := this.root
	ver := xmemprot.VERSION_INVALID
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
			if ver == xmemprot.VERSION_INVALID {
				ver = this.NextVersion()
			}
			item = new(localMemItem)
			item.version = ver

			cur.items[k] = item
			this.size += local_SIZE_ITEM
			this.count++

			this.tryInvokeListener(ctx, xmemprot.ACTION_NEW, key[:i+1], item)
		}
		r = append(r, item)
		cur = item
	}
	return r, true, ver
}

func (this *localMemGroup) LisQuery(key xmemprot.MemKey, createIfNotExists bool) ([]*localMemLis, bool) {
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

func (this *localMemGroup) Get(key xmemprot.MemKey) (*localMemItem, bool) {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return nil, false
	}
	return itemAt(list, -1), true
}

func (this *localMemGroup) CompareAndSet(key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer) xmemprot.MemVer {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return xmemprot.VERSION_INVALID
	}
	item := itemAt(list, -1)
	if item.version != ver {
		return xmemprot.VERSION_INVALID
	}

	nver := this.NextVersion()
	for _, pi := range list {
		pi.version = nver
	}
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.invokeListener(xmemprot.ACTION_UPDATE, key, item)
	return nver
}

func (this *localMemGroup) Set(key xmemprot.MemKey, val interface{}, sz int) xmemprot.MemVer {
	var ctx localActionContext
	list, _, ver := this.Query(key, &ctx)
	if ver == xmemprot.VERSION_INVALID {
		ver = this.NextVersion()
	}
	for _, pi := range list {
		pi.version = ver
	}

	item := itemAt(list, -1)
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.tryInvokeListener(&ctx, xmemprot.ACTION_UPDATE, key, item)
	ctx.Invoke()
	return ver
}

func (this *localMemGroup) SetIfAbsent(key xmemprot.MemKey, val interface{}, sz int) xmemprot.MemVer {
	_, b, _ := this.Query(key, nil)
	if b {
		return xmemprot.VERSION_INVALID
	}
	return this.Set(key, val, sz)
}

func (this *localMemGroup) InitSet(ctx *localActionContext, key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer) {
	list, _, _ := this.Query(key, ctx)
	for _, pi := range list {
		pi.version = ver
	}

	item := itemAt(list, -1)
	item.value = val
	item.size = sz
	this.size += int64(sz)

	this.tryInvokeListener(ctx, xmemprot.ACTION_UPDATE, key, item)
}

func (this *localMemGroup) doDelete(ctx *localActionContext, key xmemprot.MemKey, cur *localMemItem) {
	for k, item := range cur.items {
		nkey := append(key, k)

		delete(cur.items, k)
		this.tryInvokeListener(ctx, xmemprot.ACTION_DELETE, nkey, item)
		this.doDelete(ctx, nkey, item)
	}
	cur.items = nil
	this.size -= local_SIZE_ITEM_MAP

	cur.value = nil
	this.size -= int64(cur.size)

	this.size -= local_SIZE_ITEM
	this.count--
}

func (this *localMemGroup) Delete(key xmemprot.MemKey) xmemprot.MemVer {
	return this.CompareAndDelete(key, xmemprot.VERSION_INVALID)
}

func (this *localMemGroup) CompareAndDelete(key xmemprot.MemKey, cver xmemprot.MemVer) xmemprot.MemVer {
	list, ok, _ := this.Query(key, nil)
	if !ok {
		return xmemprot.VERSION_INVALID
	}
	if len(list) == 1 {
		// don't delete root
		return xmemprot.VERSION_INVALID
	}
	item := itemAt(list, -1)
	if cver != xmemprot.VERSION_INVALID {
		if item.version != cver {
			return xmemprot.VERSION_INVALID
		}
	}
	p := itemAt(list, -2)
	skey, _ := key.At(-1)

	ver := this.NextVersion()
	for _, pi := range list {
		pi.version = ver
	}

	var ctx localActionContext
	delete(p.items, skey)
	this.tryInvokeListener(&ctx, xmemprot.ACTION_DELETE, key, item)
	this.doDelete(&ctx, key, item)
	ctx.Invoke()
	return ver
}

func (this *localMemGroup) AddListener(key xmemprot.MemKey, id string, lis xmemprot.XMemListener) bool {
	items, _ := this.LisQuery(key, true)
	item := items[len(items)-1]
	if item.listeners == nil {
		item.listeners = make(map[string]xmemprot.XMemListener)
	}
	item.listeners[id] = lis
	return true
}

func (this *localMemGroup) RemoveListener(key xmemprot.MemKey, id string) {
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

func safeInvokeListener(lis xmemprot.XMemListener, elist []*xmemprot.XMemEvent) {
	defer func() {
		recover()
	}()
	if lis != nil {
		lis(elist)
	}
}

func (this *localMemGroup) tryInvokeListener(ctx *localActionContext, action xmemprot.Action, key xmemprot.MemKey, eitem *localMemItem) {
	list, _ := this.LisQuery(key, false)
	var ev *xmemprot.XMemEvent
	for _, item := range list {
		if len(item.listeners) > 0 {
			if ev == nil {
				ev = new(xmemprot.XMemEvent)
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

func (this *localMemGroup) invokeListener(action xmemprot.Action, key xmemprot.MemKey, eitem *localMemItem) {
	list, _ := this.LisQuery(key, false)
	var elist []*xmemprot.XMemEvent
	for _, item := range list {
		for _, lis := range item.listeners {
			if elist == nil {
				ev := new(xmemprot.XMemEvent)
				ev.Action = action
				ev.Key = key
				ev.GroupName = this.name
				ev.Value = eitem.value
				ev.Version = eitem.version
				elist = []*xmemprot.XMemEvent{ev}
			}
			safeInvokeListener(lis, elist)
		}
	}
}
