local TC = {}

TC.helpServBase = {
	Params = [
		{
			Name="n",
			Title="MemGo Name",
			Value="app"
		},
		{
			Name="key",
			Title="Key Name",
			Value="key1"
		},
		{
			Name="val",
			Title="Value",
			Value="world"
		},
		{
			Name="act",
			Title="Get(1) Set(2) Remove(3)",
			Value="1"
		}
	]
}
function TC.testServBase(T, params)
	local n,key,val,act
	n = types.string(params, "n", "")
	key = types.string(params, "key", "")
	val = types.string(params, "val", "")
	act = types.int(params, "act", 1)
	T.w("param n=%v, key=%v, val=%v, act=%v", n, key, val, act)

	local o = go.new("MemServ")
	local mg = o.Get(n)

	if act==2 then
		mg.Set(key, val)
	elseif act==3 then
		mg.Remove(key)
	else
		local v = mg.Get(key)
		T.w("val = %v", v)
	end

	local c, sz = mg.Size()
	T.w("%s -> %T, %d, %d", n, mg, c, sz)
end

return TC