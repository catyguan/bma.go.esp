local r = {}

function r.start(T, params)
end

function r.testParse(T, params)

end

require("_:service/DataModel.class.lua")

local o = class.new("FormObject")
o.FIELD = {
	id={
		type="int",
		default=0,
		valid="not0",
		view=function(v)
			return v+123
		end
	},
	f1={
		type="int",
		default=100,
		valid="x",
		view=function(v)
			return v+123
		end
	},
	f2={
		type="x",
		valid="notEmpty",
		view=function(v)
			return v+123
		end
	}
}
o.Parse_x = function(v, dv)
	return types.string(v, dv)
end
o.Valid_x = function(field, v, dv)
	return false, "test"
end

-- Parse
local formData = {
	id="123",
	f2="abc"
}
o.Parse(formData)

out(strings.format("Parse Result = %v\n", o.prop))

-- Valid
local ok = o.Valid(data)
out(strings.format("Valid Result = %v, %v\n", ok, o.validData))

if false then
	data["x"] = "hello"
	local formView = o.BuildForm(data, validData)
	out(strings.format("Form View = %v\n", formView))
end

out("</pre>")

return r