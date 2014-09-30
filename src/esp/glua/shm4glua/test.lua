print("module shm.test init")

function shm(ctx, result)
	print("lua.shm")
	req = {
		-- Set = true,
		Delete = true,
		Key = "test",
		Value = 123,
		Size = 8
	}
	greq = glua_toMap(req)
	err = glua_task("shm",greq,"endReq")
	if err~=nil then
		error(err)
	end
end

function endReq(ctx, result)
	print("lua.endReq")
	-- m = glua_getMap(result,"http")
	-- glua_setNil(m, "Content")
	-- glua_setNil(m, "Header")
	return true
end
