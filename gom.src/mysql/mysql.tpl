CREATE TABLE `<?=M.tableName?>` (
<?golua
local idx = 0
for idx, f in M.fields do
	local s = ""
	local typ = f.type
	local notNull = ""
	if f.notNull then
		notNull = " NOT NULL"
	end
	local defv = ""
	if not types.isEmpty(f.default) then
		defv = " DEFALUT "..f.default
	end
	local auto = ""
	if f.auto then
		auto = " AUTO_INCREMENT"
	end
	local comment = ""
	if not types.isEmpty(f.comment) then
		comment = " COMMENT '"..f.comment.."'"
	end
	s = s .. strings.format("\t`%s` %s%s%s%s", f.name, typ, notNull, auto, comment)
	if idx>0 then
		out(",\n")
	end
	idx = idx + 1
	out(s)
end
local ko = M.PrimaryKey
if not types.isEmpty(ko) then
	s = "PRIMARY KEY"
	if not types.isEmpty(ko.name) then
		s = s .. " "..ko.name
	end
	s = s .. "("
	for i, fn in ko.fields do
		if i>0 then
			s = s .. ","
		end
		s = s .. "`" .. fn .. "`"
	end
	s = s .. ")"
	if idx>0 then
		out(",\n")
	end
	idx = idx + 1
	out("\t")
	out(s)
end
for i, ko in M.Indexs do
	s = "INDEX"	
	if ko.unique then
		s = "UNIQUE"
	end
	if not types.isEmpty(ko.name) then
		s = s .. " "..ko.name
	end
	s = s .. "("
	for i, fn in ko.fields do
		if i>0 then
			s = s .. ","
		end
		s = s .. "`" .. fn .. "`"
	end
	s = s .. ")"
	if idx>0 then
		out(",\n")
	end
	idx = idx + 1
	out("\t")
	out(s)
end
?>
) ENGINE=<?=M.engine?> DEFAULT CHARSET=<?=M.charset?> COMMENT='<?=M.comment?>'