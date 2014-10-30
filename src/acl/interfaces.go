package acl

import (
	"bmautil/valutil"
	"bytes"
	"context"
	"fmt"
	"logger"
	"strings"
)

const (
	tag = "acl"
)

var (
	DebugMode = false
)

type User struct {
	Account string
	Domain  string
	Groups  []string
	Data    interface{}
}

func NewUser(a, d string, g []string) *User {
	r := new(User)
	r.Account = a
	r.Domain = d
	r.Groups = g
	return r
}

func NewUserS(s string) *User {
	p1 := strings.SplitN(s, ":", 2)
	p2 := strings.SplitN(p1[0], "@", 2)
	var a, d string
	var g []string
	a = p2[0]
	if len(p2) > 1 {
		d = p2[1]
	}
	if len(p1) > 1 {
		g = strings.Split(p1[1], ",")
	}
	return NewUser(a, d, g)
}

func (this *User) String() string {
	return fmt.Sprintf("%s@%s", this.Account, this.Domain)
}

func (this *User) Match(rule *User) bool {
	if rule == nil {
		return false
	}
	if rule.Account != "*" && rule.Account != this.Account {
		return false
	}
	if rule.Domain != "*" && rule.Domain != this.Domain {
		return false
	}
	return true
}

type CHECK_RESULT int

func (this CHECK_RESULT) String() string {
	switch this {
	case PASS:
		return "pass"
	case DENY:
		return "deny"
	default:
		return "unknow"
	}
}

const (
	UNKNOW = CHECK_RESULT(0)
	PASS   = CHECK_RESULT(1)
	DENY   = CHECK_RESULT(2)
)

type Rule interface {
	Check(user *User, path []string, ctx context.Context) (CHECK_RESULT, Rule, error)
}

type RuleFactory interface {
	Valid(cfg map[string]interface{}) error
	Compare(cfg map[string]interface{}, old map[string]interface{}) (same bool)
	Create(cfg map[string]interface{}) (Rule, error)
}

var (
	rflibs map[string]RuleFactory = make(map[string]RuleFactory)
)

func AddRuleFactory(n string, fac RuleFactory) {
	rflibs[n] = fac
}

func GetRuleFactory(n string) RuleFactory {
	return rflibs[n]
}

func GetRuleFactoryByType(cfg map[string]interface{}) (RuleFactory, string, error) {
	xt, ok := cfg["Type"]
	if !ok {
		return nil, "", fmt.Errorf("miss Type")
	}
	vxt := valutil.ToString(xt, "")
	if vxt == "" {
		return nil, "", fmt.Errorf("invalid Type(%v)", xt)
	}
	fac := GetRuleFactory(vxt)
	if fac == nil {
		return nil, "", fmt.Errorf("invalid Rule Type(%s)", xt)
	}
	return fac, vxt, nil
}

func CreateRule(cfg map[string]interface{}) (Rule, error) {
	fac, _, err := GetRuleFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	return fac.Create(cfg)
}

func ValidConfig(cfg map[string]interface{}) error {
	fac, _, err := GetRuleFactoryByType(cfg)
	if err != nil {
		return err
	}
	return fac.Valid(cfg)
}

func CompareConfig(cfg map[string]interface{}, old map[string]interface{}) bool {
	fac1, xt1, err1 := GetRuleFactoryByType(cfg)
	if err1 != nil {
		return false
	}
	_, xt2, err2 := GetRuleFactoryByType(old)
	if err2 != nil {
		return false
	}
	if xt1 != xt2 {
		return false
	}
	return fac1.Compare(cfg, old)
}

func InitRuleTree(rt *RuleTree) {
	ruleTree = rt
}

var (
	ruleTree *RuleTree
)

type RuleTree struct {
	nodes map[string]*RuleTree
	rules []Rule
}

func NewRuleTree() *RuleTree {
	r := new(RuleTree)
	return r
}

func (this *RuleTree) Node(n string) *RuleTree {
	if n == "" {
		return this
	}
	if this.nodes == nil {
		this.nodes = make(map[string]*RuleTree)
	}
	if node, ok := this.nodes[n]; ok {
		return node
	}
	node := new(RuleTree)
	this.nodes[n] = node
	return node
}

func (this *RuleTree) Append(rule Rule) {
	if this.rules == nil {
		this.rules = make([]Rule, 0)
	}
	this.rules = append(this.rules, rule)
}

func (this *RuleTree) Dump() string {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	this.dump(buf, "")
	return buf.String()
}

func (this *RuleTree) dump(buf *bytes.Buffer, prex string) {
	for _, r := range this.rules {
		buf.WriteString(prex)
		buf.WriteString("- ")
		buf.WriteString(fmt.Sprintf("%v", r))
		buf.WriteString("\n")
	}
	for k, ns := range this.nodes {
		buf.WriteString(prex)
		buf.WriteString(fmt.Sprintf("<%s>\n", k))
		ns.dump(buf, prex+" ")
	}
}

func (this *RuleTree) DoCheck(user *User, path []string, ctx context.Context) (CHECK_RESULT, Rule, error) {
	if this.nodes != nil && len(path) > 0 {
		n := path[0]
		child := this.Node(n)
		if child != nil {
			cr, rule, err := child.DoCheck(user, path[1:], ctx)
			if err != nil {
				return UNKNOW, rule, err
			}
			if cr != UNKNOW {
				return cr, rule, nil
			}
		}
	}
	for _, r := range this.rules {
		cr, rule, err := r.Check(user, path, ctx)
		if err != nil {
			return UNKNOW, rule, err
		}
		if cr != UNKNOW {
			return cr, rule, nil
		}
	}
	return UNKNOW, nil, nil
}

func Check(user *User, path []string, ctx context.Context, def bool) (bool, Rule, error) {
	rt := ruleTree
	if rt != nil {
		cr, rule, err := rt.DoCheck(user, path, ctx)
		if err != nil {
			return false, rule, err
		}
		switch cr {
		case PASS:
			return true, rule, nil
		case DENY:
			return false, rule, nil
		}
	}
	return def, nil, nil
}

type AclError struct {
	ErrorString string
}

func (err *AclError) Error() string { return err.ErrorString }

func Assert(user *User, path []string, ctx context.Context) error {
	ok, rule, err := Check(user, path, ctx, false)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	logger.Info(tag, "user(%s) assert(%s, %v) fail", user, strings.Join(path, "/"), rule)
	r := new(AclError)
	r.ErrorString = fmt.Sprintf("'%s' access '%s'!", user, strings.Join(path, "/"))
	return r
}

func DumpRuleTree() string {
	rt := ruleTree
	if rt != nil {
		return rt.Dump()
	}
	return "<empty>"
}
