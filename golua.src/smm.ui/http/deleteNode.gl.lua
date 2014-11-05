require("_:service/NodesManager.class.lua")

local form = httpserv.form()
local id = types.int(form,"id",0)

local s = class.new("NodesManager")
local c = s.Delete(id)

local msg = "Delete Done"
if c~=1 then
	msg = "Delete Fail"
end

httpserv.render("_:http/done.view.htm",{msg=msg})