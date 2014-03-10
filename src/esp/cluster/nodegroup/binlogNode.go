package nodegroup

import (
	"esp/cluster/clusterbase"
	"esp/cluster/nodeinfo"
)

func (this *NodeGroup) doSyncFrom(id nodeinfo.NodeId, ver clusterbase.OpVer) error {
	return nil
}

func (this *NodeGroup) doLearnFrom(id nodeinfo.NodeId, ver clusterbase.OpVer) error {
	return nil
}
