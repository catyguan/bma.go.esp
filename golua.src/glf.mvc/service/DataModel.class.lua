-- glf.mvc:service/DataModel.class.lua
-- go.deving()
local VOClass = class.define("ViewObject")

function VOClass.IFF(name, str1, str2)
	local v = self[name]
	local str
	if v==nil or v=="" then
		str = str2
	else
		str = str1
	end
	if str==nil or str=="" then return "" end
	return strings.parsef(str, v)
end

local FOClass = class.define("FormObject")

function FOClass.Value(name)
	local v
	if self.data~=nil then
		v = self.data[name]
	end
	if v~=nil then
		return v
	end
	local field = self.FIELD[name]
	if field~=nil then
		return field.default
	end
	return nil
end

function FOClass.CloneData()
	return table.clone(self.data)
end

function FOClass.Bind(data)
	self.data = data
end

function FOClass.Parse(req)
	local r = {}
	for k, field in self.FIELD do
		local v = req[k]
		local ptype = field.type
		local dv = field.default
		if types.name(ptype)=="function" then
			v = ptype(v)
		else
			if ptype~="" then
				local f = types[ptype]
				if f==nil then
					f = self["Parse_"..ptype]
				end
				if f~=nil then					
					v = f(v, dv)
				else					
					error("unknow parse type '%s'", ptype)				
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
	self.data = r
end

function FOClass.Valid()
	local data = self.data
	if data==nil then
		data = {}
	end
	local b = true
	local r = {}
	for k, field in self.FIELD do
		local v = data[k]
		local vtype = field.valid
		local msg = nil
		local ok = true
		if vtype==nil or vtype=="" then
			continue
		end		
		if types.name(vtype)=="function" then			
			ok, msg = vtype(field, v, msg)
		else
			local f = self["Valid_"..vtype]			
			if f~=nil then
				ok, msg = f(field, v, msg)
			else
				error("unknow valid type '%s'", vtype)
			end
		end
		if not ok then
			b = false
			if msg==nil then msg = "${1s} invalid" end
			r[k] = strings.parsef(msg, k, v)
		end
	end
	self.validData = r
	return b
end

function FOClass.Valid_notEmpty(field, v, msg)
	if v==nil or v=="" then		
		if msg==nil then msg = "${1s} empty" end
		return false, msg
	end
	return true
end

function FOClass.Valid_strlen(field, v, msg)
	local l = #v
	local l1 = field.valid_max[0]
	local l2 = field.valid_max[1]
	if l<l1 or l>l2 then		
		if msg==nil then msg = "${1s} invalid length" end		
		return false, msg
	end
	return true
end

function FOClass.Valid_not0(field, v, msg)
	if v==nil or v==0 then
		return false
	end
	return true
end

function FOClass.Valid_gt0(field, v, msg)
	if v==nil or v<=0 then
		return false
	end
	return true
end

function FOClass.Error(name, str1, str2)
	local v
	if self.validData~=nil then
		v = self.validData[name]
	end
	local str
	if v==nil or v=="" then
		str = str2
	else
		str = str1
	end
	if str==nil or str=="" then return "" end
	return strings.parsef(str, name, v)
end