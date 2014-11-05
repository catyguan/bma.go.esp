local Class = class.define("NodesManager")

function Class.getDB()
	if not self.db then
		local sdb = go.new("SmartDB")
		self.db = sdb.Select("smm_nodes")
	end
	return self.db
end

function Class.List()
	local db = self.getDB()
	local rs = db.Query("select * from smm_nodes order by host_name")
	local nodes = []
	local data
	local desc = {id="int",status="int"}
	while rs.Fetch(data, desc) do
		table.insert(nodes, data)
	end
	return nodes
end

function Class.Get(id)
	local db = self.getDB()
	local rs = db.Query("select * from smm_nodes where id = ?", id)
	local data
	local desc = {id="int",status="int"}
	if rs.Fetch(data, desc) then
		return data
	end
	return nil
end

function Class.Valid(data)
	local ok, vdata = self.DM_NODE.Valid(data)
	if not ok then
		return ok, vdata
	end
	-- other valid
	return ok, vdata
end

function Class.Insert(data)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.create_time = time.now().Unix()
	fv.modify_time = fv.create_time
	local id
	_, id = db.ExecInsert("smm_nodes", fv, true)
	return id
end

function Class.Update(data, id)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.modify_time = time.now().Unix()
	tj = {id=id}
	return db.ExecUpdate("smm_nodes", fv, tj)
end

function Class.Delete(id)
	local db = self.getDB()
	tj = {id=id}
	return db.ExecDelete("smm_nodes", tj)
end