require("_:mysqlcore.lua")

local M2 = mysql_m2table(gom, "HelloParams")

local M = {
	tableName = "rank_mini_rich_weekly",
	engine="InnoDB",
	charset="utf8",
	comment="mini本周富豪榜",
	fields = [
		{
			name="ranking",
			type="int(10)",
			notNull=true,
			comment="排名名次"
		},
		{
			name="uid",
			type="int(20)",
			notNull=true,
			comment="用户uid"
		},
		{
			name="nick",
			type="varchar(128)",
			notNull=true,
			comment="用户昵称"
		},
		{
			name="pay",
			type="int(20)",
			notNull=true,
			comment="身价"
		},
		{
			name="create_time",
			type="int(10)",
			notNull=true,
			comment="生成时间"
		}
	],
	PrimaryKey = {
		fields=["ranking","uid"]
	},
	Indexs = [
		{
			name="activityNewsStatusIndex",
			fields=["status"]
		}
	]
}
local vo = {
	M = M2
}
gomf.render("mysql:mysql.tpl",vo,"test.out")