package options

import (
	"github.com/kzz45/neverdown/pkg/openx/aggregator/server"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	kubecontrollerconfig "github.com/kzz45/neverdown/pkg/openx/app/config"

	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	clientgokubescheme "k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

const (
	// KubeControllerManagerUserAgent is the userAgent name when starting kube-controller managers.
	KubeControllerManagerUserAgent = "openx-controller-manager"
)

// KubeControllerManagerOptions is the main context object for the kube-controller manager.
type KubeControllerManagerOptions struct {
	Generic *kubecontrollerconfig.GenericControllerManagerConfiguration
	//KubeCloudShared   *cpoptions.KubeCloudSharedOptions
	//ServiceController *cpoptions.ServiceControllerOptions

	//AttachDetachController           *AttachDetachControllerOptions
	//CSRSigningController             *CSRSigningControllerOptions
	//DaemonSetController              *DaemonSetControllerOptions
	//DeploymentController             *DeploymentControllerOptions
	//StatefulSetController            *StatefulSetControllerOptions
	//DeprecatedFlags                  *DeprecatedControllerOptions
	//EndpointController               *EndpointControllerOptions
	//EndpointSliceController          *EndpointSliceControllerOptions
	//EndpointSliceMirroringController *EndpointSliceMirroringControllerOptions
	//GarbageCollectorController       *GarbageCollectorControllerOptions
	//HPAController                    *HPAControllerOptions
	//JobController                    *JobControllerOptions
	//CronJobController                *CronJobControllerOptions
	//NamespaceController              *NamespaceControllerOptions
	//NodeIPAMController               *NodeIPAMControllerOptions
	//NodeLifecycleController          *NodeLifecycleControllerOptions
	//PersistentVolumeBinderController *PersistentVolumeBinderControllerOptions
	//PodGCController                  *PodGCControllerOptions
	//ReplicaSetController             *ReplicaSetControllerOptions
	//ReplicationController            *ReplicationControllerOptions
	//ResourceQuotaController          *ResourceQuotaControllerOptions
	//SAController                     *SAControllerOptions
	//TTLAfterFinishedController       *TTLAfterFinishedControllerOptions
	//
	//SecureServing *apiserveroptions.SecureServingOptionsWithLoopback
	//// TODO: remove insecure serving mode
	//InsecureServing *apiserveroptions.DeprecatedInsecureServingOptionsWithLoopback
	//Authentication  *apiserveroptions.DelegatingAuthenticationOptions
	//Authorization   *apiserveroptions.DelegatingAuthorizationOptions
	//Metrics         *metrics.Options
	//Logs            *logs.Options

	Master                      string
	Kubeconfig                  string
	ShowHiddenMetricsForVersion string

	ApiServerConfig *server.Options
}

// Config return a controller manager config objective
func (s *KubeControllerManagerOptions) Config(allControllers []string, disabledByDefaultControllers []string) (*kubecontrollerconfig.Config, error) {
	//if err := s.Validate(allControllers, disabledByDefaultControllers); err != nil {
	//	return nil, err
	//}
	//
	//if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
	//	return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	//}

	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Master, s.Kubeconfig)
	if err != nil {
		return nil, err
	}
	//kubeconfig.DisableCompression = true
	//kubeconfig.ContentConfig.AcceptContentTypes = s.Generic.ClientConnection.AcceptContentTypes
	//kubeconfig.ContentConfig.ContentType = s.Generic.ClientConnection.ContentType
	//kubeconfig.QPS = s.Generic.ClientConnection.QPS
	//kubeconfig.Burst = int(s.Generic.ClientConnection.Burst)

	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeconfig, KubeControllerManagerUserAgent))
	if err != nil {
		return nil, err
	}

	// shallow copy, do not modify the kubeconfig.Timeout.
	config := *kubeconfig
	config.Timeout = s.Generic.LeaderElection.RenewDeadline.Duration
	leaderElectionClient := clientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))

	eventRecorder := createRecorder(client, KubeControllerManagerUserAgent)

	c := &kubecontrollerconfig.Config{
		Generic:              kubecontrollerconfig.DefaultGenericControllerManagerConfiguration(),
		LoopbackClientConfig: nil,
		Client:               client,
		Kubeconfig:           kubeconfig,
		AuthorityKubeConfig:  loadAuthorityKubernetes(),
		DiscoveryKubeConfig:  loadDiscoveryKubernetes(),
		EventRecorder:        eventRecorder,
		LeaderElectionClient: leaderElectionClient,
		ApiServerConfig:      s.ApiServerConfig,
	}

	//if err := s.ApplyTo(c); err != nil {
	//	return nil, err
	//}
	//s.Metrics.Apply()
	//
	//s.Logs.Apply()

	return c, nil
}

func createRecorder(kubeClient clientset.Interface, userAgent string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(clientgokubescheme.Scheme, corev1.EventSource{Component: userAgent})
}

func loadDiscoveryKubernetes() *restclient.Config {
	kubeconfig, err := restclient.InDicoveryClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return kubeconfig
}

func loadAuthorityKubernetes() *restclient.Config {
	kubeconfig, err := restclient.InAuthorityClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return kubeconfig
}
