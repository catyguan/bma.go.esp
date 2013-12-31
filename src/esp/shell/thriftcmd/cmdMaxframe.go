package thriftcmd

import (
	"bmautil/valutil"
	"esp/shell"
	"thrift"
)

type maxFrameCommand struct {
	server *thrift.ThriftServer
}

func newMaxframeCommand(s *thrift.ThriftServer) *maxFrameCommand {
	r := new(maxFrameCommand)
	r.server = s
	return r
}

func (this *maxFrameCommand) Name() string {
	return "maxframe"
}

func (this *maxFrameCommand) Process(s *shell.Session, command string) bool {
	if shell.CommandWord(command) != this.Name() {
		return false
	}

	name := this.Name()
	args := "maxFrameSize"
	fs := shell.NewFlagSet(name)
	if shell.DoParse(s, command, fs, name, args, 0, 1) {
		return true
	}
	if fs.NArg() > 0 {
		val := fs.Arg(0)
		mf, err := valutil.ToSize(val, 1024, valutil.SizeB)
		if err != nil {
			s.Writeln("ERROR: " + err.Error())
			return true
		}
		this.server.Maxframe = mf
	}
	str := valutil.SizeString(this.server.Maxframe, 1024, valutil.SizeM)
	s.Writeln("max frame -> " + str)
	return true
}
