require("_:service/NodesManager.lua")

local s = class.new("NodesManager")
local formData = s.DM_NODE.BuildForm({},nil)

httpserv.render("_:http/nodeForm.view",{id=0, node=formData})