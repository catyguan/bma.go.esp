require("_:service/NodesManager.class.lua")
require("_:http/NodeForm.class.lua")

local form = httpserv.form()
local s = class.new("NodesManager")
local fo = class.new("NodeForm")

fo.Parse(form)
local ok = fo.Valid()

if not ok then
	httpserv.render("_:http/nodeForm.view.htm",{node=fo})
	return
end

local id = fo.Value("id")
local data = fo.CloneData()
table.remove(data, "id")
if id>0 then
	s.Update(data, id)
else
	s.Insert(data)
end

httpserv.render("_:http/done.view.htm",{msg="Submit Done"})


