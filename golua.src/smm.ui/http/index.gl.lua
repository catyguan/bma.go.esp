require("_:service/NodesManager.class.lua")

local form = httpserv.form()
local s = class.new("NodesManager")
local nodes = s.List()

httpserv.render("_:http/index.view.htm",{nodes=nodes})