package etcd

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	openx "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	openxinformersv1 "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions/openx/v1"
	openxlistersv1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const (
	// maxRetries is the number of times a etcd will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a etcd is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = schema.GroupVersion{Group: openxv1.SchemeGroupVersion.Group, Version: openxv1.SchemeGroupVersion.Version}.WithKind("Etcd")

var controllerName = "etcd"

type EtcdController struct {
	// svcControl is used for adopting/releasing services.
	svcControl controller.ServiceControlInterface
	// stsControl is used for adopting/releasing statefulsets.
	stsControl controller.StatefulSetControlInterface

	client        clientset.Interface
	openx         openx.Interface
	eventRecorder record.EventRecorder

	syncHandler func(key string) error
	enqueueEtcd func(etcd *openxv1.Etcd)

	etcdLister openxlistersv1.EtcdLister
	svcLister  corelisters.ServiceLister
	stsLister  appslisters.StatefulSetLister

	// etcdListerSynced
	etcdListerSynced cache.InformerSynced
	svcListerSynced  cache.InformerSynced
	stsListerSynced  cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

// NewEtcdController creates a new EtcdController.
func NewEtcdController(etcdInformer openxinformersv1.EtcdInformer, svcInformer coreinformers.ServiceInformer, stsInformer appsinformers.StatefulSetInformer, client clientset.Interface, openx openx.Interface) (*EtcdController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events(controller.RecorderNamespace)})

	//if client != nil && client.CoreV1().RESTClient().GetRateLimiter() != nil {
	//if err := ratelimiter.RegisterMetricAndTrackRateLimiterUsage("deployment_controller", client.CoreV1().RESTClient().GetRateLimiter()); err != nil {
	//	return nil, err
	//}
	//}
	ec := &EtcdController{
		client:        client,
		openx:         openx,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "etcd-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	ec.etcdLister = etcdInformer.Lister()
	ec.svcLister = svcInformer.Lister()
	ec.stsLister = stsInformer.Lister()
	ec.etcdListerSynced = etcdInformer.Informer().HasSynced
	ec.svcListerSynced = svcInformer.Informer().HasSynced
	ec.stsListerSynced = stsInformer.Informer().HasSynced

	ec.syncHandler = ec.syncEtcd
	ec.enqueueEtcd = ec.enqueue

	etcdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ec.addEtcd,
		UpdateFunc: ec.updateEtcd,
		DeleteFunc: ec.deleteEtcd,
	})
	svcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ec.addService,
		UpdateFunc: ec.updateService,
		DeleteFunc: ec.deleteService,
	})
	stsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ec.addStatefulSet,
		UpdateFunc: ec.updateStatefulSet,
		DeleteFunc: ec.deleteStatefulSet,
	})

	ec.svcControl = &controller.RealServiceControl{
		KubeClient: client,
		Recorder:   ec.eventRecorder,
	}
	ec.stsControl = &controller.RealStatefulSetControl{
		KubeClient: client,
		Recorder:   ec.eventRecorder,
	}

	return ec, nil
}

