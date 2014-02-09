package nodegroup

import (
	"esp/cluster/clusterbase"
	"esp/cluster/nodeid"
)

func (this *NodeGroup) doSyncFrom(id nodeid.NodeId, ver clusterbase.OpVer) error {
	return nil
}

func (this *NodeGroup) doLearnFrom(id nodeid.NodeId, ver clusterbase.OpVer) error {

}
