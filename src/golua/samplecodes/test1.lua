-- local v1 = 100
-- local v2 = 1+2*3+v1
-- print("helloWorld = ", v1, v2)
-- v3 = o["e"]
local vr, v2
vr = [1,2,3,4,5]
print(vr[2])
vr[3] = nil
vr[5] = "add"
-- v2 = {a=1,b=2.0,c=true,d="hello",["e"]=3}
-- vr = v2["e"]

-- local v1
-- repeat
-- 	v1 = v1 + 1
-- 	print("in", v1)
-- until v1>3
-- print("out", v1)

-- local v1 = 0
-- while v1 < 3 do
-- 	v1 = v1 + 1
-- 	print("in", v1)
-- end
-- print("out", v1)

-- local v3,v4
-- local v1,v5 = 2

-- local v1,a,b,c
-- v1 = "hello"
-- a, b, c = 1, 2
-- print(v1, a, b, c)

-- print("hello world", v1, 1)
-- local obj = {
-- 	p = print
-- }
-- obj:p("hello world")
-- obj.parent:print(1 + 2, true, a.b)
-- a.b = 1 + 2 - 3

-- functions
-- local c = 3
-- local function f1(a, b)
-- 	closure(c)
-- 	return a + b + c
-- end
-- vr = f1(1, 2)
-- local o = {}
-- function o:f2(a,b)
-- 	b = 2
-- end
-- local function f3(a,b,...)
-- 	print(a,b,...)
-- end
-- f3(1,2,3,4,5)
-- vr = o

-- for i = 1,5,"a" do
-- 	print(i)
-- end
-- for i,v in 1,2,4,7 do
-- 	print(i, v)
-- end
-- for _,v in v2 do
-- 	print(_, v)
-- end

-- print(not false)
-- print(not false, #"abc" , -120)

-- local a, b, c
-- if a then
-- 	a = 1
-- elseif b then
-- 	b = 2
-- else
-- 	c = 3
-- end
-- print("hello",a,b,c)

-- stack overflow
-- function f()
-- 	f()
-- end
-- f()
-- while true do
-- end

-- metatable
-- local mt = { 
-- 	__index = function(t, k)
-- 		return 123
-- 	end,
-- 	__newindex = function(t, k, v)
-- 		print("newindex", k, v)
-- 		rawset(t, k, v+123)
-- 	end
-- }
-- local o = {}
-- setmetatable(o, mt)
-- o.ab = o.abc *1000
-- vr = o.ab

-- local n = 123
-- pcall(function()
-- 	closure(n)
-- 	print("here", n)
-- end)

return vr
-- return 1 + 2