package nodegroup

import (
	"esp/cluster/clusterbase"
	"esp/cluster/nodeinfo"
	"logger"
)

func (this *NodeGroup) doStartLead(old nodeinfo.NodeId) error {
	this.doStopAll()

	logger.Debug(tag, "%s doStartLead(%d)", this, old)
	this.role = clusterbase.ROLE_LEADER

	if old != nodeinfo.INVALID {
		// fast learn from old leader
	}

	return nil
}
