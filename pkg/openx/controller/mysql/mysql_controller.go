package mysql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	openx "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	openxinformersv1 "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions/openx/v1"
	openxlistersv1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"
	"github.com/kzz45/neverdown/pkg/openx/controller/loadbalancer"

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
	// maxRetries is the number of times a mysql will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a mysql is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = schema.GroupVersion{Group: openxv1.GroupVersion.Group, Version: openxv1.GroupVersion.Version}.WithKind("Mysql")

var LabelsClusterRoleKey = "mysqlclusterrole"

var controllerName = "mysql"

// MysqlController is responsible for synchronizing mysql objects stored
// in the system with actual running replica sets and pods.
type MysqlController struct {
	loadbalancer loadbalancer.Interface

	// svcControl is used for adopting/releasing services.
	svcControl controller.ServiceControlInterface
	// stsControl is used for adopting/releasing statefulsets.
	stsControl controller.StatefulSetControlInterface

	client        clientset.Interface
	openx         openx.Interface
	eventRecorder record.EventRecorder

	// To allow injection of syncMysql for testing.
	syncHandler func(key string) error
	// used for unit testing
	enqueueMysql func(mysql *openxv1.Mysql)

	// mysqlLister can list/get mysql from the shared informer's store
	mysqlLister openxlistersv1.MysqlLister
	// svcLister can list/get services from the shared informer's store
	svcLister corelisters.ServiceLister
	// stsLister can list/get mysqls from the shared informer's store
	stsLister appslisters.StatefulSetLister

	// mysqlListerSynced
	mysqlListerSynced cache.InformerSynced
	// svcListerSynced returns true if the svc store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	svcListerSynced cache.InformerSynced
	// stsListerSynced returns true if the mysql store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	stsListerSynced cache.InformerSynced

	// mysqls that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewMysqlController creates a new MysqlController.
func NewMysqlController(loadbalancher loadbalancer.Interface, mysqlInformer openxinformersv1.MysqlInformer, svcInformer coreinformers.ServiceInformer, stsInformer appsinformers.StatefulSetInformer, client clientset.Interface, openx openx.Interface) (*MysqlController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events(controller.RecorderNamespace)})

	//if client != nil && client.CoreV1().RESTClient().GetRateLimiter() != nil {
	//if err := ratelimiter.RegisterMetricAndTrackRateLimiterUsage("deployment_controller", client.CoreV1().RESTClient().GetRateLimiter()); err != nil {
	//	return nil, err
	//}
	//}
	mc := &MysqlController{
		loadbalancer:  loadbalancher,
		client:        client,
		openx:         openx,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "mysql-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	mc.mysqlLister = mysqlInformer.Lister()
	mc.svcLister = svcInformer.Lister()
	mc.stsLister = stsInformer.Lister()
	mc.mysqlListerSynced = mysqlInformer.Informer().HasSynced
	mc.svcListerSynced = svcInformer.Informer().HasSynced
	mc.stsListerSynced = stsInformer.Informer().HasSynced

	mc.syncHandler = mc.syncMysql
	mc.enqueueMysql = mc.enqueue

	mysqlInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mc.addMysql,
		UpdateFunc: mc.updateMysql,
		// This will enter the sync loop and no-op, because the mysql has been deleted from the store.
		DeleteFunc: mc.deleteMysql,
	})
	svcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mc.addService,
		UpdateFunc: mc.updateService,
		DeleteFunc: mc.deleteService,
	})
	stsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    mc.addStatefulSet,
		UpdateFunc: mc.updateStatefulSet,
		DeleteFunc: mc.deleteStatefulSet,
	})

	mc.svcControl = &controller.RealServiceControl{
		KubeClient: client,
		Recorder:   mc.eventRecorder,
	}
	mc.stsControl = &controller.RealStatefulSetControl{
		KubeClient: client,
		Recorder:   mc.eventRecorder,
	}
	return mc, nil
}

