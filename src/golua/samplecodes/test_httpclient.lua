print("vmmhttp.httpclient")
-- req = {
-- 	URL = "http://cn.bing.com/search",
-- 	Data = {
-- 		q="golang http"
-- 	}
-- }

-- local resp = httpclient.exec(req)
-- resp["Content"] = strings.substr(resp["Content"], 0, 100)
-- print(resp)

local content = httpclient.getContent("http://cn.bing.com/search?q=golang+http")
return strings.substr(content, 0, 100)