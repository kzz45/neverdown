package config

import (
	discoveryrestclient "github.com/kzz45/discovery/pkg/client-go/rest"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/server"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

// Config is the main context object for the controller manager.
type Config struct {
	//ComponentConfig kubectrlmgrconfig.KubeControllerManagerConfiguration
	Generic *GenericControllerManagerConfiguration

	//SecureServing *apiserver.SecureServingInfo
	// LoopbackClientConfig is a config for a privileged loopback connection
	LoopbackClientConfig *restclient.Config

	// TODO: remove deprecated insecure serving
	//InsecureServing *apiserver.DeprecatedInsecureServingInfo
	//Authentication  apiserver.AuthenticationInfo
	//Authorization   apiserver.AuthorizationInfo

	// the general kube client
	Client *clientset.Clientset

	// the client only used for leader election
	LeaderElectionClient *clientset.Clientset

	// the rest config for the master
	Kubeconfig *restclient.Config

	// AuthorityKubeConfig *discoveryrestclient.Config
	AuthorityKubeConfig *discoveryrestclient.Config
	// DiscoveryKubeConfig *discoveryrestclient.Config
	DiscoveryKubeConfig *discoveryrestclient.Config

	// the event sink
	EventRecorder record.EventRecorder

	ApiServerConfig *server.Options
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}

	//apiserver.AuthorizeClientBearerToken(c.LoopbackClientConfig, &c.Authentication, &c.Authorization)

	return &CompletedConfig{&cc}
}

func DefaultConfig() *Config {
	return &Config{
		Generic:              DefaultGenericControllerManagerConfiguration(),
		LoopbackClientConfig: nil,
		Client:               nil,
		LeaderElectionClient: nil,
		Kubeconfig:           nil,
		EventRecorder:        nil,
	}
}
