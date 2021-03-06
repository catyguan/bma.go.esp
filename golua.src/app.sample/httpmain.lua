local path = httpserv.path()

local m = "_"
local script = path
local dir, fname = filepath.dir(path)
-- print(dir, fname)
-- if dir=="service" then
-- 	m = "glf.testing"	
-- 	script = fname
-- end

local ext = filepath.ext(script)
if ext==".gl" then
	local nscript = filepath.changeExt(script, ".gl.lua")
	return go.exec(m..":http/"..nscript)
end

if fname=="favicon.ico" then
	return {Status=404}
end

return httpserv.writeFile("_:http.res/"..script)
