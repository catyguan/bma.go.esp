local ch = go.chan(1)
go.defer(ch)
local ech = go.chan(3)
go.defer(ech)

local v = 1
go.enableSafe(v)
local function c(n, ch, ech)
	closure(v)
	return function()
		closure(n, v, ch, ech)
		go.write(ch, true)
		for i = 1,10 do
			v = v + 1
			print(n, v)
		end
		go.write(ech, n)
	end
end

go.write(ch, true)

go.run(c("g1", ch, ech))
go.run(c("g2", ch, ech))
go.run(c("g3", ch, ech))

go.read(ch)
go.read(ch)
go.read(ch)

local c = 0
while c<3 do
	local n = go.read(ech)
	c = c + 1
	print(n, "stop", c)
end

print("end")
