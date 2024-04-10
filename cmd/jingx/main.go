package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/kzz45/neverdown/pkg/apiserver/server"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/env"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	jingxserver "github.com/kzz45/neverdown/pkg/jingx/server"

	"k8s.io/client-go/rest"

	"k8s.io/apiserver/pkg/server/dynamiccertificates"
)

const (
	StaticCertFile = "TLS_OPTION_CERT_FILE"
	StaticKeyFile  = "TLS_OPTION_KEY_FILE"
)

func main() {
	var listenPort = flag.Int("listenPort", 8083, "listenPort")
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *listenPort))
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	defer ln.Close()
	serverCert, serverKey, err := env.StaticCertKeyContent()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", serverCert, serverKey, "localhost")
	if err != nil {
		zaplogger.Sugar().Fatal("failed to generate self-signed certificate for loopback connection: %v", err)
	}
	opt := &jingxserver.Options{
		Listener:           ln,
		SNICerts:           []dynamiccertificates.SNICertKeyContentProvider{certProvider},
		StopChan:           server.SetupSignalHandler(),
		Context:            context.Background(),
		AuthorityClientSet: loadAuthorityKubernetes(),
		Api:                aggregator.New(context.Background(), loadDiscoveryKubernetes(), ""),
	}
	stopConfirm, listenConfirm, err := jingxserver.Run(opt)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	// metrics
	// go metrics.RunExporter()
	<-stopConfirm
	_ = listenConfirm
	zaplogger.Sugar().Info("gracefully shutdown jingx-apiserver")
}

func loadDiscoveryKubernetes() kubernetes.Interface {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	cfg, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return cfg
	// certFile, keyFile := os.Getenv(StaticCertFile), os.Getenv(StaticKeyFile)
	// clientset, err := kubernetes.NewForConfig(&rest.Config{
	// 	Host: "127.0.0.1:9443",
	// 	TLSClientConfig: rest.TLSClientConfig{
	// 		Insecure: true,
	// 		CertFile: certFile,
	// 		KeyFile:  keyFile,
	// 	},
	// })
	// if err != nil {
	// 	zaplogger.Sugar().Fatal(err)
	// }
	// return clientset
}

func loadAuthorityKubernetes() kubernetes.Interface {
	kubeconfig, err := rest.InAuthorityClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	cfg, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return cfg
	// certFile, keyFile := os.Getenv(StaticCertFile), os.Getenv(StaticKeyFile)
	// clientset, err := kubernetes.NewForConfig(&rest.Config{
	// 	Host: "127.0.0.1:9443",
	// 	TLSClientConfig: rest.TLSClientConfig{
	// 		Insecure: true,
	// 		CertFile: certFile,
	// 		KeyFile:  keyFile,
	// 	},
	// })
	// if err != nil {
	// 	zaplogger.Sugar().Fatal(err)
	// }
	// return clientset
}
