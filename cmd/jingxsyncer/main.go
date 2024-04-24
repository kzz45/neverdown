package main

import (
	"context"
	"flag"

	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	client "github.com/kzz45/neverdown/pkg/jingx/client-go"
	"github.com/kzz45/neverdown/pkg/jingx/client-go/syncer"
	"github.com/kzz45/neverdown/pkg/signals"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	"k8s.io/client-go/rest"
)

type args struct {
	address  *string
	username *string
	password *string
}

func (a args) validate() {
	if *a.address == "" {
		zaplogger.Sugar().Fatalf("empty address")
	}
	if *a.username == "" {
		zaplogger.Sugar().Fatalf("empty username")
	}
	if *a.password == "" {
		zaplogger.Sugar().Fatalf("empty password")
	}
}

func main() {
	zaplogger.Sugar().Info("jingx-syncer is preparing")
	args := args{
		address:  flag.String("address", "", "address"),
		username: flag.String("username", "", "username"),
		password: flag.String("password", "", "password"),
	}
	flag.Parse()
	args.validate()
	stopCh := signals.SetupSignalHandler()
	ctx, cancel := context.WithCancel(context.Background())
	handler := syncer.New(ctx, loadDiscoveryKubernetes())
	opt := &client.Option{
		Address:  *args.address,
		Username: *args.username,
		Password: *args.password,
	}
	zaplogger.Sugar().Info("jingx-syncer is starting")
	c := client.New(opt)
	c.Handle = handler.Handler
	c.DryRun()
	go c.Lister()
	zaplogger.Sugar().Info("jingx-syncer is running")
	<-stopCh
	zaplogger.Sugar().Info("jingx-syncer trigger shutdown")
	cancel()
	<-stopCh
	zaplogger.Sugar().Info("jingx-syncer shutdown gracefully")
}

func loadDiscoveryKubernetes() kubernetes.Interface {
	kubeconfig, err := rest.InDicoveryClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	cfg, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return cfg
}
