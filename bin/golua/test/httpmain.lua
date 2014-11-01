print("Hi, i'm dispatcher")

local path = httpserv.path()
if path=="/favicon.ico" then
	return {Status=404}
end
return go.exec(path..".lua")