require("_:service/NodesManager.class.lua")
require("_:http/NodeForm.class.lua")

local fo = class.new("NodeForm")
fo.Bind({})

httpserv.render("_:http/nodeForm.view.htm",{node=fo})