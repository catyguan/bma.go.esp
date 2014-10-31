package golua

func (this *GoLua) GetGlobal(n string) (interface{}, bool) {
	this.globalMutex.RLock()
	defer this.globalMutex.RUnlock()
	v, ok := this.globals[n]
	return v, ok
}

func (this *GoLua) SetGlobal(n string, v interface{}) interface{} {
	this.globalMutex.Lock()
	defer this.globalMutex.Unlock()
	old := this.globals[n]
	this.globals[n] = v
	return old
}
