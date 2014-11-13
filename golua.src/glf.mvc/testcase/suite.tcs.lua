local r = TCS
if r==nil then r = [] end
table.insert(r, {
	Script = "glf.mvc:testcase/DataModel",
	Title = "DataModelTC",
	Help = "测试mvc数据模型"
})
table.insert(r, {
	Script = "glf.mvc:testcase/Sample",
	Title = "SampleTC",
	Help = "测试用例编写范例"
})
return r