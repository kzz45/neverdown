package app

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	openxinformers "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/resources"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/server"
	"github.com/kzz45/neverdown/pkg/openx/app/config"
	"github.com/kzz45/neverdown/pkg/openx/clientbuilder"
	"github.com/kzz45/neverdown/pkg/openx/controller/loadbalancer"
	"github.com/kzz45/neverdown/pkg/openx/metrics"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
)

const (
	// ControllerStartJitter is the Jitter used when starting controller managers
	ControllerStartJitter = 1.0
	// ConfigzName is the name used for register kube-controller manager /configz, same with GroupName.
	ConfigzName = "controllermanager.config.openx.io"
)

// ControllerLoopMode is the kube-controller-manager's mode of running controller loops that are cloud provider dependent
type ControllerLoopMode int

const (
	// IncludeCloudLoops means the kube-controller-manager include the controller loops that are cloud provider dependent
	IncludeCloudLoops ControllerLoopMode = iota
	// ExternalLoops means the kube-controller-manager exclude the controller loops that are cloud provider dependent
	ExternalLoops
)

// ResyncPeriod returns a function which generates a duration each time it is
// invoked; this is so that multiple controllers don't get into lock-step and all
// hammer the apiserver with list requests simultaneously.
func ResyncPeriod(c *config.CompletedConfig) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(c.Generic.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run runs the KubeControllerManagerOptions.  This should never exit.
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	// Setup any healthz checks we will want to use.
	var checks []healthz.HealthChecker
	var electionChecker *leaderelection.HealthzAdaptor
	if c.Generic.LeaderElection.LeaderElect {
		electionChecker = leaderelection.NewLeaderHealthzAdaptor(c.Generic.LeaderElection.HealthzAdapterDuration.Duration)
		checks = append(checks, electionChecker)
	}

	run := func(ctx context.Context) {
		clientBuilder := clientbuilder.SimpleControllerClientBuilder{
			ClientConfig:          c.Kubeconfig,
			AuthxClientConfig:     c.AuthorityKubeConfig,
			DiscoveryClientConfig: c.DiscoveryKubeConfig,
		}

		controllerContext, err := CreateControllerContext(c, clientBuilder, ctx.Done())
		if err != nil {
			zaplogger.Sugar().Fatalf("error building controller context: %v", err)
		}

		if err := StartControllers(controllerContext, c.Generic.ControllerStartInterval.Duration, NewControllerInitializers(controllerContext.LoopMode)); err != nil {
			zaplogger.Sugar().Fatalf("error starting controllers: %v", err)
		}

		//controllerContext.InformerFactory.Start(controllerContext.Stop)
		// controllerContext.OpenXInformerFactory.Start(controllerContext.Stop)

		//controllerContext.ObjectOrMetadataInformerFactory.Start(controllerContext.Stop)
		close(controllerContext.InformersStarted)

		c.ApiServerConfig.Resources = resources.New(c.ApiServerConfig.Context, clientBuilder, controllerContext.InformerFactory, controllerContext.OpenXInformerFactory)

		go metrics.RunExporter()
		stop, listener, err := server.Run(c.ApiServerConfig)

		if err != nil {
			zaplogger.Sugar().Fatalf("error run apiserver: %v", err)
		}
		<-stop
		_ = listener
		//select {}
	}

	if !c.Generic.LeaderElection.LeaderElect {
		run(context.TODO())
		zaplogger.Sugar().Panic("unreachable")
	}

	id, err := os.Hostname()
	if err != nil {
		return err
	}

	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id = id + "_" + string(uuid.NewUUID())

	rl, err := resourcelock.New(c.Generic.LeaderElection.ResourceLock,
		c.Generic.LeaderElection.ResourceNamespace,
		c.Generic.LeaderElection.ResourceName,
		c.LeaderElectionClient.CoreV1(),
		c.LeaderElectionClient.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: c.EventRecorder,
		})
	if err != nil {
		zaplogger.Sugar().Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: c.Generic.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: c.Generic.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   c.Generic.LeaderElection.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				zaplogger.Sugar().Fatalf("leaderelection lost")
			},
		},
		WatchDog:        electionChecker,
		Name:            "openx-controller-manager",
		ReleaseOnCancel: c.Generic.LeaderElection.ReleaseOnCancel,
	})
	zaplogger.Sugar().Panic("unreachable")
	return nil
}

