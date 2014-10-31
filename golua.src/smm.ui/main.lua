local path = httpserv.path()
return go.exec("_:"..path..".lua")