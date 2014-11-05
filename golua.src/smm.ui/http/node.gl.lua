require("_:service/NodesManager.class.lua")

local form = httpserv.form()
local id = types.int(form,"id",0)

local s = class.new("NodesManager")
local node = s.Get(id)

if node==nil then
	error(strings.format("Node(%d) not exists", id))
end

httpserv.render("_:http/node.view.htm",{node=node})