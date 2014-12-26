package gotask

type GoTask struct {
	C chan bool
}

func NewGoTask() *GoTask {
	r := new(GoTask)
	r.Init()
	return r
}

func (this *GoTask) Init() {
	this.C = make(chan bool, 1)
}

func (this *GoTask) Close() {
	defer func() {
		recover()
	}()
	close(this.C)
}

func (this *GoTask) IsClose() bool {
	select {
	case <-this.C:
		return true
	default:
		return false
	}
}
