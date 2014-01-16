package xmemservice

import "esp/xmem/xmemprot"

type XMem4Service struct {
	service *Service
	name    string
}

func (this *XMem4Service) Init(s *Service, n string) {
	this.service = s
	this.name = n
}

func (this *XMem4Service) Get(key xmemprot.MemKey) (interface{}, xmemprot.MemVer, bool, error) {
	return this.GetAndListen(key, "", nil)
}

func (this *XMem4Service) GetAndListen(key xmemprot.MemKey, id string, lis XMemListener) (interface{}, xmemprot.MemVer, bool, error) {
	var rval interface{}
	rver := xmemprot.VERSION_INVALID
	rb := false
	err := this.service.executor.DoSync("xmemGet", func() error {
		si, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		if lis != nil {
			defer si.group.AddListener(key, id, lis)
		}
		mi, b := si.group.Get(key)
		if !b {
			return nil
		}
		rval = mi.value
		rver = mi.version
		rb = true
		return nil
	})
	return rval, rver, rb, err
}

func (this *XMem4Service) List(key xmemprot.MemKey) ([]string, bool, error) {
	return this.ListAndListen(key, "", nil)
}

func (this *XMem4Service) ListAndListen(key xmemprot.MemKey, id string, lis XMemListener) ([]string, bool, error) {
	var rlist []string
	rb := false
	err := this.service.executor.DoSync("xmemList", func() error {
		si, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		if lis != nil {
			defer si.group.AddListener(key, id, lis)
		}
		mi, b := si.group.Get(key)
		if !b {
			return nil
		}
		list := make([]string, 0, len(mi.items))
		for k, _ := range mi.items {
			list = append(list, k)
		}
		rlist = list
		rb = true
		return nil
	})
	return rlist, rb, err
}

func (this *XMem4Service) AddListener(key xmemprot.MemKey, id string, lis XMemListener) error {
	err := this.service.executor.DoSync("xmemAddListener", func() error {
		si, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		si.group.AddListener(key, id, lis)
		return nil
	})
	return err
}

func (this *XMem4Service) RemoveListener(key xmemprot.MemKey, id string) error {
	err := this.service.executor.DoSync("xmemRemoveListener", func() error {
		si, err := this.service.doGetGroup(this.name)
		if err != nil {
			return err
		}
		si.group.RemoveListener(key, id)
		return nil
	})
	return err
}

func (this *XMem4Service) Set(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error) {
	rver := xmemprot.VERSION_INVALID
	err := this.service.executor.DoSync("xmemSet", func() error {
		b, err := this.service.doSetOp(this.name, key, val, sz, xmemprot.VERSION_INVALID, false)
		if err != nil {
			return err
		}
		rver = b
		return nil
	})
	return rver, err
}

func (this *XMem4Service) CompareAndSet(key xmemprot.MemKey, val interface{}, sz int, ver xmemprot.MemVer) (xmemprot.MemVer, error) {
	rver := xmemprot.VERSION_INVALID
	err := this.service.executor.DoSync("xmemCompareAndSet", func() error {
		b, err := this.service.doSetOp(this.name, key, val, sz, ver, false)
		if err != nil {
			return err
		}
		rver = b
		return nil
	})
	return rver, err
}

func (this *XMem4Service) SetIfAbsent(key xmemprot.MemKey, val interface{}, sz int) (xmemprot.MemVer, error) {
	rver := xmemprot.VERSION_INVALID
	err := this.service.executor.DoSync("xmemSetIfAbsent", func() error {
		b, err := this.service.doSetOp(this.name, key, val, sz, xmemprot.VERSION_INVALID, true)
		if err != nil {
			return err
		}
		rver = b
		return nil
	})
	return rver, err
}

func (this *XMem4Service) Delete(key xmemprot.MemKey) (bool, error) {
	return this.CompareAndDelete(key, xmemprot.VERSION_INVALID)
}

func (this *XMem4Service) CompareAndDelete(key xmemprot.MemKey, ver xmemprot.MemVer) (bool, error) {
	rb := false
	err := this.service.executor.DoSync("xmemDelete", func() error {
		b, err := this.service.doDeleteOp(this.name, key, ver)
		if err != nil {
			return err
		}
		rb = b
		return nil
	})
	return rb, err
}
