local path = httpserv.path()
local ext = filepath.ext(path)

if ext==".gl" then
	local npath = filepath.changeExt(path, ".lua")
	return go.exec("_:http/"..npath)
end

return httpserv.writeFile("_:http/"..path)
