require("_:service/NodesManager.class.lua")
require("_:http/NodeForm.class.lua")

local form = httpserv.form()
local id = types.int(form,"id",0)

local s = class.new("NodesManager")
local node = s.Get(id)

if node==nil then
	error(strings.format("Node(%d) not exists", id))
end

local fo = class.new("NodeForm")
fo.Bind(node)
httpserv.render("_:http/nodeForm.view.htm",{node=fo})