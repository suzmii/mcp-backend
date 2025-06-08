package client

import (
	"mcp/applications/auth/authservice"
	"mcp/applications/mcp/mcpservice"
	"mcp/applications/user/userservice"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var Auth authservice.AuthService

var User userservice.UserService

var MCP mcpservice.McpService

func defaultConfig(key string) zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   key,
		},
	}
}
func init() {
	Auth = authservice.NewAuthService(zrpc.MustNewClient(defaultConfig("auth.rpc")))
	User = userservice.NewUserService(zrpc.MustNewClient(defaultConfig("user.rpc")))
	MCP = mcpservice.NewMcpService(zrpc.MustNewClient(defaultConfig("mcp.rpc")))
}
