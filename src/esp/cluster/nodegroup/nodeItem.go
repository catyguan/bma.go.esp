package nodegroup

import (
	"esp/cluster/nodeinfo"
	"fmt"
)

type NodeItem struct {
	NodeId   nodeinfo.NodeId
	NodeName string
}

func (this *NodeItem) String() string {
	return fmt.Sprintln("%d(%s)", this.NodeId, this.NodeName)
}
