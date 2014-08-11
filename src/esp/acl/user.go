package acl

import "fmt"

type key int

const KEY key = 0

type User struct {
	Id     string
	Name   string
	Groups []string
}

func (this *User) String() string {
	return fmt.Sprintf("user(%s)", this.Name)
}

func (this *User) IsWho(aclName string) bool {
	if aclName == "*" {
		return true
	}
	if this.Id == aclName {
		return true
	}
	for _, k := range this.Groups {
		if k == aclName {
			return true
		}
	}
	return false
}
