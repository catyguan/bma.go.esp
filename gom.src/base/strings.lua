function string_cap(s)
	if #s>0 then
		local fc = strings.sub(s,0,1)
		return strings.toUpper(fc) .. strings.sub(s, 1)
	end
	return s
end

function rune_iscap(rune)
	if rune>=65 and rune<=90 then
		return true
	end
	return false
end

function rune_toLower(rune)
	if rune_iscap(rune) then
		return rune - 65 + 97
	end
	return rune
end

function string_underscore(s)
	local buf1 = go.new("ByteBuffer")
	local buf2 = go.new("ByteBuffer")
	local last = false
	local first = true
	buf1.Write(s)
	while true do
		if buf1.Len()==0 then
			break
		end
		local c = buf1.ReadRune()
		if first then
			first = false
		else
			if rune_iscap(c) then
				if not last then
					buf2.WriteRune(95)
				end		
				last = true
			else
				last = false
			end
		end
		buf2.WriteRune(rune_toLower(c))
	end
	return buf2.String()
end