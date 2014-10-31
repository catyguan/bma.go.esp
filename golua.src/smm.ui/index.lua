local form = httpserv.form()
local a, b, c

a = types.int(form,"a",0)
b = types.int(form,"b",0)

httpserv.render("_:index.view",{a=a,b=b})