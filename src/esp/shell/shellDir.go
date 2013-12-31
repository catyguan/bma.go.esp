package shell

import (
	"bytes"
)

type dirInfo struct {
	pwd   ShellDir
	stack []ShellDir
}

func (this dirInfo) IsRootNow() bool {
	return len(this.stack) == 0
}

func newDirInfo(root ShellDir) *dirInfo {
	r := new(dirInfo)
	r.pwd = root
	r.stack = make([]ShellDir, 0)
	return r
}

func (this *dirInfo) CurName() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, p := range this.stack {
		buf.WriteString(p.Name())
		buf.WriteByte('/')
	}
	if this.pwd != nil {
		buf.WriteString(this.pwd.Name())
		buf.WriteByte('/')
	}
	return buf.String()
}

type ShellDir interface {
	Name() string

	GetDir(name string) ShellDir

	GetCommand(name string) ShellProcessor

	List() []string
}

type ShellDirSupported interface {
	CreateShell() ShellDir
}

type DirInfoSupport interface {
	DirInfo() string
}

// ShellDirBase
type ShellDirBase struct {
	DirName    string
	Dirs       func() map[string]ShellDir
	Commands   func() map[string]ShellProcessor
	Infos      func() []string
	GetDirInfo func() string
}

func (this *ShellDirBase) CF(f ShellCommandHandler) ShellProcessor {
	return &shellProcessorFunc{f}
}

func (this *ShellDirBase) Name() string {
	return this.DirName
}

func (this *ShellDirBase) GetDir(name string) ShellDir {
	if this.Dirs == nil {
		return nil
	}
	for n, r := range this.Dirs() {
		if n == name {
			return r
		}
	}
	return nil
}

func (this *ShellDirBase) GetCommand(name string) ShellProcessor {
	if this.Commands == nil {
		return nil
	}
	for n, r := range this.Commands() {
		if n == name {
			return r
		}
	}
	return nil
}

func (this *ShellDirBase) List() []string {
	r := make([]string, 0)
	if this.Dirs != nil {
		for n, _ := range this.Dirs() {
			r = append(r, DirName(n))
		}
	}
	if this.Commands != nil {
		for n, _ := range this.Commands() {
			r = append(r, n)
		}
	}
	if this.Infos != nil {
		for _, n := range this.Infos() {
			r = append(r, n)
		}
	}
	return r
}

func (this *ShellDirBase) DirInfo() string {
	if this.GetDirInfo != nil {
		return this.GetDirInfo()
	}
	return ""
}

// ShellDirCommon
type ShellDirCommon struct {
	name        string
	Dirs        []ShellDir
	Commands    []ShellProcessor
	DirInfoFunc func() string
}

func NewShellDirCommon(name string) *ShellDirCommon {
	r := new(ShellDirCommon)
	r.name = name
	r.Dirs = make([]ShellDir, 0)
	r.Commands = make([]ShellProcessor, 0)
	return r
}

func (this *ShellDirCommon) Name() string {
	return this.name
}

func (this *ShellDirCommon) GetDir(name string) ShellDir {
	for _, r := range this.Dirs {
		if r.Name() == name {
			return r
		}
	}
	return nil
}

func (this *ShellDirCommon) AddDir(o ShellDir) {
	this.Dirs = append(this.Dirs, o)
}

func (this *ShellDirCommon) GetCommand(name string) ShellProcessor {
	for _, r := range this.Commands {
		if o, ok := r.(NameSupport); ok {
			if o.Name() == name {
				return r
			}
		}
	}
	return nil
}

func (this *ShellDirCommon) AddCommand(o ShellProcessor) {
	this.Commands = append(this.Commands, o)
}

func (this *ShellDirCommon) List() []string {
	r := make([]string, 0, len(this.Dirs)+len(this.Commands))
	for _, o := range this.Dirs {
		r = append(r, DirName(o.Name()))
	}
	for _, o := range this.Commands {
		if n, ok := o.(NameSupport); ok {
			r = append(r, n.Name())
		}
	}
	return r
}

func (this *ShellDirCommon) DirInfo() string {
	if this.DirInfoFunc != nil {
		return this.DirInfoFunc()
	}
	return ""
}

func DirName(name string) string {
	return "<" + name + ">"
}

func InfoName(name string) string {
	return "[" + name + "]"
}
