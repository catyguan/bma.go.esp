print("Hi, i'm dispatcher")

local path = httpserv.path()
return go.exec(path..".lua")