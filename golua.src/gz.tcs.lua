local r = []
go.exec("glf.mvc:testcase/suite.tcs.lua", {TCS=r})
go.exec("go.memserv:testcase/suite.tcs.lua", {TCS=r})
return r