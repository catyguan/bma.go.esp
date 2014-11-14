require("glf.mvc:service/DataModel.class.lua")

local Class = class.define("NodeForm", "FormObject")

Class.FIELD = {
	id={
		type="int",
		default=0
	},
	name={
		type="string",
		valid="notEmpty"
	},
	host_name={
		type="string",
		valid="notEmpty"
	},
	api_url={
		type="string",
		valid="notEmpty"
	},
	code={
		type="string",
		default=""
	},
	type={
		type="int",
		default=1
	},
	remark={
		type="string",
		default=""
	}
}