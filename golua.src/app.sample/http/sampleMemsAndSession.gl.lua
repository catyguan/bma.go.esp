local form = httpserv.form()
local ajax = types.int(form, "ajax", 0)
local pdo = types.string(form, "pdo", "msg")

if ajax==1 then
	if pdo=="getmems" then
		local o = go.new("MemServ")
		local mg = o.Get("app") -- 应用关联容器
		local v = mg.Get("test")
		return {
			Content = strings.format("GetMem 'test'-> [%v]", v)
		}
	end
	if pdo=="setmems" then
		local o = go.new("MemServ")
		local mg = o.Get("app") -- 应用关联容器
		mg.Set("test", "hello world")
		return {
			Content = "OK"
		}
	end
	if pdo=="getsess" then
		local sess = go.new("HttpSession")
		local v = sess.Get("test")
		return {
			Content = strings.format("GetSess 'test'-> [%v]", v)
		}
	end
	if pdo=="setsess" then
		local sess = go.new("HttpSession")
		local v = sess.Set("test","hello session")
		return {
			Content = "OK"
		}
	end
end

httpserv.render("_:http/sampleMemsAndSession.view.html",{})