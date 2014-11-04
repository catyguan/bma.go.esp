-- glf.mfc:service/DataModel.lua
local Class = class.define("DataModel")

function Class.ctor()
	self.prop = {}
end

function Class.Create(prop)
	self.prop = prop
end

function Class.Define(propName, name, val)
	local prop = self.prop[propName]
	if prop==nil then
		prop = {}
		self.prop[propName] = prop
	end
	prop[name] = val
end

function Class.Value(propName, name)
	local prop = self.prop[propName]
	return prop[name]
end

function Class.Parse(req)
	local r = {}
	for k, prop in self.prop do
		local v = req[k]
		local ptype = prop.type
		local dv = prop.default
		if types.name(ptype)=="function" then
			v = ptype(v)
		else
			if ptype~="" then
				local f = types[ptype]
				if f~=nil then
					v = f(v, dv)
				else
					error("unknow prop type '%s'", ptype)
				end
			end
		end
		if v==nil or v=="" then			
			if dv~=nil then v = dv end
		end
		if v~=nil then
			r[k] = v
		end
	end
	return r
end

function Class.Valid(data)
	local b = true
	local r = {}
	for k, prop in self.prop do
		local v = data[k]
		local vtype = prop.valid
		local msg = prop.valid_msg
		local ok = true
		if vtype==nil then
			continue
		end
		if types.name(vtype)=="function" then
			ok, msg = vtype(v)			
		else
			if vtype=="" then
				-- do nothing
			elseif vtype=="not.empty" then
				if v==nil or v=="" then
					ok = false
					if msg==nil then msg = "%s empty" end					
				end
			elseif vtype=="strlen" then
				local l = #v
				local l1 = prop.valid_max[0]
				local l2 = prop.valid_max[1]
				if l<l1 or l>l2 then
					ok = false
					if msg==nil then msg = "%s invalid length" end		
				end
			elseif vtype=="not.0" then
				if v==nil or v==0 then
					ok = false
				end
			elseif vtype==">.0" then
				if v==nil or v>0 then
					ok = false
				end
			end
		end
		if not ok then
			b = false
			if msg==nil then msg = "%s invalid" end
			r[k] = strings.format(msg, k)
		end
	end
	return b, r
end

function Class.BuildForm(data, validData)
	local r = {}
	for k, prop in self.prop do
		local v = data[k]
		local msg = ""
		if validData~=nil then
			msg = validData[k]
		end

		local f = prop.view
		if types.name(f)=="function" then
			v = f(v)
		end
		r[k] = {
			value=v,
			msg=msg
		}
	end
	for k, v in data do
		if r[k]==nil then
			r[k] = {value=v}
		end
	end
	return r
end