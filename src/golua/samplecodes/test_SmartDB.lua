local vr

local sdb = go.new("SmartDB")
local cfg = {
	Name="test",
	Driver="mysql", DataSource="root:root@tcp(172.19.16.195:3306)/test_db",
	MaxIdleConns=10, MaxOpenConns=100,
	ReadOnly=false, Priority=5
}
sdb.Add(cfg)

local db = sdb.Select("tEST4")
db.Ping()

local rs = db.Query("select * from test4")
-- go.defer(rs)

local data
local desc = {id="int",dt="time"}
while rs.Fetch(data, desc) do
	print(data)
	print(types.name(data["id"]))
end

-- sdb.X()

return vr