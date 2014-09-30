package gelovm

import (
	"boot"

	"code.google.com/p/gelo"
)

type VMSConfig struct {
	Modules  map[string]map[string]interface{}
	Paths    []string
	Preloads []string
}

func (this *VMSConfig) Valid() error {
	return nil
}

func (this *VMSConfig) Compare(old *VMSConfig) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	// compare Modules
	if true {
		same := func() bool {
			if len(this.Modules) != len(old.Modules) {
				return false
			}
			tmp := make(map[string]bool)
			for _, s := range this.Modules {
				tmp[s] = true
			}
			for _, s := range old.Modules {
				if _, ok := tmp[s]; !ok {
					return false
				}
			}
			return true
		}()
		if !same {
			return boot.CCR_NEED_START
		}
	}

	// compare Paths
	if true {
		same := func() bool {
			if len(this.Paths) != len(old.Paths) {
				return false
			}
			tmp := make(map[string]bool)
			for _, s := range this.Paths {
				tmp[s] = true
			}
			for _, s := range old.Paths {
				if _, ok := tmp[s]; !ok {
					return false
				}
			}
			return true
		}()
		if !same {
			return boot.CCR_NEED_START
		}
	}

	return boot.CCR_NONE
}

type VMS struct {
	rvm *gelo.VM
}

func NewVMS(io gelo.Port) *VMS {
	r := new(VMS)
	r.vms = gelo.NewVM(io)
	return r
}

func (this *VMS) Root() *gelo.VM {
	return this.rvm
}

func (this *VMS) Setup(vmif VMIFactory, cfg *VMSConfig) error {

}

func (this *VMS) Close() {
	this.rvm.Destroy()
}
