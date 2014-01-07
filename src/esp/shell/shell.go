package shell

import (
	"config"
	"fmt"
	"logger"
)

const (
	tag = "shell"
)

type SessionCloseListener func()

type SessionWriter interface {
	Write(msg string) bool

	Close()

	AddCloseListner(func())
}

type ConsoleWriter struct {
	closeListener func()
}

func NewConsoleWriter() *ConsoleWriter {
	r := new(ConsoleWriter)
	return r
}

func (this *ConsoleWriter) Write(msg string) bool {
	fmt.Print(msg)
	return true
}

func (this *ConsoleWriter) Close() {
	defer func() {
		recover()
	}()
	if this.closeListener != nil {
		this.closeListener()
	}
}

func (this *ConsoleWriter) AddCloseListner(lis func()) {
	this.closeListener = lis
}

type Session struct {
	Id        string
	Vars      map[string]interface{}
	defVars   map[string]func() interface{}
	writer    SessionWriter
	listeners []SessionCloseListener
	Execute   func(command string) bool
}

func NewSession(w SessionWriter) *Session {
	r := new(Session)
	r.Vars = make(map[string]interface{})
	r.listeners = make([]SessionCloseListener, 0)
	r.writer = w
	r.writer.AddCloseListner(func() {
		r.notifyAll()
	})
	r.defVars = make(map[string]func() interface{})
	return r
}

func (this *Session) Get(name string, creator func() interface{}) interface{} {
	o, ok := this.Vars[name]
	if !ok {
		if creator == nil {
			creator = this.defVars[name]
		}
		if creator != nil {
			o = creator()
			this.Vars[name] = o
		}
	}
	return o
}

func (this *Session) HasVarFactory(name string) bool {
	if _, ok := this.defVars[name]; ok {
		return true
	}
	return false
}

func (this *Session) RegVarFactory(name string, creator func() interface{}) {
	this.defVars[name] = creator
}

func (this *Session) Write(msg string) bool {
	return this.writer.Write(msg)
}

func (this *Session) Writeln(msg string) bool {
	return this.Write(msg + "\n")
}

func notifyOne(lis SessionCloseListener) {
	defer func() {
		recover()
	}()
	lis()
}

func (this *Session) notifyAll() {
	for _, lis := range this.listeners {
		notifyOne(lis)
	}
}

func (this *Session) Close() {
	this.writer.Close()
	if len(this.Vars) > 0 {
		this.Vars = make(map[string]interface{})
	}
}

func (this *Session) AddCloseListener(lis SessionCloseListener) {
	this.listeners = append(this.listeners, lis)
}

type ShellProcessor interface {
	Process(session *Session, command string) bool
}

type ShellCommandHandler func(session *Session, cmd string) bool

type shellProcessorFunc struct {
	handler ShellCommandHandler
}

func (this *shellProcessorFunc) Process(session *Session, command string) bool {
	return this.handler(session, command)
}

type NameSupport interface {
	Name() string
}

type Shell struct {
	PresetProcessors []ShellProcessor
	root             *ShellDirCommon
	version          string
}

func NewShell(appName string) *Shell {
	r := new(Shell)
	r.root = NewShellDirCommon("root")

	r.PresetProcessors = make([]ShellProcessor, 0)
	r.PresetProcessors = append(r.PresetProcessors, NewCloseCommand())
	r.PresetProcessors = append(r.PresetProcessors, newVersionCommand(appName, func() string { return r.version }))
	r.PresetProcessors = append(r.PresetProcessors, NewCDCommand())
	r.PresetProcessors = append(r.PresetProcessors, NewLSCommand())
	r.PresetProcessors = append(r.PresetProcessors, NewRunFileCommand(r))
	return r
}

type configInfo struct {
	Version string
}

func (this *Shell) Name() string {
	return "shell"
}

func (this *Shell) Init() bool {
	var cfg configInfo
	if config.GetBeanConfig("shell", &cfg) {
		this.version = cfg.Version
	}
	return true
}

func (this *Shell) AddCommand(o ShellProcessor) {
	this.root.Commands = append(this.root.Commands, o)
}

func (this *Shell) AddDir(o ShellDir) {
	this.root.Dirs = append(this.root.Dirs, o)
}

func (this *Shell) Process(session *Session, command string) bool {
	if command == "" {
		return true
	}
	if session == nil {
		logger.Error(tag, "session is nil")
		return false
	}
	session.Execute = func(command string) bool {
		return this.doProcess(session, command, true)
	}
	return this.doProcess(session, command, false)
}

func (this *Shell) lastCommand(session *Session) string {
	v := session.Get("@LASTCOMMAND", nil)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (this *Shell) doProcess(session *Session, command string, subExecute bool) (isProcess bool) {
	isHistory := false
	if command == "/" {
		command = this.lastCommand(session)
		isHistory = true
	}
	if command == "" {
		return true
	}
	defer func() {
		if isProcess && !subExecute && !isHistory {
			session.Vars["@LASTCOMMAND"] = command
		}
	}()

	if logger.EnableDebug(tag) {
		o := session.Get("@WHO", nil)
		who := "UNKNOW"
		if o != nil {
			who = o.(string)
		}
		logger.Debug(tag, "%s >> %s", who, command)
	}
	dir := session.Get("@DIR", func() interface{} {
		r := newDirInfo(this.root)
		return r
	}).(*dirInfo)
	for _, p := range this.PresetProcessors {
		if p.Process(session, command) {
			return true
		}
	}
	editor := GetEditor(session)
	if editor != nil {
		if editor.Process(session, command) {
			return true
		}
	}

	cname := CommandWord(command)
	cobj := dir.pwd.GetCommand(cname)
	if cobj != nil {
		if cobj.Process(session, command) {
			return true
		}
	}
	session.Writeln("ERROR: unknow command '" + cname + "'")
	return false
}
