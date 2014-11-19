local vr

local mems = go.new("MemServ")
local mg = mems.Get("app")
print(mg.Call(function(mgi, x)
	mgi.Set("test", 123)
	go.sleep(100)	
	print("------->")
	print(mgi.Get("test"))
	print(mgi.Put("test2", 123))
	print(mgi.Put("test2", 234567))
	-- error("fuck")
	return "here",2,3, x
end, true))


print(mg.Size())

mg.Scan("begin", "test")
while true do
	local isend, arr = mg.Scan("do", "test", 1)
	print(isend, arr)
	for item in arr do
		print("scan", item.Key, item.Value)
	end
	if isend then
		break
	end
end

return vr