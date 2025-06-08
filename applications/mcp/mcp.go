package main

import (
	"flag"
	"fmt"

	"mcp/applications/mcp/internal/config"
	"mcp/applications/mcp/internal/server"
	"mcp/applications/mcp/internal/svc"
	"mcp/applications/mcp/pb/mcp"
	"mcp/core/log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mcp.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		mcp.RegisterMcpServiceServer(grpcServer, server.NewMcpServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	log.InitLogx()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
