print("module test init")

function hello( ctx, result )
	print("hello world, i will timeout!!!")
	-- error("fuck")
	-- print("hello world")
	-- return true
end

function add( ctx, result )
	print("lua.add")
	a = glua_getInt(ctx, "a")
	b = glua_getInt(ctx, "b")
	c = a + b
	print("a + b => ", c)
	glua_setInt(result, "c", c)
	return true
end

function async( ctx, result )
	print("lua.async")
	req = glua_newMap()
	err = glua_task("testpl",req,"endReq")
	if err~=nil then
		error(err)
	end
end

function endReq(ctx, result)
	return true
end

function all(ctx, result)
	print("lua.all")
	req = {
		Tasks={
			t1={
				Name="testpl"
			},
			t2={
				Name="testpl"
			},
		}
	}
	greq = glua_toMap(req)
	err = glua_task("all",greq,"endReq")
	if err~=nil then
		error(err)
	end
end