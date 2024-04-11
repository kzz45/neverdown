package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/kzz45/neverdown/pkg/signals"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/kzz45/neverdown/pkg/env"

	"github.com/kzz45/neverdown/pkg/openx/app"
	kubecontrollerconfig "github.com/kzz45/neverdown/pkg/openx/app/config"
	"github.com/kzz45/neverdown/pkg/openx/app/options"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
)

func main() {
	var listenPort = flag.Int("listenPort", 8080, "listenPort")
	var kubeconfig = flag.String("kubeconfig", "", "kube configuration file path")
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *listenPort))
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	serverCert, serverKey, err := env.StaticCertKeyContent()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", serverCert, serverKey, "localhost")
	if err != nil {
		zaplogger.Sugar().Fatal("failed to generate self-signed certificate for loopback connection: %v", err)
	}
	stopCh := signals.SetupSignalHandler()
	opt := options.KubeControllerManagerOptions{
		Generic:                     kubecontrollerconfig.DefaultGenericControllerManagerConfiguration(),
		Master:                      "",
		Kubeconfig:                  *kubeconfig,
		ShowHiddenMetricsForVersion: "",
		ApiServerConfig: &server.Options{
			Listener: ln,
			SNICerts: []dynamiccertificates.SNICertKeyContentProvider{certProvider},
			StopChan: stopCh,
			Context:  context.Background(),
		},
	}
	kc, err := opt.Config(nil, nil)
	if err != nil {
		zaplogger.Sugar().Panic(err)
	}
	kc.Complete()
	go func() {
		if err := app.Run(kc.Complete(), stopCh); err != nil {
			zaplogger.Sugar().Panic(err)
		}
	}()
	zaplogger.Sugar().Info("openx is running")
	<-stopCh
	zaplogger.Sugar().Info("openx trigger shutdown")
	<-stopCh
	zaplogger.Sugar().Info("openx shutdown gracefully")
}
