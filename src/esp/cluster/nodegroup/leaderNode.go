package nodegroup

import (
	"esp/cluster/clusterbase"
	"esp/cluster/nodeid"
	"logger"
)

func (this *NodeGroup) doStartLead(old nodeid.NodeId) error {
	this.doStopAll()

	logger.Debug(tag, "%s doStartLead(%d)", this, old)
	this.role = clusterbase.ROLE_LEADER

	if old != nodeid.INVALID {
		// fast learn from old leader
	}

	return nil
}
