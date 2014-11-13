package golua

import (
	"bmautil/valutil"
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"unicode"
)

func StringsModule() *VMModule {
	m := NewVMModule("strings")
	m.Init("contains", GOF_strings_contains(0))
	m.Init("hasPrefix", GOF_strings_hasPrefix(0))
	m.Init("hasSuffix", GOF_strings_hasSuffix(0))
	m.Init("index", GOF_strings_index(0))
	m.Init("lastIndex", GOF_strings_lastIndex(0))
	m.Init("replace", GOF_strings_replace(0))
	m.Init("split", GOF_strings_split(0))
	m.Init("toLower", GOF_strings_toLower(0))
	m.Init("toUpper", GOF_strings_toUpper(0))
	m.Init("trim", GOF_strings_trim(0))
	m.Init("trimLeft", GOF_strings_trimLeft(0))
	m.Init("trimRight", GOF_strings_trimRight(0))
	m.Init("trimPrefix", GOF_strings_trimPrefix(0))
	m.Init("trimSuffix", GOF_strings_trimSuffix(0))
	m.Init("substr", GOF_strings_substr(0))
	m.Init("format", GOF_strings_format(0))
	m.Init("formata", GOF_strings_formata(0))
	m.Init("parsef", GOF_strings_parsef(0))
	m.Init("md5", GOF_strings_md5(0))
	return m
}

// strings.contains(s, substr string) bool
type GOF_strings_contains int

func (this GOF_strings_contains) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, substr, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vsubstr := valutil.ToString(substr, "")
	rv := strings.Contains(vs, vsubstr)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_contains) IsNative() bool {
	return true
}

func (this GOF_strings_contains) String() string {
	return "GoFunc<strings.contains>"
}

// strings.hasPrefix(s, prefix string) bool
type GOF_strings_hasPrefix int

func (this GOF_strings_hasPrefix) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, prefix, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vprefix := valutil.ToString(prefix, "")
	rv := strings.HasPrefix(vs, vprefix)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_hasPrefix) IsNative() bool {
	return true
}

func (this GOF_strings_hasPrefix) String() string {
	return "GoFunc<strings.hasPrefix>"
}

// strings.hasSuffix(s, ss string) bool
type GOF_strings_hasSuffix int

func (this GOF_strings_hasSuffix) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	rv := strings.HasSuffix(vs, vss)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_hasSuffix) IsNative() bool {
	return true
}

func (this GOF_strings_hasSuffix) String() string {
	return "GoFunc<strings.hasSuffix>"
}

// strings.index(s, sep string) int
type GOF_strings_index int

func (this GOF_strings_index) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	rv := strings.Index(vs, vss)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_index) IsNative() bool {
	return true
}

func (this GOF_strings_index) String() string {
	return "GoFunc<strings.index>"
}

// strings.lastIndex(s, sep string) int
type GOF_strings_lastIndex int

func (this GOF_strings_lastIndex) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	rv := strings.LastIndex(vs, vss)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_lastIndex) IsNative() bool {
	return true
}

func (this GOF_strings_lastIndex) String() string {
	return "GoFunc<strings.lastIndex>"
}

// strings.replace(s, old, new string[, n int]) string
type GOF_strings_replace int

func (this GOF_strings_replace) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(3)
	if err0 != nil {
		return 0, err0
	}
	s, olds, news, n, err1 := vm.API_pop4X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	volds := valutil.ToString(olds, "")
	vnews := valutil.ToString(news, "")
	vn := valutil.ToInt(n, -1)
	rv := strings.Replace(vs, volds, vnews, vn)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_replace) IsNative() bool {
	return true
}

func (this GOF_strings_replace) String() string {
	return "GoFunc<strings.replace>"
}

// strings.Split(s, sep string[,n int]) []string
type GOF_strings_split int

func (this GOF_strings_split) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, n, err1 := vm.API_pop3X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	vn := valutil.ToInt(n, -1)
	rv := strings.SplitN(vs, vss, vn)
	if rv != nil {
		arr := make([]interface{}, len(rv))
		for i, v := range rv {
			arr[i] = v
		}
		vm.API_push(vm.API_array(arr))
	} else {
		vm.API_push(nil)
	}
	return 1, nil
}

func (this GOF_strings_split) IsNative() bool {
	return true
}

func (this GOF_strings_split) String() string {
	return "GoFunc<strings.split>"
}

// strings.toLower(s string) string
type GOF_strings_toLower int

func (this GOF_strings_toLower) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := strings.ToLower(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_toLower) IsNative() bool {
	return true
}

func (this GOF_strings_toLower) String() string {
	return "GoFunc<strings.toLower>"
}

// strings.toLower(s string) string
type GOF_strings_toUpper int

func (this GOF_strings_toUpper) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	rv := strings.ToUpper(vs)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_toUpper) IsNative() bool {
	return true
}

func (this GOF_strings_toUpper) String() string {
	return "GoFunc<strings.toUpper>"
}

// strings.trim(s string[, cutset string]) string
type GOF_strings_trim int

func (this GOF_strings_trim) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	var rv string
	if vss == "" {
		rv = strings.TrimSpace(vs)
	} else {
		rv = strings.Trim(vs, vss)
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_trim) IsNative() bool {
	return true
}

func (this GOF_strings_trim) String() string {
	return "GoFunc<strings.trim>"
}

// strings.trimLeft(s string[, cutset string]) string
type GOF_strings_trimLeft int

func (this GOF_strings_trimLeft) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	var rv string
	if vss == "" {
		rv = strings.TrimLeftFunc(vs, unicode.IsSpace)
	} else {
		rv = strings.TrimLeft(vs, vss)
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_trimLeft) IsNative() bool {
	return true
}

