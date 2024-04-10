package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/kzz45/neverdown/pkg/authx/genericserver"
	grpcserver "github.com/kzz45/neverdown/pkg/authx/grpc/server"
	httpserver "github.com/kzz45/neverdown/pkg/authx/http/server"

	"github.com/kzz45/neverdown/pkg/env"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	"k8s.io/apiserver/pkg/server"

	"k8s.io/apiserver/pkg/server/dynamiccertificates"
)

func main() {
	var httpListenPort = flag.Int("httpListenPort", 8087, "httpListenPort")
	var grpcListenPort = flag.Int("grpcListenPort", 8088, "grpcListenPort")
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *httpListenPort))
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	defer ln.Close()

	serverCert, serverKey, err := env.StaticCertKeyContent()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}

	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", serverCert, serverKey)
	if err != nil {
		zaplogger.Sugar().Fatal("failed to generate self-signed certificate for loopback connection: %v", err)
	}

	certs, err := env.StaticServerCerts()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	secret, err := env.GrpcAuthenticationSecret()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}

	opts := &genericserver.Options{
		HttpServerOptions: &httpserver.Options{
			Listener: ln,
			SNICerts: []dynamiccertificates.SNICertKeyContentProvider{certProvider},
			StopChan: server.SetupSignalHandler(),
			Context:  context.Background(),
		},
		GrpcServerOptions: &grpcserver.GrpcApiOption{
			Certificate: &certs,
			ListenPort:  int32(*grpcListenPort),
			Secret:      secret,
			Aggregator:  nil,
		},
	}
	genericServer := genericserver.New(opts)
	stopConfirm, listenConfirm, err := genericServer.DryRun()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}

	// metrics
	// go metrics.RunExporter()

	<-stopConfirm
	_ = listenConfirm
	genericServer.Shutdown()
	zaplogger.Sugar().Info("gracefully shutdown authx-apiserver")
}
