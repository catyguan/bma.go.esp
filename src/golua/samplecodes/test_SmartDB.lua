local vr

local sdb = go.new("SmartDB")
sdb.Add({Name="test", Driver="mysql", DataSource="root:root@tcp(172.19.16.195:3306)/test_db"})

local db = sdb.Select("test4")
db.Ping()

local rs = db.Query("select * from test4")
-- go.defer(rs)

local data
local desc = {id="int",dt="time"}
while rs.Fetch(data, desc) do
	print(data)
	print(types.name(data["id"]))
end

return vr