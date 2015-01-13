local project = gom.GetItem("project")
print("project.name", project["name"].Value)
for k,item in project do
	go.debug("test.lua","project.%s = %v",k, item.Value)
end

-- gomf.render("test.tpl", {gom=gom},"test.out")