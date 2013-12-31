namespace go testapp.thrift.protocol
namespace java demo.rpc

service RpcService {
	    list<string> funCall(1:i64 callTime, 2:string funCode, 3:map<string, string> paramMap),	 
}