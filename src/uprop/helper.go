package uprop

import "fmt"

type helper int

var (
	Helper helper
)

func (O helper) Set(props []*UProperty, n string, v string) error {
	for _, p := range props {
		if p.Name == n {
			return p.CallSet(v)
		}
	}
	return fmt.Errorf("invalid prop '%s'", n)
}
