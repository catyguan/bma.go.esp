local form = httpserv.form()
local n = types.string(form,"n","")
local f = types.string(form,"f","")

if n=="" or f=="" then
	return {
		Content = json.encode({
			Message = "param 'n' or 'f' empty"
		})
	}
end

local params = {}
local name
for k, v in form do
	if strings.hasPrefix(k,"p_") then
		name = strings.substr(k, 2)
		params[name] = v
	end	
end

local TC = go.exec(n..".tc.lua", {})
local func = TC["test"..f]
if func==nil then
	return {
		Content = json.encode({
			Message = "test method not exists"
		})
	}
end

require("glf.testing:service/T.class.lua")

local T = class.new("glf.testing.T")
local ok, err
ok, err = pcall(function(func, T, params)
	func(T, params)
end,func, T, params)

local r = {
	Result = "OK"
}
if not ok then
	r.Result = strings.format("%s", err)
end
r.Log = T.GetLogs()

return {
	Content = json.encode(r)
}

