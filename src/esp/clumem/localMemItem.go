package clumem

import (
	"bytes"
	"fmt"
)

const (
	local_SIZE_ITEM     = (8 + 8 + 8 + 8) * 2
	local_SIZE_ITEM_MAP = 24
)

type localMemItem struct {
	items   map[string]*localMemItem
	value   interface{}
	size    int
	version MemVer
}

func (this *localMemItem) Dump(n string, buf *bytes.Buffer, lvl int) {
	for i := 0; i < lvl; i++ {
		buf.WriteString("\t")
	}
	c := 0
	if this.items != nil {
		c = len(this.items)
	}
	buf.WriteString(fmt.Sprintf("%s(%d:%d) : %v\n", n, c, this.version, this.value))
	if this.items != nil {
		for n, item := range this.items {
			item.Dump(n, buf, lvl+1)
		}
	}
}
