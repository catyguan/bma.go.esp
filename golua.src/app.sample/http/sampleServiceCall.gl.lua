local form = httpserv.form()
local ajax = types.int(form, "ajax", 0)
local pdo = types.string(form, "do", "say")

if ajax==1 then
	if pdo=="say" then
		local scs = go.new("ServiceCall")
		local sc = scs.Assert("sample", 1000)
		local r = sc.CallTimeout("say", 1000, "Kitty")
		local msg = "NotGood"
		if r then msg = "OK" end
		return {
			Content = msg
		}
	end
	if pdo=="bind" then
		local f = function(methodName, params)
			print(strings.format("Calling %s, %v", methodName, params))
			return 123
		end
		local scs = go.new("ServiceCall")
		scs.Bind("sample2", f, true)

		local sc = scs.Assert("sample2", 1000)
		local r = sc.CallTimeout("say", 1000, "Kitty")
		return {
			Content = strings.format("%v", r)
		}	
	end
end

httpserv.render("_:http/sampleServiceCall.view.htm",{})