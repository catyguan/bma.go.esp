local vr

local db = sql.open("mysql", "root:root@tcp(172.19.16.195:3306)/test_db")
go.defer(function()
	closure(db)
	db:Close()
end)
-- sql.open("sqlite3","test.db")
db:Ping()

local rs = db:Query("select * from test4")
go.defer(rs)

local data
local desc = {id="int",dt="time"}
while rs:Fetch(data, desc) do
	print(data)
	print(types.name(data["id"]))
end

-- local res, lid
-- res = db:Exec("DELETE FROM test WHERE username='ppp'")
-- print("delete", res)

-- res, lid = db:ExecLastId("INSERT test VALUES(NULL,'ppp','123456')")
-- print("insert", res, lid)

-- local stmt = db:Prepare("UPDATE test SET password='newpass' WHERE username=?")
-- go.defer(stmt)
-- res = stmt:Exec("ppp")
-- print("Prepare & Exec", res)

-- local tx = db:Begin()
-- go.defer(tx)
-- local res = tx:Exec("UPDATE test SET password='newpass2' WHERE username=?", "ppp")
-- print(res)

return vr