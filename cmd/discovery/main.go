package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/kzz45/neverdown/pkg/apiserver/server"
	"github.com/kzz45/neverdown/pkg/env"
	"github.com/kzz45/neverdown/pkg/kubernetes/controlplane"

	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/klog/v2"
)

const (
	EtcdEndpoints     = "ETCD_ENDPOINTS"
	EtcdPrefix        = "ETCD_PREFIX"
	EtcdDefaultPrefix = "/registry"
)

func main() {
	var listenPort = flag.Int("listenPort", 9443, "listenPort")
	klog.InitFlags(nil)
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *listenPort))
	if err != nil {
		klog.Fatal(err)
	}
	defer ln.Close()
	serverCert, serverKey, err := env.StaticCertKeyContent()
	if err != nil {
		klog.Fatal(err)
	}

	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", serverCert, serverKey, "localhost")
	if err != nil {
		klog.Fatalf("failed to generate self-signed certificate for loopback connection: %v", err)
	}
	s := server.SecureServingInfo{
		Listener: ln,
		SNICerts: []dynamiccertificates.SNICertKeyContentProvider{certProvider},
	}

	c := &storagebackend.Config{
		Type: "etcd3",
		Transport: storagebackend.TransportConfig{
			ServerList: etcdEndpoints(),
		},
		Prefix: etcdPrefix(),
	}

	instance := controlplane.New(c)

	klog.Info(instance.GenericAPIServer.Handler.ListedPaths())

	sig := server.SetupSignalHandler()
	stopConfirm, _, err := s.Serve(instance.GenericAPIServer.Handler, time.Second*5, sig)
	if err != nil {
		klog.Fatal(err)
	}
	<-stopConfirm
	// _ = listenerConfirm
	klog.Infof("gracefully shutdown discovery-apiserver")
}

func etcdEndpoints() []string {
	endpoints := os.Getenv(EtcdEndpoints)
	if endpoints == "" {
		klog.Fatal("no EtcdEndpoints get from the environment")
	}
	return strings.Split(endpoints, ",")
}

func etcdPrefix() string {
	p := os.Getenv(EtcdPrefix)
	if p == "" {
		return EtcdDefaultPrefix
	}
	return p
}
