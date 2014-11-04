require("_:service/NodesManager.lua")

local form = httpserv.form()
local s = class.new("NodesManager")

local id = types.int(form, "id", 0)
local data = s.DM_NODE.Parse(form)
local ok, vdata = s.Valid(data)

if not ok then
	local formData = s.DM_NODE.BuildForm(data,vdata)
	httpserv.render("_:http/nodeForm.view",{id=id, node=formData})
	return
end

if id>0 then
	s.Update(data, id)
else
	s.Insert(data)
end

httpserv.render("_:http/done.view",{msg="Submit Done"})


