local vr

os.mkdir("e:/gocreate/sub", true)
local f = os.createFile("e:/gocreate/test.txt")
go.defer(f)
f.Write("hello kitty")

-- os.rename("e:/gocreate","e:/gocreate2")
-- os.remove("e:/gocreate2", true)

return vr