require("base:helpers.lua")
require("base:strings.lua")

function mysql_m2table(gom, name, helpers)
	-- body
	local M = {}
	local v
	local obj = gom.GetStruct(name)
	if obj==nil then
		error("invalid struct '%s'", name)
	end
	tableNameF = helper_get(helpers, "tableName", mysql_tableName)
	M.tableName = tableNameF(obj)
	M.engine = helper_annotation(obj, "mysql.engine", "InnoDB")
	M.charset = helper_annotation(obj, "mysql.charset", "utf8")
	M.comment = helper_annotation(obj, "comment", "")
	M.fields = []
	local fs = obj.Fields
	local fieldNameF = helper_get(helpers, "fieldName", mysql_fieldName)
	local typeNameF = helper_get(helpers, "fieldType", mysql_fieldType)
	for f in fs do		
		local mf = {}
		mf.name=fieldNameF(f)
		mf.type=typeNameF(f)
		mf.notNull=true
		mf.comment=helper_annotation(f, "comment", "")
		table.insert(M.fields, mf)
	end
	print("----->", M.engine)
	return M
end

function mysql_tableName(obj)
	local n
	n = obj.Annotation("mysql.name")
	if n~=nil then
		return n
	end
	n = obj.Name
	return "tbl_" .. string_underscore(n)
end

function mysql_fieldName(obj)
	local n
	n = obj.Annotation("mysql.name")
	if n~=nil then
		return n
	end
	n = obj.Name	
	return string_underscore(n)
end

local TYPES = {
	["bool"] = {name="int(1)",length=false},
	["string"] = {name="varchar",length=true}
}

function mysql_fieldType(obj)
	closure(TYPES)
	local r
	local hasl = true
	local typ = obj.Type
	local l = obj.Annotation("mysql.length")	
	r = typ
	xtyp = TYPES[typ]
	if xtyp~=nil then
		r = xtyp.name
		hasl = xtyp.length
	end
	if l~=nil and hasl then
		r = r .. "(" .. l .. ")"
	end
	return r
end