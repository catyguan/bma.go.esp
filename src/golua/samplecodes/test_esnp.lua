local vr

local esnp = go.new("ESNP")
-- local sock = esnp.Create("127.0.0.1:1091", 2000)
-- go.defer(function( ... )
-- 	closure(sock)
-- 	sock.Close()
-- end)
local sock = esnp.Open("127.0.0.1:1091","", 2000)
print("here")
sock = esnp.Open("127.0.0.1:1091","", 2000)
print("here2")
local msg = {	
	Address = {
		Service="smm.api",
		Op="invoke"
	},
	Data = {
		id = "go.server",
		aid="boot.reload",
		param=""
	}
}
vr = sock.Call(msg, 1000)

return vr