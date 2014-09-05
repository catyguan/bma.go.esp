print("module 'test/add' init")

function service_test_add( ctx, res )
	print("lua.add")
	a = glua_getInt(ctx, "a")
	b = glua_getInt(ctx, "b")
	c = a + b
	print("a + b => ", c)
	glua_setString(res, "Content", ""..c)
	return true
end