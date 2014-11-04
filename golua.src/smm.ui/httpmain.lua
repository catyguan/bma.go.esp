local path = httpserv.path()
local ext = filepath.ext(path)

if path=="run.tc" then
	local m = httpserv.formValue("m","")
	local f = httpserv.formValue("f","")
	if m=="" then m = "_" end
	return go.exec(m..":testcase/"..f.."TC.lua")	
end

if ext==".gl" then
	local npath = filepath.changeExt(path, ".lua")
	return go.exec("_:http/"..npath)
end

return httpserv.writeFile("_:http/"..path)
