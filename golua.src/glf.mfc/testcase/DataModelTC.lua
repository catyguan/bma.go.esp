require("_:service/DataModel.lua")

local out = httpserv.write
out("<pre>\n")

local o = class.new("DataModel")
o.Create({
	id={
		type="int",
		default=0,
		valid="not.0",
		view=function(v)
			return v+123
		end
	}
})

local formData = {
	
}
local data = o.Parse(formData)

out(strings.format("Parse Result = %v\n", data))

local ok, validData = o.Valid(data)
out(strings.format("Valid Result = %v, %v\n", ok, validData))

data["x"] = "hello"
local formView = o.BuildForm(data, validData)
out(strings.format("Form View = %v\n", formView))

out("</pre>")