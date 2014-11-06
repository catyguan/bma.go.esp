require("_:service/NodesManager.class.lua")
require("_:service/Node.class.lua")

local form = httpserv.form()
local id = types.int(form,"id",0)

local s = class.new("NodesManager")
local node = s.Get(id)

if node==nil then
	error(strings.format("Node(%d) not exists", id))
end

local ns = class.new("smm.ui.Node")
ns.Bind(id, "", "list")
local infoList = ns.Invoke("", nil)

httpserv.render("_:http/node.view.htm",{node=node, infoList=infoList})