func (this GOF_strings_trimLeft) String() string {
	return "GoFunc<strings.trimLeft>"
}

// strings.trimRight(s string[, cutset string]) string
type GOF_strings_trimRight int

func (this GOF_strings_trimRight) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	var rv string
	if vss == "" {
		rv = strings.TrimRightFunc(vs, unicode.IsSpace)
	} else {
		rv = strings.TrimRight(vs, vss)
	}
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_trimRight) IsNative() bool {
	return true
}

func (this GOF_strings_trimRight) String() string {
	return "GoFunc<strings.trimRight>"
}

// strings.trimPrefix(s string, cutset string) string
type GOF_strings_trimPrefix int

func (this GOF_strings_trimPrefix) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	rv := strings.TrimPrefix(vs, vss)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_trimPrefix) IsNative() bool {
	return true
}

func (this GOF_strings_trimPrefix) String() string {
	return "GoFunc<strings.trimPrefix>"
}

// strings.trimSuffix(s string, cutset string) string
type GOF_strings_trimSuffix int

func (this GOF_strings_trimSuffix) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, ss, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	vss := valutil.ToString(ss, "")
	rv := strings.TrimSuffix(vs, vss)
	vm.API_push(rv)
	return 1, nil
}

func (this GOF_strings_trimSuffix) IsNative() bool {
	return true
}

func (this GOF_strings_trimSuffix) String() string {
	return "GoFunc<strings.trimSuffix>"
}

// strings.substr(string $string , int $start [, int $length ] ) string
type GOF_strings_substr int

func (this GOF_strings_substr) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(2)
	if err0 != nil {
		return 0, err0
	}
	s, start, l, err1 := vm.API_pop3X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	clist := []rune(vs)
	vstart := valutil.ToInt(start, 0)
	vlen := valutil.ToInt(l, -1)
	if vstart < 0 || vstart >= len(clist) {
		return 0, fmt.Errorf("invalid start(%v) on string(%d)", vstart, len(clist))
	}
	if vlen < 0 {
		vlen = len(clist) - vstart
	}
	if vstart+vlen > len(clist) {
		return 0, fmt.Errorf("invalid len(%v, %v) on string(%d)", vstart, l, len(clist))
	}
	rv := clist[vstart : vstart+vlen]
	vm.API_push(string(rv))
	return 1, nil
}

func (this GOF_strings_substr) IsNative() bool {
	return true
}

func (this GOF_strings_substr) String() string {
	return "GoFunc<strings.substr>"
}

// strings.format(string $string , ... ) string
type GOF_strings_format int

func (this GOF_strings_format) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	top := vm.API_gettop()
	ps, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(ps[0], "")
	rs := fmt.Sprintf(vs, ps[1:]...)
	vm.API_push(rs)
	return 1, nil
}

func (this GOF_strings_format) IsNative() bool {
	return true
}

func (this GOF_strings_format) String() string {
	return "GoFunc<strings.format>"
}

// strings.formata(string , array ) string
type GOF_strings_formata int

func (this GOF_strings_formata) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, a, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	va := vm.API_toSlice(a)
	if va == nil {
		vm.API_push(vs)
	} else {
		rs := fmt.Sprintf(vs, va...)
		vm.API_push(rs)
	}
	return 1, nil
}

func (this GOF_strings_formata) IsNative() bool {
	return true
}

func (this GOF_strings_formata) String() string {
	return "GoFunc<strings.formata>"
}

// strings.parsef(string $string , ... ) string
type GOF_strings_parsef int

func (this GOF_strings_parsef) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	top := vm.API_gettop()
	ps, err1 := vm.API_popN(top, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(ps[0], "")
	rs, err2 := this.parsef(vs, ps[1:])
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(rs)
	return 1, nil
}

func (this GOF_strings_parsef) IsNative() bool {
	return true
}

func (this GOF_strings_parsef) String() string {
	return "GoFunc<strings.parsef>"
}

func (this GOF_strings_parsef) parsef(str string, vlist []interface{}) (string, error) {
	out := bytes.NewBuffer(make([]byte, 0))
	var c1 rune = 0
	word := bytes.NewBuffer(make([]byte, 0))

	idx := 0
	nvlist := make([]interface{}, 0, len(vlist))
	for _, c := range []rune(str) {
		switch c1 {
		case 0:
			if c == '$' {
				c1 = c
			} else {
				out.WriteRune(c)
			}
		case '$':
			if c == '{' {
				c1 = '{'
				idx = 0
			} else {
				out.WriteRune(c1)
				out.WriteRune(c)
				c1 = 0
			}
		case '{':
			if c == '}' {
				varname := word.String()
				word.Reset()

				var nv interface{}
				x := idx - 1
				if x >= 0 && x < len(vlist) {
					nv = vlist[x]
					out.WriteByte('%')
					out.WriteString(varname)
					nvlist = append(nvlist, nv)
				}
				c1 = 0
			} else {
				if idx == 0 {
					idx = int(c) - '0'
				} else {
					word.WriteRune(c)
				}
			}
		}
	}

	if word.Len() > 0 {
		return "", fmt.Errorf("invalid parse format(%s)", str)
	}

	rs := fmt.Sprintf(out.String(), nvlist...)
	return rs, nil
}

// strings.md5(string $string) string
type GOF_strings_md5 int

func (this GOF_strings_md5) Exec(vm *VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	s, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vs := valutil.ToString(s, "")
	h := md5.New()
	h.Write([]byte(vs))
	r := fmt.Sprintf("%x", h.Sum(nil))
	vm.API_push(r)
	return 1, nil
}

func (this GOF_strings_md5) IsNative() bool {
	return true
}

func (this GOF_strings_md5) String() string {
	return "GoFunc<strings.md5>"
}
