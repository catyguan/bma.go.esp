local vr

local mems = go.new("MemServ")
local mg = mems.Get("app")
mg.Set("test", 123)
go.sleep(100)
print(mg.Get("test"))
print(mg.Size())

return vr