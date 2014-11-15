local path = httpserv.path()

local m = "_"
local script = path
local dir, fname = filepath.dir(path)
print(dir, fname)
if dir=="testcase" then
	m = "glf.testing"	
	script = fname
end

local ext = filepath.ext(script)
if ext==".gl" then
	local nscript = filepath.changeExt(script, ".gl.lua")
	return go.exec(m..":http/"..nscript)
end

return httpserv.writeFile("_:http.res/"..script)
