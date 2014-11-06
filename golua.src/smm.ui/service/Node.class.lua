-- smm.ui:service/Node.class.lua
require("_:service/NodesManager.class.lua")

local Class = class.define("smm.ui.Node")

function Class.Bind(nodeId, id, aid)
	self.nodeId = nodeId
	self.id = id
	self.aid = aid
end

function Class.Invoke(param, ctx)
	local s = class.new("NodesManager")
	local nodeInfo = s.Get(self.nodeId)
	if nodeInfo==nil then
		error("unknow node(%d)", self.nodeId)
	end
	local req = {
		URL = nodeInfo.api_url
	}
	local data
	if ctx==nil then
		data = {}
	else
		data = table.clone(ctx)
	end
	data["_id"] = self.id
	data["_aid"] = self.aid
	data["_param"] = param
	req.Data = data

	local resp = httpclient.exec(req)
	local co = resp.Content
	if co==nil then co = "null" end
	if resp.Status~=200 then		
		error("invoke fail(%d, %s)", resp.Status, co)
	end
	local r = json.decode(co)
	if r==nil then
		error("invalid response content (%s)", co)
	end
	if r.Status~=200 then
		error("response fail: %v", r.Result)
	end
	return r.Result
end