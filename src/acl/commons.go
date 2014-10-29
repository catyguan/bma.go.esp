package acl

import (
	"bmautil/valutil"
	"context"
	"fmt"
)

func CreateUsers(nu []string) []*User {
	r := make([]*User, len(nu))
	for idx, u := range nu {
		r[idx] = NewUserS(u)
	}
	return r
}

func CompareUsers(nu []string, ou []string) bool {
	if len(nu) != len(ou) {
		return false
	}
	tmp := make(map[string]bool)
	for _, n := range nu {
		tmp[n] = true
	}
	for _, n := range ou {
		if _, ok := tmp[n]; !ok {
			return false
		}
	}
	return true
}

type ConstRule struct {
	Users  []*User
	Result CHECK_RESULT
}

func (this *ConstRule) Check(user *User, path []string, ctx context.Context) (CHECK_RESULT, error) {
	for _, r := range this.Users {
		res := user.Match(r)
		if DebugMode {
			fmt.Println("checking", user, r, "->", this.Result, res)
		}
		if res {
			return this.Result, nil
		}
	}
	return UNKNOW, nil
}

func (this *ConstRule) String() string {
	return fmt.Sprintf("%s%v", this.Result, this.Users)
}

type ConstRuleFactory int

type configConstRule struct {
	Users []string
}

func (this ConstRuleFactory) Valid(cfg map[string]interface{}) error {
	var co configConstRule
	if valutil.ToBean(cfg, &co) {
		if len(co.Users) == 0 {
			return fmt.Errorf("Users empty")
		}
		return nil
	}
	return fmt.Errorf("invalid ConstRule config")
}

func (this ConstRuleFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo configConstRule
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if !CompareUsers(co.Users, oo.Users) {
		return false
	}
	return true
}

func (this ConstRuleFactory) Create(cfg map[string]interface{}) (Rule, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co configConstRule
	valutil.ToBean(cfg, &co)
	r := new(ConstRule)
	r.Users = CreateUsers(co.Users)
	r.Result = CHECK_RESULT(this)
	return r, nil
}

func init() {
	AddRuleFactory("pass", ConstRuleFactory(PASS))
	AddRuleFactory("deny", ConstRuleFactory(DENY))
}
