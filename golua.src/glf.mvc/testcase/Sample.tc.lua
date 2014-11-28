local TC = {}

--[[@title:加法]]
--[[@help:测试加法实现]]
--[[@desc:
传入参数a, b
算出a + b结果
]]
--[[@params:{"Name":"a"}]]
--[[@params:{"Name":"b","Title":"B", "Value":2}]]
function TC.testAdd(T, params)
	local a, b
	a = types.int(params, "a", 0)
	b = types.int(params, "b", 0)
	T.w("a + b = %d", a+b)
end

function TC.testSub(T, params)
	local a, b
	a = types.int(params, "a", 0)
	b = types.int(params, "b", 0)
	T.error("a - b = %d", a-b)
end

return TC