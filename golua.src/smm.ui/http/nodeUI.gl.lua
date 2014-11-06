require("_:service/NodesManager.class.lua")
require("_:service/Node.class.lua")

local form = httpserv.form()
local nid = types.int(form,"nid",0)
local id = types.string(form,"id","")
local aid = types.string(form,"aid","")
local uin = types.string(form,"uin","")

local s = class.new("NodesManager")
local node = s.Get(nid)

if node==nil then
	error(strings.format("Node(%d) not exists", nid))
end
if aid=="" then
	error(strings.format("Node(%d) action empty", nid))
end
if uin=="" then
	error(strings.format("Node(%d) uin empty", nid))
end

return go.exec(uin, {node=node, nodeId=nid, actionId=aid, id=id})