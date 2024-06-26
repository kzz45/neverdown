package install

import (
	"context"

	openxclientset "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/openx/clientbuilder"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	cmap "github.com/orcaman/concurrent-map"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	NamespaceKubeDiscovery = "kube-discovery"
	NamespaceKubeApi       = "kube-api"
	NamespaceKubeAuthx     = "kube-authx"
)

type Installer interface {
	Run()
}

type Options struct {
	ClientBuilder             clientbuilder.ControllerClientBuilder
	DockerSecrets             []DockerSecret
	JingxImage                *JingxImage
	CertificateOption         *CertificateOption
	ServiceLoadBalancerOption *ServiceLoadBalancerOption
	NativeAppHost             *NativeAppHost
	AuthoritySecret           string
}

type installer struct {
	ctx                    context.Context
	opts                   *Options
	clientSet              kubernetes.Interface
	openxClientSet         openxclientset.Interface
	acrCredentialSecrets   cmap.ConcurrentMap
	discoverEtcdEndpoints  string
	discoveryEndpoints     string
	discoveryPort          int32
	authorityEtcdEndpoints string
	authorityEndpoints     string
	authorityPort          int32
}

func NewInstaller(ctx context.Context, options *Options) Installer {
	b := &installer{
		ctx:                  ctx,
		opts:                 options,
		clientSet:            options.ClientBuilder.ClientOrDie("builder"),
		openxClientSet:       options.ClientBuilder.OpenxClientOrDie("builder"),
		acrCredentialSecrets: cmap.New(),
	}
	return b
}

func (b *installer) Run() {
	b.installNamespaces()
	b.installDockerSecrets()
	b.installClusterRoles()
	// b.installCloudSlb()
	// b.installDiscovery()
	// b.installAuthx()
	// b.installOpenx()
	// b.installJingx()
}

func (b *installer) supportedNamespaces() []string {
	return []string{
		NamespaceKubeDiscovery,
		NamespaceKubeApi,
		NamespaceKubeAuthx,
	}
}

func (b *installer) installNamespaces() {
	for _, namespace := range b.supportedNamespaces() {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		if _, err := b.clientSet.CoreV1().Namespaces().Create(b.ctx, ns, metav1.CreateOptions{}); err != nil {
			if !errors.IsAlreadyExists(err) {
				zaplogger.Sugar().Fatal(err)
			}
			zaplogger.Sugar().Infof("Already exist Namespace:%s", ns.Name)
		} else {
			zaplogger.Sugar().Infof("Successful create Namespace:%s", ns.Name)
		}
	}
	// todo wait
	// for _, namespace := range b.supportedNamespaces() {
	// 	b.waitAcrSecrets(namespace)
	// }
}
