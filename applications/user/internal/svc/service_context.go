package svc

import (
	"mcp/applications/auth/authservice"
	"mcp/applications/user/internal/config"
	"mcp/dao"
	"mcp/dao/query"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	AuthService authservice.AuthService
}

func NewServiceContext(c config.Config) *ServiceContext {
	query.SetDefault(dao.MustNewDB(c.DB))

	return &ServiceContext{
		Config:      c,
		AuthService: authservice.NewAuthService(zrpc.MustNewClient(c.AuthService)),
	}
}
