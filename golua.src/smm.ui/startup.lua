local sdb = go.new("SmartDB")
local cfg = {
	Name="test",
	Driver="mysql", DataSource="root:root@tcp(172.19.16.195:3306)/test_db",
	MaxIdleConns=10, MaxOpenConns=100,
	ReadOnly=false, Priority=5
}
sdb.Add(cfg, true)