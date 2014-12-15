print(request.String)

tname = "opendao"
if request.Name=="get_sid_by_asid" then
	tname = "opendao"
end

if tname=="" then
	error("unknow request(%s)", request.Name)
end

return {
	Action="forward",
	Target=tname
}