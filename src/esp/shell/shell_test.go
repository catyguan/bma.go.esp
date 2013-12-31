package shell

import (
	"fmt"
	"testing"
)

type testProcessor struct {
}

func (this testProcessor) Process(session *Session, command string) bool {
	fmt.Println("handle ", command)
	return true
}

func (this testProcessor) Name() string {
	return "test"
}

func TestShellParse(t *testing.T) {
	var a bool
	fs := NewFlagSet("test")
	fs.BoolVar(&a, "a", false, "show all")
	t.Log(GetHelp("ls", "files ...", fs))
	cmd, err := Parse(fs, "ls -a test test2 $$select * from table1 where a='hello'")
	t.Log("command = ", cmd)
	t.Log("a = ", a)
	t.Log(fs.Args())
	t.Log(fs.Arg(2))
	if err != nil {
		t.Log(err)
	}
}

func TestShell(t *testing.T) {
	// logger.Config().InitLogger()

	sh := NewShell("app")

	dir1 := NewShellDirCommon("dir1")
	sh.AddDir(dir1)

	sh.AddCommand(new(testProcessor))

	session := NewSession(NewConsoleWriter())
	session.Vars["@WHO"] = "TestCase"
	session.AddCloseListener(func() {
		fmt.Println("i know session close!!!")
	})

	sh.Process(session, "hello")
	sh.Process(session, "/")

	sh.Process(session, "cd")
	sh.Process(session, "cd dir1")

	sh.Process(session, "cd ..")
	sh.Process(session, "ls")

	sh.Process(session, "test hello")

	session.Close()
}

func TestShellVersion(t *testing.T) {
	// logger.Config().InitLogger()

	sh := NewShell("app")
	sh.version = "1.2.3"

	session := NewSession(NewConsoleWriter())
	sh.Process(session, "version -m -c app 1.3.1")
	session.Close()
}
