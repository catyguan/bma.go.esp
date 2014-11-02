local vr

local cls1 = class.define("test1")
function cls1.ctor()
	print("in ctor1")
end
function cls1.hello()
	print("hello world")
end
print(cls1.ClassName)

local cls2 = class.define("test2", ["test1"])
function cls2.ctor(n)
	print("in ctor2")
	self._name = n
end
function cls2.name()
	return self._name
end
cls2.hello()

local cls3 = class.define("i3")
cls3.name = 0
cls3.helloWorld = 0

local o = class.new("test2")
o._name = "abc"
print(o.name())

local o = cls2.New("ppp")
print(o.name())

print("is", class.is(o, "test1"))
print("check", class.check(o,"i3"))

return vr