package svc

import (
	"mcp/applications/mcp/internal/config"
)

type ServiceContext struct {
	Config config.Config
	// MCPClient *mcpclient.MCPClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// MCPClient: mcpclient.NewMCPClient(c.MCPClientBaseURL),
	}
}
