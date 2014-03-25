package lua51

import "errors"

func (this *State) SetPath(path string) {
	this.GetGlobal("package") /* stack: package */
	this.PushString(path)     /* stack: package newpath */
	this.SetField(-2, "path") /* package.path = newpath, stack: package */
	this.Pop(1)               /* stack: - */
}

func (this *State) Eval(s string) error {
	if !this.DoString(s) {
		err := this.ToString(-1)
		this.Pop(1)
		return errors.New(err)
	}
	return nil
}

func (this *State) PushGValue(v interface{}) int {
	id := this.valueId + 1
	this.valueId = id
	this.values[id] = v
	this.NewTable()
	this.PushInteger(id)
	this.SetField(-2, "_gid")
	return id
}

func (this *State) ToGValue(index int) (interface{}, int, bool) {
	if !this.IsTable(index) {
		return nil, 0, false
	}
	this.GetField(index, "_gid")
	if this.IsNoneOrNil(-1) {
		this.Pop(1)
		return nil, 0, false
	}
	id := this.ToInteger(-1)
	this.Pop(1)
	r, ok := this.values[id]
	return r, id, ok
}

func (this *State) RemoveGValue(id int) {
	delete(this.values, id)
}

func (this *State) QueryGValue(gid int) (interface{}, bool) {
	r, ok := this.values[gid]
	return r, ok
}

func (this *State) ReplaceGValue(gid int, v interface{}) {
	this.values[gid] = v
}

func (this *State) ClearGValues() {
	for k, _ := range this.values {
		delete(this.values, k)
	}
}
