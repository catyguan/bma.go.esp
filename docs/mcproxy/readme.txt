h1. MCProxy

h2. 功能说明

目的：同时对多个MC做读写，实现多个MC的数据的同步更新

支持标准命令:
version
quit
set,add,delete -- 同时更新访问所有MC，所有成功响应成功，否则错误
get -- 同时get访问所有MC，响应最新获得的结果，丢弃其他结果

支持扩展命令:
reload -- 重载配置文件
getall xxx -- 同时get访问所有MC，用于检测各个MC的某个KEY值情况

h2. 部署说明

h3. 安装Go

http://golang.org/doc/install

h3. 安装第三方库

# 安装 go, git
# 配置公共 GOPATH

h3. 安装MCProxy

# 建立应用目录，如 /data/webapps/mcproxy
# 建立源代码目录，如 /xxxx/
# git clone https://github.com/catyguan/bma.go.esp.git
# git checkout 版本号
# bma.go.esp/deploy/build.sh app/mcproxy /data/webapps/mcproxy mcproxy-config.json

h3. 配置启动

启动配置
# vi config/mcproxy-config.json
# ./mcproxy > mcproxy.log &

配置项：
global.GOMAXPROCS -- go 携程数
logger.RootLevel -- 日志信息等级
mcPoint.Address -- 代理地址，如：127.0.0.1:11213
mcPoint.Port -- 代理端口
service.PoolMax -- 远程MC的连接池大小，缺省10
service.Remotes -- 目标MC地址，如 ["172.19.16.97:11211", "172.19.16.195:11211"]

h4. 关闭

* kill xxxx