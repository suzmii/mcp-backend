package svc

import (
	"mcp/applications/auth/internal/config"
	"mcp/core/middleware/cache"
)

type ServiceContext struct {
	Config config.Config
	Auth   *cache.Auth
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Auth:   cache.NewAuth(c.Cache),
	}
}
