local r = TCS
if r==nil then r = [] end
table.insert(r, {
	Script = "go.memserv:testcase/Base",
	Title = "go.memserv base TC",
	Help = "基础功能"
})
return r