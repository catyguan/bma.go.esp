local path = httpserv.path()

if path=="run.tc" then
	local m = httpserv.formValue("m","")
	local f = httpserv.formValue("f","")
	if m=="" then m = "_" end
	return go.exec(m..":testcase/"..f..".tc.lua")	
end

local ext = filepath.ext(path)
if ext==".gl" then
	local npath = filepath.changeExt(path, ".gl.lua")
	return go.exec("_:http/"..npath)
end

return httpserv.writeFile("_:http.res/"..path)