// Run begins watching and syncing.
func (ec *EtcdController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer ec.queue.ShutDown()

	zaplogger.Sugar().Infof("Starting etcd controller")
	defer zaplogger.Sugar().Infof("Shutting down etcd controller")

	if !cache.WaitForNamedCacheSync(controllerName, stopCh, ec.etcdListerSynced, ec.svcListerSynced, ec.stsListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(ec.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (ec *EtcdController) worker() {
	for ec.processNextWorkItem() {
	}
}

func (ec *EtcdController) processNextWorkItem() bool {
	key, quit := ec.queue.Get()
	if quit {
		return false
	}
	defer ec.queue.Done(key)

	err := ec.syncHandler(key.(string))
	ec.handleErr(err, key)

	return true
}

func (ec *EtcdController) handleErr(err error, key interface{}) {
	if err == nil || errors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
		ec.queue.Forget(key)
		return
	}

	ns, name, keyErr := cache.SplitMetaNamespaceKey(key.(string))
	if keyErr != nil {
		zaplogger.Sugar().Errorw("Failed to split meta namespace cache key", zap.Any("key", key), zap.Error(err))
	}

	if ec.queue.NumRequeues(key) < maxRetries {
		zaplogger.Sugar().Infow("Error syncing",
			zap.String("controller", controllerName),
			zap.String("namespace", ns),
			zap.String("name", name),
			zap.Error(err))
		ec.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	zaplogger.Sugar().Infof("Dropping etcd %q out of the queue: %v", key, err)
	ec.queue.Forget(key)
}

// resolveControllerRef returns the controller referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching controller
// of the correct Kind.
func (ec *EtcdController) resolveControllerRef(namespace string, controllerRef *metav1.OwnerReference) *openxv1.Etcd {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != controllerKind.Kind {
		return nil
	}
	d, err := ec.etcdLister.Etcds(namespace).Get(controllerRef.Name)
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

func (ec *EtcdController) addEtcd(obj interface{}) {
	d := obj.(*openxv1.Etcd)
	zaplogger.Sugar().Infof("Adding etcd:%s obj: %#v", d.Name, d)
	ec.enqueueEtcd(d)
}

func (ec *EtcdController) updateEtcd(old, cur interface{}) {
	oldD := old.(*openxv1.Etcd)
	curD := cur.(*openxv1.Etcd)
	if curD.ResourceVersion == oldD.ResourceVersion {
		return
	}
	zaplogger.Sugar().Infof("Updating trigger etcd namespace:%s name:%s", oldD.Namespace, oldD.Name)
	ec.enqueueEtcd(curD)
}

func (ec *EtcdController) deleteEtcd(obj interface{}) {
	d, ok := obj.(*openxv1.Etcd)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		d, ok = tombstone.Obj.(*openxv1.Etcd)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a etcd %#v", obj))
			return
		}
	}
	zaplogger.Sugar().Infof("Deleting etcd %s", d.Name)
	ec.enqueueEtcd(d)
}

func (ec *EtcdController) enqueue(etcd *openxv1.Etcd) {
	key, err := controller.KeyFunc(etcd)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", etcd, err))
		return
	}

	ec.queue.Add(key)
}

func (ec *EtcdController) enqueueRateLimited(etcd *openxv1.Etcd) {
	key, err := controller.KeyFunc(etcd)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", etcd, err))
		return
	}

	ec.queue.AddRateLimited(key)
}

// enqueueAfter will enqueue a deployment after the provided amount of time.
func (ec *EtcdController) enqueueAfter(etcd *openxv1.Etcd, after time.Duration) {
	key, err := controller.KeyFunc(etcd)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", etcd, err))
		return
	}

	ec.queue.AddAfter(key, after)
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

