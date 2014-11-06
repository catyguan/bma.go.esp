require("_:service/NodesManager.class.lua")
require("_:service/Node.class.lua")

local form = httpserv.form()
local nid = types.int(form,"nid",0)
local id = types.string(form,"id","")
local aid = types.string(form,"aid","")
local param = types.string(form,"param","")

local s = class.new("NodesManager")
local node = s.Get(nid)

if node==nil then
	error(strings.format("Node(%d) not exists", nid))
end
if aid=="" then
	error(strings.format("Node(%d) action empty", nid))
end

local ns = class.new("smm.ui.Node")
ns.Bind(nid, id, aid)
local res = ns.Invoke("", nil)
local r = {
	Result = res
}
return {
	["Content-Type"] = "application/json",
	Content = json.encode(r)
}