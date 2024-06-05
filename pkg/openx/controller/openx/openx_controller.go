package openx

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	"github.com/kzz45/discovery/pkg/jingx/aggregator"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	openx "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	openxinformersv1 "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions/openx/v1"
	openxlistersv1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"
	"github.com/kzz45/neverdown/pkg/openx/controller/loadbalancer"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	autoscalinginformers "k8s.io/client-go/informers/autoscaling/v2"
	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	autoscalinglisters "k8s.io/client-go/listers/autoscaling/v2"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const (
	// maxRetries is the number of times a openx will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a openx is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = schema.GroupVersion{Group: openxv1.GroupVersion.Group, Version: openxv1.GroupVersion.Version}.WithKind("Openx")

var controllerName = "openx"

// OpenxController is responsible for synchronizing openx objects stored
// in the system with actual running replica sets and pods.
type OpenxController struct {
	jingx        *aggregator.Aggregator
	loadbalancer loadbalancer.Interface
	// svcControl is used for adopting/releasing services.
	svcControl controller.ServiceControlInterface
	// dpControl is used for adopting/releasing deployments.
	dpControl controller.DeploymentControlInterface
	// hpaControl is used for adopting/releasing HorizontalPodAutoscalers.
	hpaControl controller.HorizontalPodAutoscalerControlInterface

	client        clientset.Interface
	openx         openx.Interface
	eventRecorder record.EventRecorder

	syncHandler  func(key string) error
	enqueueOpenx func(openx *openxv1.Openx)

	openxLister openxlistersv1.OpenxLister
	svcLister   corelisters.ServiceLister
	dpLister    appslisters.DeploymentLister
	hpaLister   autoscalinglisters.HorizontalPodAutoscalerLister

	openxListerSynced cache.InformerSynced
	svcListerSynced   cache.InformerSynced
	dpListerSynced    cache.InformerSynced
	hpaListerSynced   cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

// NewOpenxController creates a new OpenxController.
func NewOpenxController(
	jingx *aggregator.Aggregator,
	loadbalancher loadbalancer.Interface,
	openxInformer openxinformersv1.OpenxInformer,
	svcInformer coreinformers.ServiceInformer,
	dpInformer appsinformers.DeploymentInformer,
	hpaInformer autoscalinginformers.HorizontalPodAutoscalerInformer,
	client clientset.Interface,
	openx openx.Interface,
) (*OpenxController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events(controller.RecorderNamespace)})

	//if client != nil && client.CoreV1().RESTClient().GetRateLimiter() != nil {
	//if err := ratelimiter.RegisterMetricAndTrackRateLimiterUsage("deployment_controller", client.CoreV1().RESTClient().GetRateLimiter()); err != nil {
	//	return nil, err
	//}
	//}
	tc := &OpenxController{
		jingx:         jingx,
		loadbalancer:  loadbalancher,
		client:        client,
		openx:         openx,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "openx-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	tc.openxLister = openxInformer.Lister()
	tc.svcLister = svcInformer.Lister()
	tc.dpLister = dpInformer.Lister()
	tc.hpaLister = hpaInformer.Lister()
	tc.openxListerSynced = openxInformer.Informer().HasSynced
	tc.svcListerSynced = svcInformer.Informer().HasSynced
	tc.dpListerSynced = dpInformer.Informer().HasSynced
	tc.hpaListerSynced = hpaInformer.Informer().HasSynced

	tc.syncHandler = tc.syncOpenx
	tc.enqueueOpenx = tc.enqueue

	openxInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    tc.addOpenx,
		UpdateFunc: tc.updateOpenx,
		// This will enter the sync loop and no-op, because the openx has been deleted from the store.
		DeleteFunc: tc.deleteOpenx,
	})
	svcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    tc.addService,
		UpdateFunc: tc.updateService,
		DeleteFunc: tc.deleteService,
	})
	dpInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    tc.addDeployment,
		UpdateFunc: tc.updateDeployment,
		DeleteFunc: tc.deleteDeployment,
	})
	hpaInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    tc.addHorizontalPodAutoscaler,
		UpdateFunc: tc.updateHorizontalPodAutoscaler,
		DeleteFunc: tc.deleteHorizontalPodAutoscaler,
	})

	tc.svcControl = &controller.RealServiceControl{
		KubeClient: client,
		Recorder:   tc.eventRecorder,
	}
	tc.dpControl = &controller.RealDeploymentControl{
		KubeClient: client,
		Recorder:   tc.eventRecorder,
	}
	tc.hpaControl = &controller.RealHorizontalPodAutoscalerControl{
		KubeClient: client,
		Recorder:   tc.eventRecorder,
	}
	return tc, nil
}

