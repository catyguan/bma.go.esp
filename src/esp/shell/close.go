package shell

import ()

type CloseCommand struct {
}

func NewCloseCommand() *CloseCommand {
	r := new(CloseCommand)
	return r
}

func (this *CloseCommand) Process(s *Session, command string) bool {
	cname := CommandWord(command)
	if cname == "close" || cname == "quit" || cname == "exit" || cname == "bye" {
		s.Close()
		return true
	}
	return false
}
