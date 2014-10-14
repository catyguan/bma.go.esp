local ech = go.chan(3)
go.defer(function()
	closure(ech)
	go.close(ech)
end)

local v = 1
local function c(n, ech)
	closure(v)
	return function()
		closure(n, v, ech)
		for i = 1,10 do
			v = v + 1
			print(n, v)
		end
		-- go.write(ech, n)
	end
end

go.run(c("g1", ech))
go.run(c("g2", ech))
go.run(c("g3", ech))

local c = 0
while c<3 do
	local n = go.read(ech)
	c = c + 1
	print(n, "stop", c)
end

print("end")
