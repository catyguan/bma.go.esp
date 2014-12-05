print(request.String)

tname = ""
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