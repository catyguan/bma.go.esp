package clumem

import "esp/shell"

func (this *Service) NewShellDir() shell.ShellDir {
	r := new(dirService)
	r.InitDir(this)
	return r
}

func (this *Service) BuildMemGroupCommands(name string, dir *shell.ShellDirCommon) {
	// dir.DirInfoFunc = func() string {
	// 	r, err := this.QueryStats(name)
	// 	if err != nil {
	// 		return ""
	// 	}
	// 	return r
	// }

	// cmd1 := &cmdDelete{this, name}
	// dir.AddCommand(cmd1)
	// cmd2 := &cmdGet{this, name}
	// dir.AddCommand(cmd2)
	// cmd3 := &cmdSet{this, name}
	// dir.AddCommand(cmd3)
	// cmd4 := &cmdStats{this, name}
	// dir.AddCommand(cmd4)
}
