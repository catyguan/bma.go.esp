local Class = class.define("NodesManager")

Class.M = go.enableSafe({}, true)

function Class.getDB()
	if not self.db then
		local sdb = go.new("SmartDB")
		self.db = sdb.Select("smm_nodes")
	end
	return self.db
end

function Class.List(refresh)
	local vb = types.bool(refresh, false)
	if vb then self.M.Nodes = nil end

	if self.M.Nodes~=nil then
		return self.M.Nodes
	end

	local db = self.getDB()
	local rs = db.Query("select * from smm_nodes order by host_name")
	local nodes = []
	local data
	local desc = {id="int",status="int"}
	while rs.Fetch(data, desc) do
		table.insert(nodes, data)
	end
	self.M.Nodes = nodes
	return nodes
end

function Class.Get(id)
	if self.M.Nodes == nil then
		self.List(true)
	end

	local nodes = self.M.Nodes
	if nodes~=nil then
		for _, node in nodes do
			if node.id == id then
				return node
			end
		end
		return nil
	end
	-- local db = self.getDB()
	-- local rs = db.Query("select * from smm_nodes where id = ?", id)
	-- local data
	-- local desc = {id="int",status="int"}
	-- if rs.Fetch(data, desc) then
	-- 	return data
	-- end
	return nil
end

function Class.Insert(data)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.create_time = time.now().Unix()
	fv.modify_time = fv.create_time
	local id
	_, id = db.ExecInsert("smm_nodes", fv, true)
	self.M.Nodes = nil
	return id
end

function Class.Update(data, id)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.modify_time = time.now().Unix()
	tj = {id=id}
	self.M.Nodes = nil
	return db.ExecUpdate("smm_nodes", fv, tj)
end

function Class.Delete(id)
	local db = self.getDB()
	tj = {id=id}
	self.M.Nodes = nil
	return db.ExecDelete("smm_nodes", tj)
end