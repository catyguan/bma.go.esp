package snapshot

import (
	"esp/cluster/clusterbase"
	"fmt"
	"path/filepath"
)

type SnapshotConfig struct {
	DataDir       string
	FileFormatter string
}

func (this *SnapshotConfig) Valid() error {
	if this.DataDir == "" {
		return fmt.Errorf("snapshot data dir name empty")
	}
	if this.FileFormatter == "" {
		this.FileFormatter = "snapshot.%d.blog"
	}
	return nil
}

func (this *SnapshotConfig) GetFileName(ver clusterbase.OpVer) string {
	return filepath.Join(this.DataDir, fmt.Sprintf(this.FileFormatter, ver))
}
