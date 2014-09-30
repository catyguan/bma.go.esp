print("module 'test/shm' init")

function service_test_shm(ctx, result)
	print("lua.shm")
	v = glua_getInt(ctx, "v")
	req = {
		Key = "test",
	}
	if q==1 then
		req["Set"] = true
		req["Value"] = v
		req["Size"] = 4
	end
	
	greq = glua_toMap(req)
	err = glua_task("shm",greq,"test_shm_endReq")
	if err~=nil then
		error(err)
	end
end

function test_shm_endReq(ctx, result)
	print("test_shm_endReq")
	m = glua_getMap(result,"shm")
	str = glua_getString(m, "Content")
	-- hs = {}
	-- hs["Content-Type"] = "text/html; charset=utf-8"
	-- glua_setMap(result, "Header", hs)
	glua_setString(result, "Content-Type", "text/html; charset=utf-8")
	glua_setString(result, "Content", str)
	return true
end