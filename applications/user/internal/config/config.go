package config

import (
	"mcp/dao"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	AuthService zrpc.RpcClientConf
	DB          dao.DBConfig
}