func (ec *EtcdController) getServicesForEtcd(m *openxv1.Etcd) ([]*corev1.Service, error) {
	// List all Services to find those we own but that no longer match our
	// selector. They will be orphaned by ClaimServices().
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd %s/%s has invalid label selector: %v", m.Namespace, m.Name, err)
	}
	svcList, err := ec.svcLister.Services(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := ec.openx.OpenxV1().Etcds(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original etcd %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewServiceControllerRefManager(ec.svcControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimServices(svcList)
}

func (ec *EtcdController) getStatefulSetsForEtcd(m *openxv1.Etcd) ([]*appsv1.StatefulSet, error) {
	// Get all Pods that potentially belong to this Deployment.
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, err
	}
	stsList, err := ec.stsLister.StatefulSets(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := ec.openx.OpenxV1().Etcds(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original etcd %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewStatefulSetControllerRefManager(ec.stsControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimStatefulSets(stsList)
}

// syncEtcd will sync the etcd with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (ec *EtcdController) syncEtcd(key string) error {
	startTime := time.Now()
	zaplogger.Sugar().Infof("Started syncing etcd %q (%v)", key, startTime)
	defer func() {
		zaplogger.Sugar().Infof("Finished syncing etcd %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	etcd, err := ec.etcdLister.Etcds(namespace).Get(name)
	if errors.IsNotFound(err) {
		zaplogger.Sugar().Infof("Deployment %v has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	m := etcd.DeepCopy()
	m.Labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: m.Name,
	}

	// List Services owned by this Etcd, while reconciling ControllerRef
	// through adoption/orphaning.
	svcList, err := ec.getServicesForEtcd(m)
	if err != nil {
		return err
	}
	// List all Pods owned by this Deployment, grouped by their ReplicaSet.
	// Current uses of the podMap are:
	//
	// * check if a Pod is labeled correctly with the pod-template-hash label.
	// * check that no old Pods are running in the middle of Recreate Deployments.
	stsList, err := ec.getStatefulSetsForEtcd(m)
	if err != nil {
		return err
	}

	if m.DeletionTimestamp != nil {
		return ec.syncStatusOnly(m, stsList)
	}

	// compare and update svc
	if err = ec.syncServices(m, svcList); err != nil {
		zaplogger.Sugar().Error(err)
		ec.eventRecorder.Event(etcd.Recorder(), corev1.EventTypeWarning, controller.EtcdSyncAborted, err.Error())
		return err
	}

	// compare and update sts
	if err = ec.syncStatefulSets(m, stsList); err != nil {
		ec.eventRecorder.Event(etcd.Recorder(), corev1.EventTypeWarning, controller.EtcdSyncAborted, err.Error())
		return err
	}
	return ec.syncStatusOnly(m, stsList)
}

// addService enqueues the etcd that manages a Service when the Service is created.
func (ec *EtcdController) addService(obj interface{}) {
	svc := obj.(*corev1.Service)

	if svc.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		ec.deleteService(svc)
		return
	}

	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(svc)
	if controllerRef == nil {
		//zaplogger.Sugar().Infof("Orphan Service:%s namespace:%s", svc.Name, svc.Namespace)
		return
	}
	//zaplogger.Sugar().Infof("Service:%s namespace:%s controllerRef:%#v", svc.Name, svc.Namespace, controllerRef)
	m := ec.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Service %s namespace:%s added.", svc.Name, svc.Namespace)
	ec.enqueueEtcd(m)
}

// getDeploymentsForReplicaSet returns a list of Deployments that potentially
// match a ReplicaSet.
func (ec *EtcdController) getEtcdsForService(svc *corev1.Service) []*openxv1.Etcd {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ms, err := ec.etcdLister.Etcds(svc.Namespace).List(selector)
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
		zaplogger.Sugar().Infof("user error! more than one etcd is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

// updateService figures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (ec *EtcdController) updateService(old, cur interface{}) {
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
		if d := ec.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			ec.enqueueEtcd(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := ec.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("Service %s updated.", curRS.Name)
		ec.enqueueEtcd(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := ec.getEtcdsForService(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan Service %s updated.", curRS.Name)
		for _, d := range ds {
			ec.enqueueEtcd(d)
		}
	}
}

// deleteReplicaSet enqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (ec *EtcdController) deleteService(obj interface{}) {
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
	m := ec.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Service %s deleted.", svc.Name)
	ec.enqueueEtcd(m)
}

// addStatefulSet enqueues the etcd that manages a Service when the StatefulSet is created.
func (ec *EtcdController) addStatefulSet(obj interface{}) {
	sts := obj.(*appsv1.StatefulSet)

	if sts.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		ec.deleteService(sts)
		return
	}
	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(sts)
	if controllerRef == nil {
		return
	}
	m := ec.resolveControllerRef(sts.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("StatefulSet %s namespace:%s added.", sts.Name, sts.Namespace)
	ec.enqueueEtcd(m)
}

// getEtcdsForStatefulSet returns a list of Deployments that potentially
// match a ReplicaSet.
func (ec *EtcdController) getEtcdsForStatefulSet(svc *appsv1.StatefulSet) []*openxv1.Etcd {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ms, err := ec.etcdLister.Etcds(svc.Namespace).List(selector)
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
		zaplogger.Sugar().Infof("user error! more than one etcd is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

// updateStatefulSet figures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (ec *EtcdController) updateStatefulSet(old, cur interface{}) {
	curRS := cur.(*appsv1.StatefulSet)
	oldRS := old.(*appsv1.StatefulSet)
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
		if d := ec.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			ec.enqueueEtcd(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := ec.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("StatefulSet %s updated.", curRS.Name)
		ec.enqueueEtcd(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := ec.getEtcdsForStatefulSet(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan StatefulSet %s updated.", curRS.Name)
		for _, d := range ds {
			ec.enqueueEtcd(d)
		}
	}
}

// deleteStatefulSet enqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (ec *EtcdController) deleteStatefulSet(obj interface{}) {
	svc, ok := obj.(*appsv1.StatefulSet)

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
		svc, ok = tombstone.Obj.(*appsv1.StatefulSet)
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
	m := ec.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("StatefulSet %s deleted.", svc.Name)
	ec.enqueueEtcd(m)
}
