-- glf.testing:T.class.lua
local Class = class.define("glf.testing.T")

function Class.add(msg, err)
	if self.logs==nil then self.logs = [] end
	table.insert(self.logs, {msg=msg, err=err})
end

function Class.GetLogs()
	return self.logs
end

function Class.w(s, ...)
	local msg = strings.formata(s, ...)
	self.add(msg, false)
	go.debug("T", msg)
end

function Class.error(s, ...)
	local msg = strings.formata(s, ...)
	self.add(msg, true)
	go.warn("T", msg)
end
