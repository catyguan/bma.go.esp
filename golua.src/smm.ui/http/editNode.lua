require("_:service/NodesManager.lua")

local form = httpserv.form()
local id = types.int(form,"id",0)

local s = class.new("NodesManager")
local node = s.Get(id)

if node==nil then
	error(strings.format("Node(%d) not exists", id))
end

local formData = s.DM_NODE.BuildForm(node,nil)
httpserv.render("_:http/nodeForm.view",{id=id, node=formData})