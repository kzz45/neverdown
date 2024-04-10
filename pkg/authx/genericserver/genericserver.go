package genericserver

import (
	"context"

	"github.com/kzz45/neverdown/pkg/authx/grpc/aggregator"
	grpcserver "github.com/kzz45/neverdown/pkg/authx/grpc/server"
	httpserver "github.com/kzz45/neverdown/pkg/authx/http/server"
	"github.com/kzz45/neverdown/pkg/authx/rbac/handler"
)

type Options struct {
	HttpServerOptions *httpserver.Options
	GrpcServerOptions *grpcserver.GrpcApiOption
}

type GenericServer struct {
	httpServer *httpserver.Server
	grpcServer *grpcserver.Server
}

func New(opts *Options) *GenericServer {
	handlers := handler.NewHandler(context.Background())
	opts.HttpServerOptions.Handler = handlers

	opts.GrpcServerOptions.Aggregator = aggregator.New(handlers.AdminApp())

	gs := &GenericServer{
		httpServer: httpserver.Init(opts.HttpServerOptions),
		grpcServer: grpcserver.Init(opts.GrpcServerOptions),
	}
	return gs
}

func (s *GenericServer) DryRun() (<-chan struct{}, <-chan struct{}, error) {
	s.grpcServer.DryRun()
	return s.httpServer.DryRun()
}

func (s *GenericServer) Shutdown() {
	s.httpServer.Shutdown()
	s.grpcServer.Shutdown()
}
