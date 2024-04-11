package config

import (
	"time"

	leaderconfig "github.com/kzz45/neverdown/pkg/openx/leader/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenericControllerManagerConfiguration holds configuration for a generic controller-manager
type GenericControllerManagerConfiguration struct {
	// port is the port that the controller-manager's http service runs on.
	Port int32
	// address is the IP address to serve on (set to 0.0.0.0 for all interfaces).
	Address string
	// minResyncPeriod is the resync period in reflectors; will be random between
	// minResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod metav1.Duration
	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	//ClientConnection componentbaseconfig.ClientConnectionConfiguration
	// How long to wait between starting controller managers
	ControllerStartInterval metav1.Duration
	// leaderElection defines the configuration of leader election client.
	LeaderElection *leaderconfig.LeaderElectionConfiguration
	// Controllers is the list of controllers to enable or disable
	// '*' means "all enabled by default controllers"
	// 'foo' means "enable 'foo'"
	// '-foo' means "disable 'foo'"
	// first item for a particular name wins
	Controllers []string
	// DebuggingConfiguration holds configuration for Debugging related features.
	//Debugging componentbaseconfig.DebuggingConfiguration
}

func DefaultGenericControllerManagerConfiguration() *GenericControllerManagerConfiguration {
	return &GenericControllerManagerConfiguration{
		Port:                    0,
		Address:                 "",
		MinResyncPeriod:         metav1.Duration{Duration: 60 * time.Second},
		ControllerStartInterval: metav1.Duration{Duration: 10 * time.Second},
		LeaderElection:          leaderconfig.DefaultLeaderConfig(),
		Controllers:             nil,
	}
}
