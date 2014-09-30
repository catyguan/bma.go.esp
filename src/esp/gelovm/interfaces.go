package gelovm

import "code.google.com/p/gelo"

type VMInitor func(vm *gelo.VM) error

type VMIFactory func(cfg *VMSConfig) (VMInitor, error)

func SimpleVMIFactory(vmi VMInitor) VMIFactory {
	return func(cfg *VMSConfig) (VMInitor, error) {
		return vmi, nil
	}
}
