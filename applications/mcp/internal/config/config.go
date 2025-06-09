package config

import (
	"mcp/dao"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// MCPClientBaseURL string
	DB dao.DBConfig
}
