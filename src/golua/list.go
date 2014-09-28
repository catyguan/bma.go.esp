package golua

type list struct {
	data []interface{}
	pos  int
}

func newList() *list {
	r := new(list)
	r.data = make([]interface{}, 0)
	r.pos = 0
	return r
}

func (this *list) clear() {
	this.data = make([]interface{}, 0)
	this.pos = 0
}

func (this *list) get(i int) interface{} {
	return this.data[i]
}

func (this *list) add(n interface{}) {
	if this.pos < len(this.data) {
		this.data[this.pos] = n
	} else {
		this.data = append(this.data, n)
	}
	this.pos++
}

func (this *list) remove(i int) interface{} {
	r := this.data[i]
	this.data[i] = nil
	copy(this.data[i:], this.data[i+1:])
	this.pos--
	return r
}

func (this *list) size() int {
	return this.pos
}
