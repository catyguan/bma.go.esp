package xmemservice

import (
	"bytes"
	"esp/xmem/xmemprot"
	"fmt"
	"strings"
)

const (
	local_SIZE_ITEM     = (8 + 8 + 8 + 8 + 8) * 2
	local_SIZE_ITEM_MAP = 24
)

type localMemItem struct {
	items   map[string]*localMemItem
	value   interface{}
	size    int
	version xmemprot.MemVer
}

func (this *localMemItem) Clear() {
	this.value = nil
	for _, item := range this.items {
		item.Clear()
	}
}

func (this *localMemItem) Len() int {
	return len(this.items)
}

func (this *localMemItem) ToString(n string, lvl int) string {
	return fmt.Sprintf("%s%s(%d:%d) : %v", strings.Repeat("\t", lvl), n, this.Len(), this.version, this.value)
}

func (this *localMemItem) Walk(key xmemprot.MemKey, w XMemWalker) WalkStep {
	ws := w(key, this.value, this.version)
	if ws == WALK_STEP_NEXT {
		for k, item := range this.items {
			nkey := append(key, k)
			nws := item.Walk(nkey, w)
			switch nws {
			case WALK_STEP_END:
				return nws
			case WALK_STEP_OUT:
				return ws
			case WALK_STEP_OVER, WALK_STEP_NEXT:
				continue
			}
		}
	}
	return ws
}

func (this *localMemItem) Dump(n string, buf *bytes.Buffer, lvl int, all bool) {
	buf.WriteString(this.ToString(n, lvl))
	buf.WriteString("\n")
	if this.items != nil {
		for k, item := range this.items {
			if all {
				item.Dump(k, buf, lvl+1, all)
			} else {
				buf.WriteString(item.ToString(k, lvl+1))
				buf.WriteString("\n")
			}
		}
	}
}
