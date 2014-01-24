package xmemservice

import "esp/shell"

func (this *Service) NewShellDir() shell.ShellDir {
	r := new(dirService)
	r.InitDir(this)
	return r
}
