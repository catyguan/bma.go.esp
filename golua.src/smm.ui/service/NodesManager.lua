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
	local rs = db.Query("select * from smm_nodes order where id = ?", id)
	local data
	local desc = {id="int",status="int"}
	if rs.Fetch(data, desc) then
		return data
	end
	return nil
end

function Class.Insert(data)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.create_time = time.unix()
	fv.modify_time = fv.create_time
	local id
	_, id = db.ExecInsert("smm_nodes", fv, true)
	return id
end

function Class.Update(data, id)
	local db = self.getDB()
	local fv = table.clone(data)
	fv.modify_time = time.unix()
	tj = {id=id}
	return db.ExecUpdate("smm_nodes", fv, tj)
end

function Class.Delete(id)
	local db = self.getDB()
	tj = {id=id}
	return db.ExecDelete("smm_nodes", tj)
end