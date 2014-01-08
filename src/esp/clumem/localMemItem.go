package clumem

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	local_SIZE_ITEM     = (8 + 8 + 8 + 8 + 8) * 2
	local_SIZE_ITEM_MAP = 24
)

type localMemItem struct {
	items     map[string]*localMemItem
	value     interface{}
	size      int
	version   MemVer
	listeners map[string]IMemListener
}

func (this *localMemItem) Len() int {
	return len(this.items)
}

func (this *localMemItem) ToString(n string, lvl int) string {
	return fmt.Sprintf("%s%s(%d:%d) : %v", strings.Repeat("\t", lvl), n, this.Len(), this.version, this.value)
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
