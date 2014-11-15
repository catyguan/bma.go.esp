local form = httpserv.form()
local n = types.string(form,"n","")

if n=="" then
	return {
		Content = json.encode({
			Message = "param 'n' empty"
		})
	}
end

local suite = go.exec(n..".tcs.lua", {})
local r = {
	Result = suite
}
return {
	Content = json.encode(r)
}

