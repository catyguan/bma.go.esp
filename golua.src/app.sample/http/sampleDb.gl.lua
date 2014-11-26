local form = httpserv.form()
local ajax = types.int(form, "ajax", 0)
local pdo = types.string(form, "pdo", "msg")

if ajax==1 then
	if pdo=="query" then
		local db = sql.open("mysql", "root:root@tcp(172.19.16.195:3306)/test_db")
		local rows = db.Query("SELECT * FROM smm_nodes")
		local msg = "<empty resultSet>"
		local rs
		if rows.Fetch(rs) then
			msg = strings.format("%v", rs)
		end
		return {
			Content = msg
		}
	end
end

httpserv.render("_:http/sampleDb.view.html",{})