package controller

import (
	"context"
	"sync"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"

	openx "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	openxinformers "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions"
	openxlistersv1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"
)

type Option struct {
	LeaseLock      *resourcelock.LeaseLock
	ClientSet      kubernetes.Interface
	Openx          openx.Interface
	ResyncInterval time.Duration
}

type Controller struct {
	option    *Option
	name      string
	namespace string
	clientSet kubernetes.Interface

	wg sync.WaitGroup

	// recorder event
	recorder record.EventRecorder

	// leaderContextCancelFunc will be called when apisix-ingress-controller
	// decides to give up its leader role.
	leaderContextCancelFunc context.CancelFunc

	// core informers and listers
	namespaceInformer cache.SharedIndexInformer
	namespaceLister   corelistersv1.NamespaceLister
	podInformer       cache.SharedIndexInformer
	podLister         corelistersv1.PodLister
	svcInformer       cache.SharedIndexInformer
	svcLister         corelistersv1.ServiceLister
	ingressInformer   cache.SharedIndexInformer
	secretInformer    cache.SharedIndexInformer
	secretLister      corelistersv1.SecretLister

	// apps informers and listers
	deploymentInformer  cache.SharedIndexInformer
	deploymentLister    appslistersv1.DeploymentLister
	statefulSetInformer cache.SharedIndexInformer
	statefulSetLister   appslistersv1.StatefulSetLister

	// crd informers and listers
	mysqlInformer cache.SharedIndexInformer
	mysqlLister   openxlistersv1.MysqlLister
	redisInformer cache.SharedIndexInformer
	RedisLister   openxlistersv1.RedisLister
	openxInformer cache.SharedIndexInformer
	OpenxLister   openxlistersv1.OpenxLister
}

func NewController(opt *Option) *Controller {

	kubeFactory := informers.NewSharedInformerFactory(opt.ClientSet, opt.ResyncInterval)
	openxFactory := openxinformers.NewSharedInformerFactory(opt.Openx, opt.ResyncInterval)

	c := &Controller{
		option:                  opt,
		name:                    opt.LeaseLock.LeaseMeta.Name,
		namespace:               opt.LeaseLock.LeaseMeta.Namespace,
		clientSet:               nil,
		recorder:                nil,
		leaderContextCancelFunc: nil,
		namespaceInformer:       nil,
		namespaceLister:         nil,
		podInformer:             kubeFactory.Core().V1().Pods().Informer(),
		podLister:               kubeFactory.Core().V1().Pods().Lister(),
		svcInformer:             kubeFactory.Core().V1().Services().Informer(),
		svcLister:               kubeFactory.Core().V1().Services().Lister(),
		ingressInformer:         nil,
		secretInformer:          kubeFactory.Core().V1().Secrets().Informer(),
		secretLister:            kubeFactory.Core().V1().Secrets().Lister(),
		deploymentInformer:      kubeFactory.Apps().V1().Deployments().Informer(),
		deploymentLister:        kubeFactory.Apps().V1().Deployments().Lister(),
		statefulSetInformer:     kubeFactory.Apps().V1().StatefulSets().Informer(),
		statefulSetLister:       kubeFactory.Apps().V1().StatefulSets().Lister(),
		mysqlInformer:           openxFactory.Openx().V1().Mysqls().Informer(),
		mysqlLister:             openxFactory.Openx().V1().Mysqls().Lister(),
		redisInformer:           openxFactory.Openx().V1().Redises().Informer(),
		RedisLister:             openxFactory.Openx().V1().Redises().Lister(),
		openxInformer:           openxFactory.Openx().V1().Openxes().Informer(),
		OpenxLister:             openxFactory.Openx().V1().Openxes().Lister(),
	}
	return c
}

// recorderEvent recorder events for resources
func (c *Controller) recorderEventS(object runtime.Object, eventtype, reason string, msg string) {
	c.recorder.Event(object, eventtype, reason, msg)
}

