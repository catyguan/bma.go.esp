function helper_get(helpers, name, defv)
	if helpers==nil then
		return defv
	end
	local v = helpers[name]
	if v==nil then
		return defv
	end
	return v
end

function helper_annotation(obj, name, defv)
	local v
	v = obj.Annotation(name)
	if v==nil then
		return defv
	end
	return v
end