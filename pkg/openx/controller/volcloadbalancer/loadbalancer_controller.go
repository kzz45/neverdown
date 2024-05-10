package volcloadbalancer

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

var controllerName = "volcloadbalancer"

type Interface interface {
	Annotations(svcType corev1.ServiceType, cnc openxv1.CloudNetworkConfig, namespace string) (map[string]string, error)
}

// VolcLoadBalancerController is responsible for synchronizing mysql objects stored
// in the system with actual running replica sets and pods.
type VolcLoadBalancerController struct {
	namespace string

	openx         openx.Interface
	eventRecorder record.EventRecorder

	lbLister openxlistersv1.VolcLoadBalancerLister
	acLister openxlistersv1.VolcAccessControlLister

	lbListerSynced cache.InformerSynced
	acListerSynced cache.InformerSynced

	// mysqls that need to be synced
	queue workqueue.RateLimitingInterface
}

func NewVolcLoadBalancerController(lbInformer openxinformersv1.VolcLoadBalancerInformer, acInformer openxinformersv1.VolcAccessControlInformer, openx openx.Interface) (*VolcLoadBalancerController, error) {
	//eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartStructuredLogging(0)
	//eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events("")})

	lbc := &VolcLoadBalancerController{
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

func (lbc *VolcLoadBalancerController) Annotations(svcType corev1.ServiceType, cnc openxv1.CloudNetworkConfig, namespace string) (map[string]string, error) {
	if svcType != corev1.ServiceTypeLoadBalancer {
		return map[string]string{}, nil
	}
	if cnc.AliyunSLB != nil {
		return lbc.volc(cnc.AliyunSLB, namespace)
	}
	return map[string]string{}, nil
}

func (lbc *VolcLoadBalancerController) volc(slb *openxv1.AliyunSLB, namespace string) (map[string]string, error) {
	annotations := make(map[string]string)
	if slb.LoadBalancerId == "" {
		return nil, fmt.Errorf("VolcSLB LoadBalancerId must be specified")
	}
	lb, err := lbc.lbLister.VolcLoadBalancers(namespace).Get(slb.LoadBalancerId)
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
		return nil, fmt.Errorf("VolcSLB AccessControlId must be specified when the Staus is `on`")
	}
	ac, err := lbc.acLister.VolcAccessControls(namespace).Get(slb.AccessControlId)
	if err != nil {
		return nil, err
	}
	if ac.Spec.Instance.Key == "" {
		return nil, fmt.Errorf("VolcSLB AccessControl spec Instance key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Instance.Key] = ac.Spec.Instance.Value
	if ac.Spec.Status.Key == "" {
		return nil, fmt.Errorf("VolcSLB AccessControl spec Status key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Status.Key] = string(slb.Status)
	if ac.Spec.Type.Key == "" {
		return nil, fmt.Errorf("VolcSLB AccessControl spec Type key was nil name:%s", slb.AccessControlId)
	}
	annotations[ac.Spec.Type.Key] = ac.Spec.Type.Value
	return annotations, nil
}