func (c *Controller) goAttach(handler func()) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		handler()
	}()
}

// Eventf implements the resourcelock.EventRecorder interface.
func (c *Controller) Eventf(_ runtime.Object, eventType string, reason string, message string, _ ...interface{}) {
	zaplogger.Sugar().Infow(reason, zap.String("message", message), zap.String("event_type", eventType))
}

// Run launches the controller.
func (c *Controller) Run(stop chan struct{}) error {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()
	go func() {
		<-stop
		rootCancel()
	}()
	//c.MetricsCollector.ResetLeader(false)

	//go func() {
	//	if err := c.apiServer.Run(rootCtx.Done()); err != nil {
	//		zaplogger.Sugar().Errorf("failed to launch API Server: %s", err)
	//	}
	//}()

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Namespace: c.option.LeaseLock.LeaseMeta.Namespace,
			Name:      c.option.LeaseLock.LeaseMeta.Name,
		},
		Client: c.option.ClientSet.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      c.option.LeaseLock.LockConfig.Identity,
			EventRecorder: c,
		},
	}
	cfg := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 5 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: c.start,
			OnNewLeader: func(identity string) {
				zaplogger.Sugar().Warnf("found a new leader %s", identity)
				if identity != c.option.LeaseLock.Identity() {
					zaplogger.Sugar().Infow("controller now is running as a candidate",
						zap.String("namespace", c.namespace),
						zap.String("pod", c.name),
					)
					//c.MetricsCollector.ResetLeader(false)
					// delete the old APISIX cluster, so that the cached state
					// like synchronization won't be used next time the candidate
					// becomes the leader again.
					//c.apisix.DeleteCluster(c.cfg.APISIX.DefaultClusterName)
				}
			},
			OnStoppedLeading: func() {
				zaplogger.Sugar().Infow("controller now is running as a candidate",
					zap.String("namespace", c.namespace),
					zap.String("pod", c.name),
				)
				//c.MetricsCollector.ResetLeader(false)
				// delete the old APISIX cluster, so that the cached state
				// like synchronization won't be used next time the candidate
				// becomes the leader again.
				//c.apisix.DeleteCluster(c.cfg.APISIX.DefaultClusterName)
			},
		},
		ReleaseOnCancel: true,
		Name:            "openx-controller",
	}

	elector, err := leaderelection.NewLeaderElector(cfg)
	if err != nil {
		zaplogger.Sugar().Errorf("failed to create leader elector: %s", err.Error())
		return err
	}

election:
	curCtx, cancel := context.WithCancel(rootCtx)
	c.leaderContextCancelFunc = cancel
	elector.Run(curCtx)
	select {
	case <-rootCtx.Done():
		return nil
	default:
		goto election
	}
}

func (c *Controller) runInformer(ctx context.Context) {
	//c.goAttach(func() {
	//	c.namespaceInformer.Run(ctx.Done())
	//})
	c.goAttach(func() {
		c.podInformer.Run(ctx.Done())
	})
	c.goAttach(func() {
		c.svcInformer.Run(ctx.Done())
	})
	//c.goAttach(func() {
	//	c.ingressInformer.Run(ctx.Done())
	//})
	c.goAttach(func() {
		c.secretInformer.Run(ctx.Done())
	})
}

func (c *Controller) start(ctx context.Context) {
	zaplogger.Sugar().Infow("controller tries to leading ...",
		zap.String("namespace", c.namespace),
		zap.String("pod", c.name),
	)

	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(ctx)
	defer cancelFunc()

	// give up leader
	defer c.leaderContextCancelFunc()

	c.runInformer(ctx)
	//c.MetricsCollector.ResetLeader(true)

	zaplogger.Sugar().Infow("controller now is running as leader",
		zap.String("namespace", c.namespace),
		zap.String("pod", c.name),
	)

	<-ctx.Done()
	c.wg.Wait()
}
