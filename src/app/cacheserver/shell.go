package cacheserver

import (
	"esp/shell"
)

type dirService struct {
	service *CacheService
}

func (this *dirService) Name() string {
	return this.service.Name()
}

func (this *dirService) GetCommand(name string) shell.ShellProcessor {
	switch name {
	case commandNameDeleteCache:
		return &cmdDeleteCache{this.service, ""}
	case commandNameNewCache:
		return &cmdNewCache{this.service}
	case commandNameSave:
		return &cmdSave{this.service}
	case commandNameStartCache:
		return &cmdStartCache{this.service}
	case commandNameStopCache:
		return &cmdStopCache{this.service}
	}
	return nil
}

func (this *dirService) GetDir(name string) shell.ShellDir {
	cache, _ := this.service.GetCache(name, false)
	if cache != nil {
		ss, ok := cache.(shell.ShellDirSupported)
		if ok {
			return ss.CreateShell()
		}
	}
	return nil
}

func (this *dirService) List() []string {
	clist := this.service.ListCacheName()
	r := make([]string, 0)
	for _, k := range clist {
		r = append(r, shell.DirName(k))
	}
	cmds := []string{
		commandNameDeleteCache,
		commandNameNewCache,
		commandNameSave,
		commandNameStartCache,
		commandNameStopCache,
	}
	for _, k := range cmds {
		r = append(r, k)
	}
	return r
}

func (this *CacheService) NewShellDir() shell.ShellDir {
	return &dirService{this}
}

func (this *CacheService) BuildCacheCommands(name string, dir *shell.ShellDirCommon) {
	dir.DirInfoFunc = func() string {
		r, err := this.QueryStats(name)
		if err != nil {
			return ""
		}
		return r
	}

	cmd1 := &cmdDelete{this, name}
	dir.AddCommand(cmd1)
	cmd2 := &cmdGet{this, name}
	dir.AddCommand(cmd2)
	cmd3 := &cmdSet{this, name}
	dir.AddCommand(cmd3)
	cmd4 := &cmdStats{this, name}
	dir.AddCommand(cmd4)
}