// Run begins watching and syncing.
func (tc *OpenxController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer tc.queue.ShutDown()

	zaplogger.Sugar().Infof("Starting openx controller")
	defer zaplogger.Sugar().Infof("Shutting down openx controller")

	if !cache.WaitForNamedCacheSync(controllerName, stopCh, tc.openxListerSynced, tc.svcListerSynced, tc.dpListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(tc.worker, time.Second, stopCh)
	}

	go tc.watch()
	<-stopCh
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (tc *OpenxController) worker() {
	for tc.processNextWorkItem() {
	}
}

func (tc *OpenxController) processNextWorkItem() bool {
	key, quit := tc.queue.Get()
	if quit {
		return false
	}
	defer tc.queue.Done(key)

	err := tc.syncHandler(key.(string))
	tc.handleErr(err, key)

	return true
}

func (tc *OpenxController) handleErr(err error, key interface{}) {
	if err == nil || errors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
		tc.queue.Forget(key)
		return
	}

	ns, name, keyErr := cache.SplitMetaNamespaceKey(key.(string))
	if keyErr != nil {
		zaplogger.Sugar().Errorw("Failed to split meta namespace cache key", zap.Any("key", key), zap.Error(err))
	}

	if tc.queue.NumRequeues(key) < maxRetries {
		zaplogger.Sugar().Infow("Error syncing",
			zap.String("controller", controllerName),
			zap.String("namespace", ns),
			zap.String("name", name),
			zap.Error(err))
		tc.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	zaplogger.Sugar().Infof("Dropping openx %q out of the queue: %v", key, err)
	tc.queue.Forget(key)
}

// resolveControllerRef returns the controller referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching controller
// of the correct Kind.
func (tc *OpenxController) resolveControllerRef(namespace string, controllerRef *metav1.OwnerReference) *openxv1.Openx {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != controllerKind.Kind {
		return nil
	}
	d, err := tc.openxLister.Openxes(namespace).Get(controllerRef.Name)
	if err != nil {
		return nil
	}
	if d.UID != controllerRef.UID {
		// The controller we found with this Name is not the same one that the
		// ControllerRef points to.
		return nil
	}
	return d
}

func (tc *OpenxController) addOpenx(obj interface{}) {
	d := obj.(*openxv1.Openx)
	zaplogger.Sugar().Infof("Adding openx:%s obj: %#v", d.Name, d)
	tc.enqueueOpenx(d)
}

func (tc *OpenxController) updateOpenx(old, cur interface{}) {
	oldD := old.(*openxv1.Openx)
	curD := cur.(*openxv1.Openx)
	if curD.ResourceVersion == oldD.ResourceVersion {
		return
	}
	zaplogger.Sugar().Infof("Updating trigger openx namespace:%s name:%s", oldD.Namespace, oldD.Name)
	tc.enqueueOpenx(curD)
}

func (tc *OpenxController) deleteOpenx(obj interface{}) {
	d, ok := obj.(*openxv1.Openx)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		d, ok = tombstone.Obj.(*openxv1.Openx)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a openx %#v", obj))
			return
		}
	}
	zaplogger.Sugar().Infof("Deleting openx %s", d.Name)
	tc.enqueueOpenx(d)
}

func (tc *OpenxController) enqueue(openx *openxv1.Openx) {
	key, err := controller.KeyFunc(openx)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", openx, err))
		return
	}

	tc.queue.Add(key)
}

func (tc *OpenxController) enqueueRateLimited(openx *openxv1.Openx) {
	key, err := controller.KeyFunc(openx)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", openx, err))
		return
	}

	tc.queue.AddRateLimited(key)
}

// enqueueAfter will enqueue a deployment after the provided amount of time.
func (tc *OpenxController) enqueueAfter(openx *openxv1.Openx, after time.Duration) {
	key, err := controller.KeyFunc(openx)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", openx, err))
		return
	}

	tc.queue.AddAfter(key, after)
}

// RecheckDeletionTimestamp returns a CanAdopt() function to recheck deletion.
//
// The CanAdopt() function calls getObject() to fetch the latest value,
// and denies adoption attempts if that object has a non-nil DeletionTimestamp.
func RecheckDeletionTimestamp(getObject func() (metav1.Object, error)) func() error {
	return func() error {
		obj, err := getObject()
		if err != nil {
			return fmt.Errorf("can't recheck DeletionTimestamp: %v", err)
		}
		if obj.GetDeletionTimestamp() != nil {
			return fmt.Errorf("%v/%v has just been deleted at %v", obj.GetNamespace(), obj.GetName(), obj.GetDeletionTimestamp())
		}
		return nil
	}
}

