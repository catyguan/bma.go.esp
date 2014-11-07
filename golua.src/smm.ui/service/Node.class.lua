-- smm.ui:service/Node.class.lua
require("_:service/NodesManager.class.lua")

local Class = class.define("smm.ui.Node")

function Class.Bind(nodeId, id, aid)
	self.nodeId = nodeId
	self.id = id
	self.aid = aid
end

function Class.Invoke(param)
	local s = class.new("NodesManager")
	local nodeInfo = s.Get(self.nodeId)
	if nodeInfo==nil then
		error("unknow node(%d)", self.nodeId)
	end
	local req = {
		URL = nodeInfo.api_url
	}
	local strparam = ""
	if not types.isEmpty(param) then
		strparam = json.encode(param)
		req.Post = true
	end
	local data = {}
	data["id"] = self.id
	data["aid"] = self.aid
	data["param"] = strparam
	if nodeInfo.code~='' then
		local tmp = strings.format("%d/%s/%s/%s", self.id, self.aid, strparam, nodeInfo.code)
		data["code"] = strings.md5(tmp)
	end
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