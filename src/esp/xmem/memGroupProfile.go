package xmem

import "fmt"

type MemGroupProfile struct {
	Name  string
	Coder XMemCoder
}

func (this *MemGroupProfile) Valid() error {
	if this.Name == "" {
		return fmt.Errorf("name empty")
	}
	return nil
}
