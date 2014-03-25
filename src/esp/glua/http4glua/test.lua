print("module http.test init")

function http(ctx, result)
	print("lua.http")
	req = {
		URL = "http://cn.bing.com/search",
		Data = {
			q="golang http"
		}
	}
	greq = glua_toMap(req)
	err = glua_task("http",greq,"endReq")
	if err~=nil then
		error(err)
	end
end

function endReq(ctx, result)
	print("lua.endReq")
	m = glua_getMap(result,"http")
	glua_setNil(m, "Content")
	glua_setNil(m, "Header")
	return true
end
