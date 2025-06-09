package svc

import (
	"mcp/applications/mcp/internal/config"
	"mcp/dao"
	"mcp/dao/query"
)

type ServiceContext struct {
	Config config.Config
	// MCPClient *mcpclient.MCPClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	query.SetDefault(dao.MustNewDB(c.DB))
	return &ServiceContext{
		Config: c,
		// MCPClient: mcpclient.NewMCPClient(c.MCPClientBaseURL),
	}
}
