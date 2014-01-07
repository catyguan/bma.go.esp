package clumem

type MemKey []string
type MemVer uint32

type IService interface {
	Save() error
}
