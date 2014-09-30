print("module 'test/http' init")

function service_test_http(ctx, result)
	print("lua.http")
	req = {
		URL = "http://cn.bing.com/search",
		Data = {
			q="golang http"
		}
	}
	greq = glua_toMap(req)
	err = glua_task("http",greq,"test_http_endReq")
	if err~=nil then
		error(err)
	end
end

function test_http_endReq(ctx, result)
	print("lua.endReq2")
	m = glua_getMap(result,"http")
	str = glua_getString(m, "Content")
	-- hs = {}
	-- hs["Content-Type"] = "text/html; charset=utf-8"
	-- glua_setMap(result, "Header", hs)
	glua_setString(result, "Content-Type", "text/html; charset=utf-8")
	glua_setString(result, "Content", str)
	return true
end