package espnet

import "bytes"

type Address []string

func NewAddress(s string) Address {
	return Address([]string{s})
}

func NewAddressN(s ...string) Address {
	return Address(s)
}

func (this Address) Size() int {
	return len(this)
}

func (this Address) Identity() string {
	if len(this) > 0 {
		return this[0]
	}
	return ""
}

func (this Address) ListIdentity() []string {
	c := len(this)
	r := make([]string, c)
	copy(r, this)
	return r
}

func (this Address) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("Address[")
	for i, s := range this {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s)
	}
	buf.WriteString("]")
	return buf.String()
}

func (this Address) Add(a string) Address {
	r := append(this, a)
	return Address(r)
}

func (this Address) AddUnique(a string) Address {
	for _, s := range this {
		if a == s {
			return this
		}
	}
	return this.Add(a)
}

func (this Address) AddAll(a []string) Address {
	r := make(Address, len(this)+len(a))
	copy(r, this)
	copy(r[len(this):], a)
	return r
}

func (this Address) AddAllUnique(a []string) Address {
	m := make(map[string]bool)
	for _, s := range this {
		m[s] = true
	}
	for _, s := range a {
		m[s] = true
	}
	r := make(Address, 0, len(m))
	for s, _ := range m {
		r = append(r, s)
	}
	return r
}

func (this Address) Append(a Address) Address {
	return this.AddAllUnique([]string(a))
}
