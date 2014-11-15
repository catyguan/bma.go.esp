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
for k, v in TC do
	print(k, v, strings.hasPrefix(v, "test"))
	if strings.hasPrefix(k, "test") then		
		name = strings.substr(k, 4)
		title = name
		help = ""
		params = nil
		desc = nil
		ho = TC["help"..name]
		if ho~=nil then
			if not types.isEmpty(ho.Name) then name = ho.Name end
			if not types.isEmpty(ho.Title) then title = ho.Title end
			if not types.isEmpty(ho.Help) then help = ho.Help end
			params = ho.Params
			desc = ho.Desc
		end
		table.insert(rs, {
			Name = name,
			Title = title,
			Help = help,
			Params = params,
			Desc = desc
		})
	end
end
local r = {
	Result = rs
}
return {
	Content = json.encode(r)
}

