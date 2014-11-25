local form = httpserv.form()
local method = types.string(form, "m", "say")
local content = types.string(form, "p", "")
local params
if types.isEmpty(content) then
	params = ["world"]
else
	params = json.decode(content)
end

if method=="say" then
	print(strings.format("hello '%s'", params[0]))
	r = json.encode(true)
	return {
		Content = r
	}
end

return {Status=502}