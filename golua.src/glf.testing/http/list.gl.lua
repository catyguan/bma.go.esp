local form = httpserv.form()
local n = types.string(form,"n","")

if n=="" then
	return {
		Content = json.encode({
			Message = "param 'n' empty"
		})
	}
end

local TC = go.exec(n..".tc.lua", {})
local rs = []
local name, title, help,params, ho, desc
for k, f in TC do
	-- print(k, v, strings.hasPrefix(v, "test"))
	if strings.hasPrefix(k, "test") then		
		name = strings.substr(k, 4)
		title = go.annoGet(f,"title")
		if title=="" then
			title = name
		end
		help = go.annoGet(f,"help")
		params = []
		annParams = go.annoList(f,"params")
		for _,annP in annParams do
			table.insert(params, json.decode(annP))
		end
		desc = go.annoGet(f,"desc")
		table.insert(rs, {
			Name = name,
			Title = title,
			Help = help,
			Params = params,
			Desc = desc
		})
	end
end
table.sort(rs, function(o1, o2)
	return o1.Name < o2.Name
end)
local r = {
	Result = rs
}
return {
	Content = json.encode(r)
}

