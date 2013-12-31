namespace go app.cacheserver.thriftpoint

struct TCacheRequest {
	1: string groupName 		// 缓冲分组名称
	2: string key				// 查询的键
	3: optional bool trace		// 是否返回调试信息
	4: optional bool notLoad   	// 不命中的时候，是否不加载数据
	5: optional i32 Timeout   	// 请求超时时间，毫秒ms, <=0 表示无超时	
}

struct TCacheResult {
	1: bool done,						// 是否有结果
	2: binary value,					// 查询的值，二进制数据
	3: i32 length,						// 查询的值的长度
	4: optional string error,			// 错误信息，如果有的话
	5: optional list<string> traces,	// 调试信息
}

service TCacheServer {

	TCacheResult cacheServerGet(1:TCacheRequest req, 2:map<string, string> options),

	void cacheServerLoad(1:string groupName, 2:string key),

	void cacheServerPut(1:string groupName, 2:string key, 3:binary value, 4:i32 length),

	void cacheServerErase(1:string groupName, 2:string key),
	
}