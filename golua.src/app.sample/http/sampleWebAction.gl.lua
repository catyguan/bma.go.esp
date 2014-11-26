local form = httpserv.form()
local ajax = types.int(form, "ajax", 0)
local pdo = types.string(form, "pdo", "msg")

if ajax==1 then
	if pdo=="msg" then
		local msg = types.string(form, "msg", "<empty message>")
		return {
			Content = "Server Message ["..msg.."]"
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

httpserv.render("_:http/sampleWebAction.view.html",{})