require("local:config.lua")

local sdb = go.new("SmartDB")
local sdbConfigList = config.get("smartdb.list")
for cfg in sdbConfigList do
	sdb.Add(cfg, true)
end