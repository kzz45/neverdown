package loadbalancer

import (
	"fmt"
	"strconv"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	openx "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	openxinformersv1 "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions/openx/v1"
	openxlistersv1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

var controllerName = "loadbalancer"

type Interface interface {
	Annotations(svcType corev1.ServiceType, cnc openxv1.CloudNetworkConfig, namespace string) (map[string]string, error)
}

// LoadBalancerController is responsible for synchronizing mysql objects stored
// in the system with actual running replica sets and pods.
type LoadBalancerController struct {
	namespace string

	openx         openx.Interface
	eventRecorder record.EventRecorder

	lbLister openxlistersv1.LoadBalancerLister
	acLister openxlistersv1.AccessControlLister

	lbListerSynced cache.InformerSynced
	acListerSynced cache.InformerSynced

	// mysqls that need to be synced
	queue workqueue.RateLimitingInterface
}

func NewLoadBalancerController(lbInformer openxinformersv1.LoadBalancerInformer, acInformer openxinformersv1.AccessControlInformer, openx openx.Interface) (*LoadBalancerController, error) {
	//eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartStructuredLogging(0)
	//eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events("")})

	lbc := &LoadBalancerController{
		namespace: "kube-neverdown",
		openx:     openx,
		//eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "loadbalancer-controller"}),
		//queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	lbc.lbLister = lbInformer.Lister()
	lbc.acLister = acInformer.Lister()
	lbc.lbListerSynced = lbInformer.Informer().HasSynced
	lbc.acListerSynced = acInformer.Informer().HasSynced

	//lbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    lbc.addMysql,
	//	UpdateFunc: lbc.updateMysql,
	//	DeleteFunc: lbc.deleteMysql,
	//})
	//acInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    lbc.addService,
	//	UpdateFunc: lbc.updateService,
	//	DeleteFunc: lbc.deleteService,
	//})

	return lbc, nil
}

func (lbc *LoadBalancerController) Annotations(svcType corev1.ServiceType, cnc openxv1.CloudNetworkConfig, namespace string) (map[string]string, error) {
	if svcType != corev1.ServiceTypeLoadBalancer {
		return map[string]string{}, nil
	}
	if cnc.CloudSLB != nil {
		return lbc.aliyun(cnc.CloudSLB, namespace)
	}
	return map[string]string{}, nil
}

func (lbc *LoadBalancerController) aliyun(slb *openxv1.CloudSLB, namespace string) (map[string]string, error) {
	annotations := make(map[string]string)
	if slb.LoadBalancerId == "" {
		return nil, fmt.Errorf("AliyunSLB LoadBalancerId must be specified")
	}
	lb, err := lbc.lbLister.LoadBalancers(namespace).Get(slb.LoadBalancerId)
	if err != nil {
		return nil, err
	}
	if lb.Spec.Instance.Key != "" {
		annotations[lb.Spec.Instance.Key] = lb.Spec.Instance.Value
	}
	if lb.Spec.OverrideListeners.Key != "" {
		annotations[lb.Spec.OverrideListeners.Key] = strconv.FormatBool(slb.OverrideListeners)
	}
	if slb.Status == "" || slb.Status == openxv1.CloudLoadBalancerOff {
		return annotations, nil
	}
	if slb.AccessControlId == "" {
		return nil, fmt.Errorf("SLB AccessControlId must be specified when the Staus is `on`")
	}
	ac, err := lbc.acLister.AccessControls(namespace).Get(slb.AccessControlId)
	if err != nil {
		return nil, err
	}
	if ac.Spec.Instance.Key == "" {
		return nil, fmt.Errorf("SLB AccessControl spec Instance key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Instance.Key] = ac.Spec.Instance.Value
	if ac.Spec.Status.Key == "" {
		return nil, fmt.Errorf("SLB AccessControl spec Status key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Status.Key] = string(slb.Status)
	if ac.Spec.Type.Key == "" {
		return nil, fmt.Errorf("SLB AccessControl spec Type key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Type.Key] = ac.Spec.Type.Value
	return annotations, nil
}
