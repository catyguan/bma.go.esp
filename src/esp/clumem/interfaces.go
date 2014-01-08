package clumem

type MemKey []string
type MemVer uint32

type Action int

func (O Action) String() string {
	switch O {
	case ACTION_NONE:
		return "NONE"
	case ACTION_NEW:
		return "NEW"
	case ACTION_UPDATE:
		return "UDPATE"
	case ACTION_DELETE:
		return "DELETE"
	default:
		return "UNKNOW"
	}
}

const (
	ACTION_NONE = iota
	ACTION_NEW
	ACTION_UPDATE
	ACTION_DELETE
)

type IMemListener func(action Action, groupName string, key MemKey, val interface{}, ver MemVer)