// ControllerContext defines the context object for controller
type ControllerContext struct {
	jingx *aggregator.Aggregator

	loadBalancher loadbalancer.Interface

	// ClientBuilder will provide a client for this controller to use
	ClientBuilder clientbuilder.ControllerClientBuilder

	// InformerFactory gives access to informers for the controller.
	InformerFactory informers.SharedInformerFactory

	// OpenXInformerFactory gives access to openx informers for the controller.
	OpenXInformerFactory openxinformers.SharedInformerFactory

	// ObjectOrMetadataInformerFactory gives access to informers for typed resources
	// and dynamic resources by their metadata. All generic controllers currently use
	// object metadata - if a future controller needs access to the full object this
	// would become GenericInformerFactory and take a dynamic client.
	//ObjectOrMetadataInformerFactory informerfactory.InformerFactory

	// ComponentConfig provides access to init options for a given controller
	//ComponentConfig kubectrlmgrconfig.KubeControllerManagerConfiguration

	// DeferredDiscoveryRESTMapper is a RESTMapper that will defer
	// initialization of the RESTMapper until the first mapping is
	// requested.
	//RESTMapper *restmapper.DeferredDiscoveryRESTMapper

	// AvailableResources is a map listing currently available resources
	AvailableResources map[schema.GroupVersionResource]bool

	// Cloud is the cloud provider interface for the controllers to use.
	// It must be initialized and ready to use.
	//Cloud cloudprovider.Interface

	// Control for which control loops to be run
	// IncludeCloudLoops is for a kube-controller-manager running all loops
	// ExternalLoops is for a kube-controller-manager running with a cloud-controller-manager
	LoopMode ControllerLoopMode

	// Stop is the stop channel
	Stop <-chan struct{}

	// InformersStarted is closed after all of the controllers have been initialized and are running.  After this point it is safe,
	// for an individual controller to start the shared informers. Before it is closed, they should not.
	InformersStarted chan struct{}

	// ResyncPeriod generates a duration each time it is invoked; this is so that
	// multiple controllers don't get into lock-step and all hammer the apiserver
	// with list requests simultaneously.
	ResyncPeriod func() time.Duration
}

// IsControllerEnabled checks if the context's controllers enabled or not
func (c ControllerContext) IsControllerEnabled(name string) bool {
	//return genericcontrollermanager.IsControllerEnabled(name, ControllersDisabledByDefault, c.ComponentConfig.Generic.Controllers)
	return true
}

// InitFunc is used to launch a particular controller.  It may run additional "should I activate checks".
// Any error returned will cause the controller process to `Fatal`
// The bool indicates whether the controller was enabled.
type InitFunc func(ctx *ControllerContext) (debuggingHandler http.Handler, enabled bool, err error)

type ControllerObject struct {
	Name string
	InitFunc
}

// NewControllerInitializers is a public map of named controller groups (you can start more than one in an init func)
// paired to their InitFunc.  This allows for structured downstream composition and subdivision.
func NewControllerInitializers(loopMode ControllerLoopMode) []ControllerObject {
	return []ControllerObject{
		{
			Name:     "loadbalancer",
			InitFunc: startLoadBalancerController,
		},
		{
			Name:     "mysql",
			InitFunc: startMysqlController,
		},
		{
			Name:     "redis",
			InitFunc: starRedisController,
		},
		{
			Name:     "openx",
			InitFunc: starOpenxController,
		},
		{
			Name:     "etcd",
			InitFunc: startNewEtcdController,
		},
	}
}

// CreateControllerContext creates a context struct containing references to resources needed by the
// controllers such as the cloud provider and clientBuilder. rootClientBuilder is only used for
// the shared-informers client and token controller.
func CreateControllerContext(s *config.CompletedConfig, clientBuilder clientbuilder.ControllerClientBuilder, stop <-chan struct{}) (*ControllerContext, error) {
	ctx := &ControllerContext{
		jingx:                aggregator.New(context.TODO(), clientBuilder.DiscoveryClientOrDie("openx-discovery-controller"), ""),
		ClientBuilder:        clientBuilder,
		InformerFactory:      informers.NewSharedInformerFactory(clientBuilder.ClientOrDie("shared-informers"), ResyncPeriod(s)()),
		OpenXInformerFactory: openxinformers.NewSharedInformerFactory(clientBuilder.OpenxClientOrDie("shared-informers-sageras"), ResyncPeriod(s)()),
		//ObjectOrMetadataInformerFactory: informerfactory.NewInformerFactory(sharedInformers, metadataInformers),
		//ComponentConfig:                 s.ComponentConfig,
		//RESTMapper:                      restMapper,
		//AvailableResources:              availableResources,
		//Cloud:                           cloud,
		//LoopMode:                        loopMode,
		Stop:             stop,
		InformersStarted: make(chan struct{}),
		ResyncPeriod:     ResyncPeriod(s),
	}
	ctx.runInformer(stop)
	return ctx, nil
}

// StartControllers starts a set of controllers with a specified ControllerContext
func StartControllers(ctx *ControllerContext, duration time.Duration, controllers []ControllerObject) error {
	for _, v := range controllers {
		//if !ctx.IsControllerEnabled(controllerName) {
		//	zaplogger.Sugar().Infof("%q is disabled", controllerName)
		//	continue
		//}

		controllerName := v.Name

		time.Sleep(wait.Jitter(duration, ControllerStartJitter))

		zaplogger.Sugar().Infof("Starting %q", controllerName)
		_, started, err := v.InitFunc(ctx)
		if err != nil {
			zaplogger.Sugar().Errorf("Error starting %q", controllerName)
			return err
		}
		if !started {
			zaplogger.Sugar().Infof("Skipping %q", controllerName)
			continue
		}

		zaplogger.Sugar().Infof("Started %q", controllerName)
	}

	return nil
}

func (c *ControllerContext) runSharedIndexInformer(gvk schema.GroupVersionKind, s cache.SharedIndexInformer, stopCh <-chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			zaplogger.Sugar().Errorw("GroupVersionKind Recovered",
				zap.Any("gvk", gvk),
				zap.Any("err", r))
		}
		c.runSharedIndexInformer(gvk, s, stopCh)
	}()
	s.Run(stopCh)
}
