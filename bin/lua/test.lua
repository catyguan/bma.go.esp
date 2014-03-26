print("module 'test' init")

function service_add( ctx, res )
	print("lua.add")
	a = glua_getInt(ctx, "a")
	b = glua_getInt(ctx, "b")
	c = a + b
	print("a + b => ", c)
	glua_setInt(res, "Content", c)
	return true
end