local vr

-- vr = strings.contains("seafood", "foo")
-- vr = strings.hasPrefix("seafood", "sea")
-- vr = strings.hasSuffix("seafood", "ood")
-- vr = strings.index("seafood", "ood")
-- vr = strings.lastIndex("seafood", "ood")
-- vr = strings.replace("oink oink oink", "k", "ky", 2)
-- vr = strings.split("a,b,c", ",", 2)
-- vr = strings.toLower("ABVDX123")
-- vr = strings.toUpper("a,b,c")
-- vr = strings.trimSuffix(" a,b,c ","c ")
-- print(#vr)
-- vr = strings.substr("中文abcde", 1, 2)
-- vr = strings.format("a = %d, b=%v",1, "string")
vr = strings.parsef("a = ${2d}, b= ${1v}", "string", 1)
vr = strings.parsef("a = ${2d}, b= const", "string", 1)
vr = strings.parsef("a = ${2d}, b= ${1v}", "string")

return vr