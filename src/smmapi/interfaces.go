package smmapi

type SMA_TYPE int

const (
	SMA_API    = SMA_TYPE(0)
	SMA_HTTPUI = SMA_TYPE(1)
)

type SMAction struct {
	Id    string
	Title string
	Type  SMA_TYPE
	UIN   string
}

type SMInfo struct {
	Id      string
	Title   string
	Content string
	Actions []*SMAction
}

type smlist []*SMInfo

func (this smlist) Len() int {
	return len(this)
}

func (this smlist) Less(i, j int) bool {
	id1 := this[i].Id
	id2 := this[j].Id
	return id1 < id2
}

func (this smlist) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type SMMObject interface {
	GetInfo() (*SMInfo, error)
	// Result, error
	ExecuteAction(aid string, param map[string]interface{}) (interface{}, error)
}
