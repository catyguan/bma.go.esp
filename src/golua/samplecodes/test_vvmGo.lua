local vr
-- go module
-- local function f()
-- 	print("run in go")
-- end
-- go.run(f, "test")
-- go.defer(function()
-- 	print("run in defer")
-- end)

-- local ch = go.chan(1)
-- go.defer(ch)
-- local ch2 = go.chan(1)
-- go.deferClose(ch2)

-- go.write(ch2, 123)
-- print(go.read([ch, ch2], 50))
-- vr = go.read(ch, 50)
-- go.close(ch)

-- local mux = go.mutex(true)
-- mux:Lock()
-- go.defer(mux)
-- local rmux = mux:RLocker()
-- rmux:Lock()
-- vr = rmux:Sync(function()
-- 	return 1
-- end)

go.sleep(100)

return vr