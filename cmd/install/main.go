package main

import (
	"context"
	"flag"

	"github.com/kzz45/neverdown/pkg/install"
	"github.com/kzz45/neverdown/pkg/openx/clientbuilder"
	"github.com/kzz45/neverdown/pkg/signals"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig = flag.String("kubeconfig", "", "kube configuration file path")
	flag.Parse()
	kc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	stop := signals.SetupSignalHandler()
	ctx, cancel := context.WithCancel(context.Background())

	b := install.NewInstaller(ctx,
		&install.Options{
			ClientBuilder: clientbuilder.SimpleControllerClientBuilder{
				ClientConfig: kc,
			},
			DockerSecrets:     install.DefaultDockerSecrets(),
			CertificateOption: install.DefaultCertificateOption(),
		},
	)
	zaplogger.Sugar().Infof("installer is running")
	b.Run()
	<-stop
	zaplogger.Sugar().Infof("installer is cancelling")
	cancel()
	<-stop
	zaplogger.Sugar().Infof("installer shutdown")
}
