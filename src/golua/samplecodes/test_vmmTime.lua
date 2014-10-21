local vr

-- local du = time.parseDuration("1.2m")
-- vr = du:Seconds()

-- local tm = time.date(2014,10,20,16,07,54)
-- local tm = time.now()
-- local tm = time.parse("2014-10-20 18:59:18")
-- local tm = time.unix(1)
local tm = time.now()
-- tm = tm:Add("1h")
-- tm = tm:AddDate(1,2,3)
-- vr = tm:Add("1h"):Before(tm)
-- vr = tm:Sub(tm)
-- vr = tm:Truncate("1h")
-- vr = tm:Unix()
-- vr = tm:YearDay()
vr = tm:ToMap()

return vr