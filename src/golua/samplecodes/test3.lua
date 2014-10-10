-- struct.lua
--- defining a struct constructor ---
local struct_mt = {
	-- instances can be created by calling the struct object
	__call = function(s,t)
		local obj = t or {}  -- pass it a table (or nothing)
		local fields = s._fields
		-- attempt to set a non-existent field in ctor?
		for k,v in pairs(obj) do
			if not fields[k] then
				s._error_nf(nil,k)
			end
		end
		-- fill in any default values if not supplied
		for k,v in pairs(fields) do
			if not obj[k] then
				obj[k] = v
			end
		end
		setmetatable(obj,s._mt)
		return obj
	end
}

-- creating a new struct triggered by struct.STRUCTNAME
struct = setmetatable({},{
	__index = function(tbl,sname)
		-- so we create a new struct object with a name
		local s = {_name = sname}
		-- and put the struct in the enclosing context
		_G[sname] = s
		-- the not-found error
		s._error_nf = function (tbl,key)
			error("field '"..key.."' is not in "..s._name)
		end
		-- reading or writing an undefined field of this struct is an error
		s._mt = {
			_name = s._name,
			__index = s._error_nf,
			__newindex = s._error_nf,
		}
		-- the struct has a ctor
		setmetatable(s,struct_mt)
		-- return a function that sets the struct's fields
		return function(t)
			s._fields = t
		end
	end
})