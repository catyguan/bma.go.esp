package snapshot

import (
	"bytes"
	"esp/cluster/clusterbase"
	"fmt"
	"io/ioutil"
	"logger"
	"os"
	"sort"
)

const (
	tag = "binlog"
)

type SnapshotInfo struct {
	ver clusterbase.OpVer
}

type SortInfo []*SnapshotInfo

func (a SortInfo) Len() int           { return len(a) }
func (a SortInfo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortInfo) Less(i, j int) bool { return a[i].ver < a[j].ver }

// Service
type Service struct {
	name   string
	config *SnapshotConfig

	infos []*SnapshotInfo
}

func NewService(n string, cfg *SnapshotConfig) *Service {
	this := new(Service)
	this.name = n
	this.config = cfg
	return this
}

func (this *Service) String() string {
	return fmt.Sprintf("Snapshot[%s]", this.name)
}

func (this *Service) Setup() error {
	if true {
		err := os.MkdirAll(this.config.DataDir, 0664)
		if err != nil {
			return err
		}
	}
	if true {
		err := this.ReadInfos()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Service) ReadInfos() error {
	if true {
		fs, err := ioutil.ReadDir(this.config.DataDir)
		if err != nil {
			return err
		}
		this.infos = make([]*SnapshotInfo, 0)
		for _, f := range fs {
			if f.IsDir() {
				continue
			}
			var num int64
			n, err2 := fmt.Sscanf(f.Name(), this.config.FileFormatter, &num)
			if err2 != nil || n != 1 {
				continue
			}
			info := new(SnapshotInfo)
			info.ver = clusterbase.OpVer(num)
			this.infos = append(this.infos, info)
		}
		sort.Sort(SortInfo(this.infos))
	}

	if logger.EnableDebug(tag) {
		buf := bytes.NewBuffer([]byte{})
		for i, info := range this.infos {
			if i != 0 {
				buf.WriteString(",")
			}
			buf.WriteString(fmt.Sprintf("%d", info.ver))
		}
		logger.Debug(tag, "%s setupInfo [%s]", this, buf.String())
	}
	return nil
}

func (this *Service) Match(ver clusterbase.OpVer) clusterbase.OpVer {
	r := clusterbase.OpVer(0)
	for _, info := range this.infos {
		if info.ver > ver {
			break
		}
		r = info.ver
	}
	return r
}

func (this *Service) NewWriter(ver clusterbase.OpVer) (*Writer, error) {
	r := new(Writer)
	err := r.InitWriter(this, ver)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (this *Service) NewReader(ver clusterbase.OpVer, flag int) (*os.File, error) {
	fn := this.config.GetFileName(ver)
	logger.Debug(tag, "open snapshot reader(%s)", fn)
	if flag == 0 {
		flag = os.O_RDONLY
	}
	fd, err := os.OpenFile(fn, flag, 0664)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func (this *Service) Read(ver clusterbase.OpVer, pos, sz int64) ([]byte, error) {
	fd, err1 := this.NewReader(ver, 0)
	if err1 != nil {
		return nil, err1
	}
	b := make([]byte, sz)
	_, err2 := fd.ReadAt(b, pos)
	if err2 != nil {
		return nil, err2
	}
	return b, nil
}
