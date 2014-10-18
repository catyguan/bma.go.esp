local vr

local tb = [1,2,3,4,5]
-- vr = table.concat(tb)
-- table.insert(tb, 2, 6)
-- table.remove(tb, 2)
-- print(tb)
-- tb = {a=1,b=2}
-- table.remove(tb, "a")
local tb2 = table.subtable(tb, 1, 3)
print(tb2)

return vr