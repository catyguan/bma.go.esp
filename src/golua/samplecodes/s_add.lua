require("m_add.lua")
-- require("m_add.lua")

local a, b, c
-- print(_REQUEST)
a = types.int(_REQUEST,"a",0)
b = types.int(_REQUEST,"b",0)
c = add(a, b)
print("中文", types.string(c))
return c