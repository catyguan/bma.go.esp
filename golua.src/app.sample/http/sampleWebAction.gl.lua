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
end

httpserv.render("_:http/sampleWebAction.view.htm",{})