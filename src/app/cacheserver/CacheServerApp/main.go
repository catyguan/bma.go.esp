package main

import (
	"app/cacheserver"
	"app/cacheserver/httppoint"
	// "app/cacheserver/sqlloader"
	"app/cacheserver/mcpoint"
	"app/cacheserver/thriftpoint"
	"boot"
	// "esp/namedsql"
	"esp/shell"
	"esp/shell/mcservercmd"
	"esp/shell/telnetcmd"
	"esp/shell/thriftcmd"
	"esp/sqlite"
	"mcserver"
	"telnetserver"
	"thrift"
	/*
		plugin modules
	*/
	// "github.com/ziutek/mymysql/godrv"
	_ "app/cacheserver/eloader"
	_ "app/cacheserver/elru"
	_ "app/cacheserver/esimple"
	_ "app/cacheserver/estable"
	// _ "github.com/go-sql-driver/mysql"
)

func main() {
	// mymysql
	// godrv.Register("set names utf8")
	// tcp:172.19.16.195:3306*db_live2/root/root
	// mysql
	// root:root@tcp(172.19.16.195:3306)/db_live2

	cfile := "config/cacheserver-config.json"

	shl := shell.NewShell("cacheserver")
	boot.Install("shell", shl)

	// namedsql
	// namedSQL := namedsql.NewNamedSQL("namedSQL")
	// namedSQL.DefaultBoot()
	// shl.AddDir(namedSQL.NewShellDir())

	// sqliteServer
	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()
	shl.AddDir(sqliteServer.NewShellDir())

	// CacheServer
	cacheService := cacheserver.NewCacheService("cacheService", sqliteServer)
	cacheService.DefaultBoot()
	shl.AddDir(cacheService.NewShellDir())

	// ThriftPoint
	handler := thriftpoint.NewTCacheServerImpl(cacheService)
	processor := thriftpoint.NewTCacheServerProcessor(handler)
	thriftServer := thrift.NewThriftServer("thriftPoint", processor)
	thriftServer.DefaultBoot(true)
	shl.AddDir(thriftcmd.NewShellDir(thriftServer))

	// HttpPoint
	httpPoint := httppoint.NewHttpPoint("httpPoint", cacheService)
	httpPoint.DefaultBoot()

	// MemcachePoint
	mcPoint := mcpoint.NewMemcachePoint("mcPoint", cacheService)
	mcServer := mcserver.NewMemcacheServer("mcPoint", mcPoint.Handle)
	mcServer.DefaultBoot()
	mcPoint.DefaultBoot()
	mcDir := mcservercmd.NewShellDir(mcServer)
	mcPoint.BuildShellDir(mcDir)
	shl.AddDir(mcDir)

	// telnetServer
	tServer := telnetserver.NewTelnetServer("telnetServer", telnetcmd.NewHandler(shl))
	tServer.DefaultBoot(true)
	shl.AddDir(telnetcmd.NewShellDir(tServer))

	// sqlloader.InitSQLLoader(namedSQL, "")

	boot.Go(cfile)
}