func (tc *OpenxController) getServicesForOpenx(m *openxv1.Openx) ([]*corev1.Service, error) {
	// List all Services to find those we own but that no longer match our
	// selector. They will be orphaned by ClaimServices().
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, fmt.Errorf("openx %s/%s has invalid label selector: %v", m.Namespace, m.Name, err)
	}
	svcList, err := tc.svcLister.Services(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := tc.openx.OpenxV1().Openxes(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original openx %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewServiceControllerRefManager(tc.svcControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimServices(svcList)
}

func (tc *OpenxController) getDeploymentsForOpenx(m *openxv1.Openx) ([]*appsv1.Deployment, error) {
	// Get all Pods that potentially belong to this Deployment.
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, err
	}
	dpList, err := tc.dpLister.Deployments(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := tc.openx.OpenxV1().Openxes(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original openx %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewDeploymentControllerRefManager(tc.dpControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimDeployments(dpList)
}

func (tc *OpenxController) getHorizontalPodAutoscalersForOpenx(m *openxv1.Openx) ([]*autoscalingv2.HorizontalPodAutoscaler, error) {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, err
	}
	hpaList, err := tc.hpaLister.HorizontalPodAutoscalers(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := tc.openx.OpenxV1().Openxes(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original openx %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewHorizontalPodAutoscalerControllerRefManager(tc.hpaControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimHorizontalPodAutoscalers(hpaList)
}

// syncOpenx will sync the openx with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (tc *OpenxController) syncOpenx(key string) error {
	startTime := time.Now()
	zaplogger.Sugar().Debugf("Started syncing openx %q (%v)", key, startTime)
	defer func() {
		zaplogger.Sugar().Debugf("Finished syncing openx %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	openx, err := tc.openxLister.Openxes(namespace).Get(name)
	if errors.IsNotFound(err) {
		zaplogger.Sugar().Debugf("Deployment %v has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	m := openx.DeepCopy()
	m.Labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: m.Name,
	}

	// List Services owned by this openx, while reconciling ControllerRef
	// through adoption/orphaning.
	svcList, err := tc.getServicesForOpenx(m)
	if err != nil {
		return err
	}
	// List all Pods owned by this Deployment, grouped by their ReplicaSet.
	// Current uses of the podMap are:
	//
	// * check if a Pod is labeled correctly with the pod-template-hash label.
	// * check that no old Pods are running in the middle of Recreate Deployments.
	dpList, err := tc.getDeploymentsForOpenx(m)
	if err != nil {
		return err
	}
	hpaList, err := tc.getHorizontalPodAutoscalersForOpenx(m)
	if err != nil {
		return err
	}

	if m.DeletionTimestamp != nil {
		return tc.syncStatusOnly(m, dpList, hpaList)
	}

	// compare and update hpa
	if err = tc.syncHorizontalPodAutoscalers(m, hpaList); err != nil {
		tc.eventRecorder.Event(openx.Recorder(), corev1.EventTypeWarning, controller.OpenxSyncAborted, err.Error())
		return err
	}
	// compare and update svc
	if err = tc.syncServices(m, svcList); err != nil {
		tc.eventRecorder.Event(openx.Recorder(), corev1.EventTypeWarning, controller.OpenxSyncAborted, err.Error())
		zaplogger.Sugar().Error(err)
		return err
	}
	// compare and update sts
	if err = tc.syncDeployments(m, dpList, hpaScaleReplicas(hpaList)); err != nil {
		tc.eventRecorder.Event(openx.Recorder(), corev1.EventTypeWarning, controller.OpenxSyncAborted, err.Error())
		return err
	}
	return tc.syncStatusOnly(m, dpList, hpaList)
}

// addService enqueues the openx that manages a Service when the Service is created.
func (tc *OpenxController) addService(obj interface{}) {
	svc := obj.(*corev1.Service)

	if svc.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		tc.deleteService(svc)
		return
	}

	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(svc)
	if controllerRef == nil {
		//zaplogger.Sugar().Infof("Orphan Service:%s namespace:%s", svc.Name, svc.Namespace)
		return
	}
	//zaplogger.Sugar().Infof("Service:%s namespace:%s controllerRef:%#v", svc.Name, svc.Namespace, controllerRef)
	m := tc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Debugf("Service %s namespace:%s added.", svc.Name, svc.Namespace)
	tc.enqueueOpenx(m)
}

func (tc *OpenxController) getOpenxsForService(svc *corev1.Service) []*openxv1.Openx {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ts, err := tc.openxLister.Openxes(svc.Namespace).List(selector)
	if err != nil || len(ts) == 0 {
		return nil
	}
	// Because all ReplicaSet's belonging to a deployment should have a unique label key,
	// there should never be more than one deployment returned by the above method.
	// If that happens we should probably dynamically repair the situation by ultimately
	// trying to clean up one of the controllers, for now we just return th	e older one
	if len(ts) > 1 {
		// ControllerRef will ensure we don't do anything crazy, but more than one
		// item in this list nevertheless constitutes user error.
		zaplogger.Sugar().Infof("user error! more than one openx is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ts[0].Namespace, ts[0].Name)
	}
	return ts
}

// updateService figures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (tc *OpenxController) updateService(old, cur interface{}) {
	curRS := cur.(*corev1.Service)
	oldRS := old.(*corev1.Service)
	if curRS.ResourceVersion == oldRS.ResourceVersion {
		// Periodic resync will send update events for all known replica sets.
		// Two different versions of the same replica set will always have different RVs.
		return
	}

	curControllerRef := metav1.GetControllerOf(curRS)
	oldControllerRef := metav1.GetControllerOf(oldRS)
	controllerRefChanged := !reflect.DeepEqual(curControllerRef, oldControllerRef)
	if controllerRefChanged && oldControllerRef != nil {
		// The ControllerRef was changed. Sync the old controller, if any.
		if d := tc.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			tc.enqueueOpenx(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := tc.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("Service %s updated.", curRS.Name)
		tc.enqueueOpenx(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := tc.getOpenxsForService(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan Service %s updated.", curRS.Name)
		for _, d := range ds {
			tc.enqueueOpenx(d)
		}
	}
}

// deleteReplicaSet enqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (tc *OpenxController) deleteService(obj interface{}) {
	svc, ok := obj.(*corev1.Service)

	// When a delete is dropped, the relist will notice a pod in the store not
	// in the list, leading to the insertion of a tombstone object which contains
	// the deleted key/value. Note that this value might be stale. If the ReplicaSet
	// changed labels the new deployment will not be woken up till the periodic resync.
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		svc, ok = tombstone.Obj.(*corev1.Service)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Service %#v", obj))
			return
		}
	}

	controllerRef := metav1.GetControllerOf(svc)
	if controllerRef == nil {
		// No controller should care about orphans being deleted.
		return
	}
	m := tc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Service %s deleted.", svc.Name)
	tc.enqueueOpenx(m)
}

// addDeployment enqueues the openx that manages a Service when the Deploymentis created.
func (tc *OpenxController) addDeployment(obj interface{}) {
	dp := obj.(*appsv1.Deployment)

	if dp.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		tc.deleteService(dp)
		return
	}
	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(dp)
	if controllerRef == nil {
		return
	}
	m := tc.resolveControllerRef(dp.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Deployment%s namespace:%s added.", dp.Name, dp.Namespace)
	tc.enqueueOpenx(m)
}

func (tc *OpenxController) getOpenxsForDeployment(svc *appsv1.Deployment) []*openxv1.Openx {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ms, err := tc.openxLister.Openxes(svc.Namespace).List(selector)
	if err != nil || len(ms) == 0 {
		return nil
	}
	// Because all ReplicaSet's belonging to a deployment should have a unique label key,
	// there should never be more than one deployment returned by the above method.
	// If that happens we should probably dynamically repair the situation by ultimately
	// trying to clean up one of the controllers, for now we just return th	e older one
	if len(ms) > 1 {
		// ControllerRef will ensure we don't do anything crazy, but more than one
		// item in this list nevertheless constitutes user error.
		zaplogger.Sugar().Infof("user error! more than one openx is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

// updateDeploymentfigures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (tc *OpenxController) updateDeployment(old, cur interface{}) {
	curRS := cur.(*appsv1.Deployment)
	oldRS := old.(*appsv1.Deployment)
	if curRS.ResourceVersion == oldRS.ResourceVersion {
		// Periodic resync will send update events for all known replica sets.
		// Two different versions of the same replica set will always have different RVs.
		return
	}

	curControllerRef := metav1.GetControllerOf(curRS)
	oldControllerRef := metav1.GetControllerOf(oldRS)
	controllerRefChanged := !reflect.DeepEqual(curControllerRef, oldControllerRef)
	if controllerRefChanged && oldControllerRef != nil {
		// The ControllerRef was changed. Sync the old controller, if any.
		if d := tc.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			tc.enqueueOpenx(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := tc.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		//zaplogger.Sugar().Infof("Deployment %s updated.", curRS.Name)
		tc.enqueueOpenx(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := tc.getOpenxsForDeployment(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan Deployment %s updated.", curRS.Name)
		for _, d := range ds {
			tc.enqueueOpenx(d)
		}
	}
}

// deleteDeploymentenqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (tc *OpenxController) deleteDeployment(obj interface{}) {
	svc, ok := obj.(*appsv1.Deployment)

	// When a delete is dropped, the relist will notice a pod in the store not
	// in the list, leading to the insertion of a tombstone object which contains
	// the deleted key/value. Note that this value might be stale. If the ReplicaSet
	// changed labels the new deployment will not be woken up till the periodic resync.
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		svc, ok = tombstone.Obj.(*appsv1.Deployment)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Service %#v", obj))
			return
		}
	}

	controllerRef := metav1.GetControllerOf(svc)
	if controllerRef == nil {
		// No controller should care about orphans being deleted.
		return
	}
	m := tc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Deployment %s deleted.", svc.Name)
	tc.enqueueOpenx(m)
}

func (tc *OpenxController) addHorizontalPodAutoscaler(obj interface{}) {
	hpa := obj.(*autoscalingv2.HorizontalPodAutoscaler)

	if hpa.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		tc.deleteHorizontalPodAutoscaler(hpa)
		return
	}
	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(hpa)
	if controllerRef == nil {
		return
	}
	m := tc.resolveControllerRef(hpa.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Deployment%s namespace:%s added.", hpa.Name, hpa.Namespace)
	tc.enqueueOpenx(m)
}

func (tc *OpenxController) getOpenxsForHorizontalPodAutoscaler(hpa *autoscalingv2.HorizontalPodAutoscaler) []*openxv1.Openx {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: hpa.Labels,
	})
	ms, err := tc.openxLister.Openxes(hpa.Namespace).List(selector)
	if err != nil || len(ms) == 0 {
		return nil
	}
	// Because all ReplicaSet's belonging to a deployment should have a unique label key,
	// there should never be more than one deployment returned by the above method.
	// If that happens we should probably dynamically repair the situation by ultimately
	// trying to clean up one of the controllers, for now we just return th	e older one
	if len(ms) > 1 {
		// ControllerRef will ensure we don't do anything crazy, but more than one
		// item in this list nevertheless constitutes user error.
		zaplogger.Sugar().Infof("user error! more than one openx is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			hpa.Namespace, hpa.Name, hpa.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

func (tc *OpenxController) updateHorizontalPodAutoscaler(old, cur interface{}) {
	curHPA := cur.(*autoscalingv2.HorizontalPodAutoscaler)
	oldHPA := old.(*autoscalingv2.HorizontalPodAutoscaler)
	if curHPA.ResourceVersion == oldHPA.ResourceVersion {
		// Periodic resync will send update events for all known replica sets.
		// Two different versions of the same replica set will always have different RVs.
		return
	}

	curControllerRef := metav1.GetControllerOf(curHPA)
	oldControllerRef := metav1.GetControllerOf(oldHPA)
	controllerRefChanged := !reflect.DeepEqual(curControllerRef, oldControllerRef)
	if controllerRefChanged && oldControllerRef != nil {
		// The ControllerRef was changed. Sync the old controller, if any.
		if d := tc.resolveControllerRef(oldHPA.Namespace, oldControllerRef); d != nil {
			tc.enqueueOpenx(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := tc.resolveControllerRef(curHPA.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("HorizontalPodAutoscaler %s updated.", curHPA.Name)
		tc.enqueueOpenx(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curHPA.Labels, oldHPA.Labels)
	if labelChanged || controllerRefChanged {
		ds := tc.getOpenxsForHorizontalPodAutoscaler(curHPA)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan HorizontalPodAutoscaler %s updated.", curHPA.Name)
		for _, d := range ds {
			tc.enqueueOpenx(d)
		}
	}
}

func (tc *OpenxController) deleteHorizontalPodAutoscaler(obj interface{}) {
	hpa, ok := obj.(*autoscalingv2.HorizontalPodAutoscaler)

	// When a delete is dropped, the relist will notice a pod in the store not
	// in the list, leading to the insertion of a tombstone object which contains
	// the deleted key/value. Note that this value might be stale. If the ReplicaSet
	// changed labels the new deployment will not be woken up till the periodic resync.
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		hpa, ok = tombstone.Obj.(*autoscalingv2.HorizontalPodAutoscaler)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Service %#v", obj))
			return
		}
	}

	controllerRef := metav1.GetControllerOf(hpa)
	if controllerRef == nil {
		// No controller should care about orphans being deleted.
		return
	}
	m := tc.resolveControllerRef(hpa.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("HorizontalPodAutoscaler %s deleted.", hpa.Name)
	tc.enqueueOpenx(m)
}