// Run begins watching and syncing.
func (mc *MysqlController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer mc.queue.ShutDown()

	zaplogger.Sugar().Infof("Starting mysql controller")
	defer zaplogger.Sugar().Infof("Shutting down mysql controller")

	if !cache.WaitForNamedCacheSync(controllerName, stopCh, mc.mysqlListerSynced, mc.svcListerSynced, mc.stsListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(mc.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (mc *MysqlController) worker() {
	for mc.processNextWorkItem() {
	}
}

func (mc *MysqlController) processNextWorkItem() bool {
	key, quit := mc.queue.Get()
	if quit {
		return false
	}
	defer mc.queue.Done(key)

	err := mc.syncHandler(key.(string))
	mc.handleErr(err, key)

	return true
}

func (mc *MysqlController) handleErr(err error, key interface{}) {
	if err == nil || errors.HasStatusCause(err, corev1.NamespaceTerminatingCause) {
		mc.queue.Forget(key)
		return
	}

	ns, name, keyErr := cache.SplitMetaNamespaceKey(key.(string))
	if keyErr != nil {
		zaplogger.Sugar().Errorw("Failed to split meta namespace cache key", zap.Any("key", key), zap.Error(err))
	}

	if mc.queue.NumRequeues(key) < maxRetries {
		zaplogger.Sugar().Infow("Error syncing",
			zap.String("controller", controllerName),
			zap.String("namespace", ns),
			zap.String("name", name),
			zap.Error(err))
		mc.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	zaplogger.Sugar().Infof("Dropping mysql %q out of the queue: %v", key, err)
	mc.queue.Forget(key)
}

// resolveControllerRef returns the controller referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching controller
// of the correct Kind.
func (mc *MysqlController) resolveControllerRef(namespace string, controllerRef *metav1.OwnerReference) *openxv1.Mysql {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != controllerKind.Kind {
		return nil
	}
	d, err := mc.mysqlLister.Mysqls(namespace).Get(controllerRef.Name)
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

func (mc *MysqlController) addMysql(obj interface{}) {
	d := obj.(*openxv1.Mysql)
	zaplogger.Sugar().Infof("Adding mysql:%s obj: %#v", d.Name, d)
	mc.enqueueMysql(d)
}

func (mc *MysqlController) updateMysql(old, cur interface{}) {
	oldD := old.(*openxv1.Mysql)
	curD := cur.(*openxv1.Mysql)
	if curD.ResourceVersion == oldD.ResourceVersion {
		return
	}
	zaplogger.Sugar().Infof("Updating trigger mysql namespace:%s name:%s", oldD.Namespace, oldD.Name)
	mc.enqueueMysql(curD)
}

func (mc *MysqlController) deleteMysql(obj interface{}) {
	d, ok := obj.(*openxv1.Mysql)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		d, ok = tombstone.Obj.(*openxv1.Mysql)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Mysql %#v", obj))
			return
		}
	}
	zaplogger.Sugar().Infof("Deleting mysql %s", d.Name)
	mc.enqueueMysql(d)
}

func (mc *MysqlController) enqueue(mysql *openxv1.Mysql) {
	key, err := controller.KeyFunc(mysql)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", mysql, err))
		return
	}

	mc.queue.Add(key)
}

func (mc *MysqlController) enqueueRateLimited(mysql *openxv1.Mysql) {
	key, err := controller.KeyFunc(mysql)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", mysql, err))
		return
	}

	mc.queue.AddRateLimited(key)
}

// enqueueAfter will enqueue a deployment after the provided amount of time.
func (mc *MysqlController) enqueueAfter(mysql *openxv1.Mysql, after time.Duration) {
	key, err := controller.KeyFunc(mysql)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", mysql, err))
		return
	}

	mc.queue.AddAfter(key, after)
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

// getServicesForMysql uses ControllerRefManager to reconcile
// ControllerRef by adopting and orphaning.
// It returns the list of Services that this Mysql should manage.
func (mc *MysqlController) getServicesForMysql(m *openxv1.Mysql) ([]*corev1.Service, error) {
	// List all Services to find those we own but that no longer match our
	// selector. They will be orphaned by ClaimServices().
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql %s/%s has invalid label selector: %v", m.Namespace, m.Name, err)
	}
	svcList, err := mc.svcLister.Services(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := mc.openx.OpenxV1().Mysqls(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original Mysql %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewServiceControllerRefManager(mc.svcControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimServices(svcList)
}

// getStatefulSetsForMysql returns the StatefulSet managed by a Mysql.
//
// It returns a map from ReplicaSet UID to a list of Pods controlled by that RS,
// according to the Pod's ControllerRef.
// NOTE: The pod pointers returned by this method point the pod objects in the cache and thus
// shouldn't be modified in any way.
func (mc *MysqlController) getStatefulSetsForMysql(m *openxv1.Mysql) ([]*appsv1.StatefulSet, error) {
	// Get all Pods that potentially belong to this Deployment.
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: m.Labels,
	})
	if err != nil {
		return nil, err
	}
	stsList, err := mc.stsLister.StatefulSets(m.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing ReplicaSets (see #42639).
	canAdoptFunc := RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := mc.openx.OpenxV1().Mysqls(m.Namespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != m.UID {
			return nil, fmt.Errorf("original Mysql %v/%v is gone: got uid %v, wanted %v", m.Namespace, m.Name, fresh.UID, m.UID)
		}
		return fresh, nil
	})
	cm := controller.NewStatefulSetControllerRefManager(mc.stsControl, m, selector, controllerKind, canAdoptFunc)
	return cm.ClaimStatefulSets(stsList)
}

// syncMysql will sync the mysql with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (mc *MysqlController) syncMysql(key string) error {
	startTime := time.Now()
	zaplogger.Sugar().Infof("Started syncing mysql %q (%v)", key, startTime)
	defer func() {
		zaplogger.Sugar().Infof("Finished syncing mysql %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	mysql, err := mc.mysqlLister.Mysqls(namespace).Get(name)
	if errors.IsNotFound(err) {
		zaplogger.Sugar().Infof("Deployment %v has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	m := mysql.DeepCopy()
	m.Labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: m.Name,
	}

	//everything := metav1.LabelSelector{}
	//if reflect.DeepEqual(d.Spec.Selector, &everything) {
	//	mc.eventRecorder.Eventf(d, corev1.EventTypeWarning, "SelectingAll", "This deployment is selecting all pods. A non-empty selector is required.")
	//	if d.Status.ObservedGeneration < d.Generation {
	//		d.Status.ObservedGeneration = d.Generation
	//		mc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(context.TODO(), d, metav1.UpdateOptions{})
	//	}
	//	return nil
	//}

	// List Services owned by this Mysql, while reconciling ControllerRef
	// through adoption/orphaning.
	svcList, err := mc.getServicesForMysql(m)
	if err != nil {
		return err
	}
	// List all Pods owned by this Deployment, grouped by their ReplicaSet.
	// Current uses of the podMap are:
	//
	// * check if a Pod is labeled correctly with the pod-template-hash label.
	// * check that no old Pods are running in the middle of Recreate Deployments.
	stsList, err := mc.getStatefulSetsForMysql(m)
	if err != nil {
		return err
	}

	if m.DeletionTimestamp != nil {
		return mc.syncStatusOnly(m, stsList)
	}

	// compare and update svc
	if err = mc.syncServices(m, svcList); err != nil {
		zaplogger.Sugar().Error(err)
		mc.eventRecorder.Event(mysql.Recorder(), corev1.EventTypeWarning, controller.MysqlSyncAborted, err.Error())
		return err
	}

	// compare and update sts
	if err = mc.syncStatefulSets(m, stsList); err != nil {
		mc.eventRecorder.Event(mysql.Recorder(), corev1.EventTypeWarning, controller.MysqlSyncAborted, err.Error())
		return err
	}
	return mc.syncStatusOnly(m, stsList)

	//// Update deployment conditions with an Unknown condition when pausing/resuming
	//// a deployment. In this way, we can be sure that we won't timeout when a user
	//// resumes a Deployment with a set progressDeadlineSeconds.
	//if err = mc.checkPausedConditions(d); err != nil {
	//	return err
	//}
	//
	//if d.Spec.Paused {
	//	return mc.sync(d, rsList)
	//}
	//
	//// rollback is not re-entrant in case the underlying replica sets are updated with a new
	//// revision so we should ensure that we won't proceed to update replica sets until we
	//// make sure that the deployment has cleaned up its rollback spec in subsequent enqueues.
	//if getRollbackTo(d) != nil {
	//	return mc.rollback(d, rsList)
	//}
	//
	//scalingEvent, err := mc.isScalingEvent(d, rsList)
	//if err != nil {
	//	return err
	//}
	//if scalingEvent {
	//	return mc.sync(d, rsList)
	//}
	//
	//switch d.Spec.Strategy.Type {
	//case apps.RecreateDeploymentStrategyType:
	//	return mc.rolloutRecreate(d, rsList, podMap)
	//case apps.RollingUpdateDeploymentStrategyType:
	//	return mc.rolloutRolling(d, rsList)
	//}
	//return fmt.Errorf("unexpected deployment strategy type: %s", d.Spec.Strategy.Type)
}

// addService enqueues the mysql that manages a Service when the Service is created.
func (mc *MysqlController) addService(obj interface{}) {
	svc := obj.(*corev1.Service)

	if svc.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		mc.deleteService(svc)
		return
	}

	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(svc)
	if controllerRef == nil {
		//zaplogger.Sugar().Infof("Orphan Service:%s namespace:%s", svc.Name, svc.Namespace)
		return
	}
	//zaplogger.Sugar().Infof("Service:%s namespace:%s controllerRef:%#v", svc.Name, svc.Namespace, controllerRef)
	m := mc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Service %s namespace:%s added.", svc.Name, svc.Namespace)
	mc.enqueueMysql(m)
}

// getDeploymentsForReplicaSet returns a list of Deployments that potentially
// match a ReplicaSet.
func (mc *MysqlController) getMysqlsForService(svc *corev1.Service) []*openxv1.Mysql {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ms, err := mc.mysqlLister.Mysqls(svc.Namespace).List(selector)
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
		zaplogger.Sugar().Infof("user error! more than one mysql is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

// updateService figures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (mc *MysqlController) updateService(old, cur interface{}) {
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
		if d := mc.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			mc.enqueueMysql(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := mc.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("Service %s updated.", curRS.Name)
		mc.enqueueMysql(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := mc.getMysqlsForService(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan Service %s updated.", curRS.Name)
		for _, d := range ds {
			mc.enqueueMysql(d)
		}
	}
}

// deleteReplicaSet enqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (mc *MysqlController) deleteService(obj interface{}) {
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
	m := mc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("Service %s deleted.", svc.Name)
	mc.enqueueMysql(m)
}

// addStatefulSet enqueues the mysql that manages a Service when the StatefulSet is created.
func (mc *MysqlController) addStatefulSet(obj interface{}) {
	sts := obj.(*appsv1.StatefulSet)

	if sts.DeletionTimestamp != nil {
		// On a restart of the controller manager, it's possible for an object to
		// show up in a state that is already pending deletion.
		mc.deleteService(sts)
		return
	}
	// If it has a ControllerRef, that's all that matters.
	controllerRef := metav1.GetControllerOf(sts)
	if controllerRef == nil {
		return
	}
	m := mc.resolveControllerRef(sts.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("StatefulSet %s namespace:%s added.", sts.Name, sts.Namespace)
	mc.enqueueMysql(m)
}

// getMysqlsForStatefulSet returns a list of Deployments that potentially
// match a ReplicaSet.
func (mc *MysqlController) getMysqlsForStatefulSet(svc *appsv1.StatefulSet) []*openxv1.Mysql {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Labels,
	})
	ms, err := mc.mysqlLister.Mysqls(svc.Namespace).List(selector)
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
		zaplogger.Sugar().Infof("user error! more than one mysql is selecting replica set %s/%s with labels: %#v, returning %s/%s",
			svc.Namespace, svc.Name, svc.Labels, ms[0].Namespace, ms[0].Name)
	}
	return ms
}

// updateStatefulSet figures out what deployment(s) manage a ReplicaSet when the ReplicaSet
// is updated and wake them up. If the anything of the ReplicaSets have changed, we need to
// awaken both the old and new deployments. old and cur must be *apps.ReplicaSet
// types.
func (mc *MysqlController) updateStatefulSet(old, cur interface{}) {
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
		if d := mc.resolveControllerRef(oldRS.Namespace, oldControllerRef); d != nil {
			mc.enqueueMysql(d)
		}
	}

	// If it has a ControllerRef, that's all that matters.
	if curControllerRef != nil {
		d := mc.resolveControllerRef(curRS.Namespace, curControllerRef)
		if d == nil {
			return
		}
		zaplogger.Sugar().Infof("StatefulSet %s updated.", curRS.Name)
		mc.enqueueMysql(d)
		return
	}

	// Otherwise, it's an orphan. If anything changed, sync matching controllers
	// to see if anyone wants to adopt it now.
	labelChanged := !reflect.DeepEqual(curRS.Labels, oldRS.Labels)
	if labelChanged || controllerRefChanged {
		ds := mc.getMysqlsForStatefulSet(curRS)
		if len(ds) == 0 {
			return
		}
		zaplogger.Sugar().Infof("Orphan StatefulSet %s updated.", curRS.Name)
		for _, d := range ds {
			mc.enqueueMysql(d)
		}
	}
}

// deleteStatefulSet enqueues the deployment that manages a ReplicaSet when
// the ReplicaSet is deleted. obj could be an *apps.ReplicaSet, or
// a DeletionFinalStateUnknown marker item.
func (mc *MysqlController) deleteStatefulSet(obj interface{}) {
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
	m := mc.resolveControllerRef(svc.Namespace, controllerRef)
	if m == nil {
		return
	}
	zaplogger.Sugar().Infof("StatefulSet %s deleted.", svc.Name)
	mc.enqueueMysql(m)
}
