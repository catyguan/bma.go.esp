require("_:service/NodesManager.class.lua")

local form = httpserv.form()
local rf = types.int(form, "r", 0)
local s = class.new("NodesManager")
local nodes = s.List(rf>0)

httpserv.render("_:http/index.view.htm",{nodes=nodes})