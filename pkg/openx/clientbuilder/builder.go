package clientbuilder

import (
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	nativakubernetes "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// ControllerClientBuilder allows you to get clients and configs for controllers
// Please note a copy also exists in staging/src/k8s.io/cloud-provider/cloud.go
// TODO: Extract this into a separate controller utilities repo (issues/68947)
type ControllerClientBuilder interface {
	Config(name string) (*restclient.Config, error)
	ConfigOrDie(name string) *restclient.Config
	Client(name string) (nativakubernetes.Interface, error)
	ClientOrDie(name string) nativakubernetes.Interface
	MetricsClient(name string) (metricsclientset.Interface, error)
	MetricsClientOrDie(name string) metricsclientset.Interface
	OpenxClient(name string) (kubernetes.Interface, error)
	OpenxClientOrDie(name string) kubernetes.Interface
	// DiscoveryConfig was the discovery config
	DiscoveryConfig(name string) (*restclient.Config, error)
	DiscoveryConfigOrDie(name string) *restclient.Config
	DiscoveryClient(name string) (kubernetes.Interface, error)
	DiscoveryClientOrDie(name string) kubernetes.Interface
	// AuthxConfig was the authority config
	AuthxConfig(name string) (*restclient.Config, error)
	AuthxConfigOrDie(name string) *restclient.Config
	AuthxClient(name string) (kubernetes.Interface, error)
	AuthxClientOrDie(name string) kubernetes.Interface
}

// SimpleControllerClientBuilder returns a fixed client with different user agents
type SimpleControllerClientBuilder struct {
	// ClientConfig is a skeleton config to clone and use as the basis for each controller client
	ClientConfig *restclient.Config
	// DiscoveryClientConfig is a skeleton config to clone and use as the basis for each controller client
	DiscoveryClientConfig *restclient.Config
	// AuthxClientConfig is a skeleton config to clone and use as the basis for each controller client
	AuthxClientConfig *restclient.Config
}

// Config returns a client config for a fixed client
func (b SimpleControllerClientBuilder) Config(name string) (*restclient.Config, error) {
	clientConfig := *b.ClientConfig
	return restclient.AddUserAgent(&clientConfig, name), nil
}

// ConfigOrDie returns a client config if no error from previous config func.
// If it gets an error getting the client, it will log the error and kill the process it's running in.
func (b SimpleControllerClientBuilder) ConfigOrDie(name string) *restclient.Config {
	clientConfig, err := b.Config(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return clientConfig
}

// Client returns a clientset.Interface built from the ClientBuilder
func (b SimpleControllerClientBuilder) Client(name string) (nativakubernetes.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return nativakubernetes.NewForConfig(clientConfig)
}

// ClientOrDie returns a clientset.interface built from the ClientBuilder with no error.
// If it gets an error getting the client, it will log the error and kill the process it's running in.
func (b SimpleControllerClientBuilder) ClientOrDie(name string) nativakubernetes.Interface {
	client, err := b.Client(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return client
}

// MetricsClient returns a metricsclientset.Interface built from the ClientBuilder
func (b SimpleControllerClientBuilder) MetricsClient(name string) (metricsclientset.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return metricsclientset.NewForConfig(clientConfig)
}

// MetricsClientOrDie returns a metricsclientset.interface built from the ClientBuilder with no error.
// If it gets an error getting the client, it will log the error and kill the process it's running in.
func (b SimpleControllerClientBuilder) MetricsClientOrDie(name string) metricsclientset.Interface {
	client, err := b.MetricsClient(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return client
}

// OpenxClient returns a kubernetes.Interface built from the ClientBuilder
func (b SimpleControllerClientBuilder) OpenxClient(name string) (kubernetes.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(clientConfig)
}

// OpenxClientOrDie returns a kubernetes.interface built from the ClientBuilder with no error.
// If it gets an error getting the client, it will log the error and kill the process it's running in.
func (b SimpleControllerClientBuilder) OpenxClientOrDie(name string) kubernetes.Interface {
	client, err := b.OpenxClient(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return client
}

func (b SimpleControllerClientBuilder) DiscoveryConfig(name string) (*restclient.Config, error) {
	clientConfig := *b.DiscoveryClientConfig
	return restclient.AddUserAgent(&clientConfig, name), nil
}

func (b SimpleControllerClientBuilder) DiscoveryConfigOrDie(name string) *restclient.Config {
	clientConfig, err := b.DiscoveryConfig(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return clientConfig
}

func (b SimpleControllerClientBuilder) DiscoveryClient(name string) (kubernetes.Interface, error) {
	clientConfig, err := b.DiscoveryConfig(name)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(clientConfig)
}

func (b SimpleControllerClientBuilder) DiscoveryClientOrDie(name string) kubernetes.Interface {
	client, err := b.DiscoveryClient(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return client
}

func (b SimpleControllerClientBuilder) AuthxConfig(name string) (*restclient.Config, error) {
	clientConfig := *b.AuthxClientConfig
	return restclient.AddUserAgent(&clientConfig, name), nil
}

func (b SimpleControllerClientBuilder) AuthxConfigOrDie(name string) *restclient.Config {
	clientConfig, err := b.AuthxConfig(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return clientConfig
}

func (b SimpleControllerClientBuilder) AuthxClient(name string) (kubernetes.Interface, error) {
	clientConfig, err := b.AuthxConfig(name)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(clientConfig)
}

func (b SimpleControllerClientBuilder) AuthxClientOrDie(name string) kubernetes.Interface {
	client, err := b.AuthxClient(name)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return client
}
