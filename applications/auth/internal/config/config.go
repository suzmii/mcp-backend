package config

import (
	"mcp/core/middleware/cache"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache      cache.Config
	AccessExp  time.Duration
	RefreshExp time.Duration
}
