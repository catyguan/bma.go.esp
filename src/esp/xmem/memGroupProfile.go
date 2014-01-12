package xmem

import "fmt"

type memGroupProfile struct {
	Name  string
	Coder XMemCoder
}

func (this *memGroupProfile) Valid() error {
	if this.Name == "" {
		return fmt.Errorf("name empty")
	}
	return nil
}
