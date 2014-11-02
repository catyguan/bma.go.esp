local vr

config.set("abc", 123)
-- vr = config.get("abc")
-- vr = config.get("!global.Debug")
vr = config.parse("my${abc}_${!global.Debug}")

return vr