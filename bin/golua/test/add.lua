local form = httpserv.form()
local a, b, c

a = types.int(form,"a",0)
b = types.int(form,"b",0)
c = a+b

acclog.log("ask-c", c)

-- print("here", httpserv.header())

local str = "result = "..types.string(c)
httpserv.write(str)
-- return str