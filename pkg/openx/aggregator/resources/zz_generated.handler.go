package resources

import (
	"context"
	"fmt"
	"reflect"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

const ErrGVKNotExist = "unable to match GroupVersionKind:%v"

func (r *Resources) registerInformer() {
	r.InformerFactory.Apps().V1().DaemonSets().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Apps().V1().Deployments().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Apps().V1().StatefulSets().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Autoscaling().V2().HorizontalPodAutoscalers().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().ConfigMaps().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Endpoints().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Events().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Namespaces().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().PersistentVolumes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().PersistentVolumeClaims().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Secrets().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().Services().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Core().V1().ServiceAccounts().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Rbac().V1().ClusterRoles().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.InformerFactory.Rbac().V1().ClusterRoleBindings().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Affinities().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().AliyunAccessControls().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().VolcAccessControls().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().AliyunLoadBalancers().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().VolcLoadBalancers().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Etcds().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Mysqls().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().NodeSelectors().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Redises().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Openxes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
	r.OpenXInformerFactory.Openx().V1().Tolerations().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.addFunc,
		UpdateFunc: r.updateFunc,
		DeleteFunc: r.deleteFunc,
	})
}

func (r *Resources) create(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	createOpts := v1.CreateOptions{}
	switch *gvk {
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}:
		obj := &appsv1.DaemonSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().DaemonSets(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}:
		obj := &appsv1.Deployment{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().Deployments(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}:
		obj := &appsv1.StatefulSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().StatefulSets(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: autoscalingv2.GroupName, Version: autoscalingv2.SchemeGroupVersion.Version, Kind: "HorizontalPodAutoscaler"}:
		obj := &autoscalingv2.HorizontalPodAutoscaler{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ConfigMap"}:
		obj := &corev1.ConfigMap{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().ConfigMaps(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Endpoints"}:
		obj := &corev1.Endpoints{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Endpoints(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &corev1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Events(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Namespace"}:
		obj := &corev1.Namespace{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Namespaces().Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Node"}:
		obj := &corev1.Node{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Nodes().Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolume"}:
		obj := &corev1.PersistentVolume{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().PersistentVolumes().Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolumeClaim"}:
		obj := &corev1.PersistentVolumeClaim{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}:
		obj := &corev1.Pod{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Pods(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Secret"}:
		obj := &corev1.Secret{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Secrets(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Service"}:
		obj := &corev1.Service{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Services(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ServiceAccount"}:
		obj := &corev1.ServiceAccount{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().ServiceAccounts(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRole"}:
		obj := &rbacv1.ClusterRole{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.RbacV1().ClusterRoles().Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRoleBinding"}:
		obj := &rbacv1.ClusterRoleBinding{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.RbacV1().ClusterRoleBindings().Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Affinity"}:
		obj := &openxv1.Affinity{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Affinities(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunAccessControl"}:
		obj := &openxv1.AliyunAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().AliyunAccessControls(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcAccessControl"}:
		obj := &openxv1.VolcAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().VolcAccessControls(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunLoadBalancer"}:
		obj := &openxv1.AliyunLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().AliyunLoadBalancers(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcLoadBalancer"}:
		obj := &openxv1.VolcLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().VolcLoadBalancers(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Etcd"}:
		obj := &openxv1.Etcd{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Etcds(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Mysql"}:
		obj := &openxv1.Mysql{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Mysqls(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "NodeSelector"}:
		obj := &openxv1.NodeSelector{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().NodeSelectors(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Redis"}:
		obj := &openxv1.Redis{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Redises(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Openx"}:
		obj := &openxv1.Openx{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Openxes(namespace).Create(ctx, obj, createOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Toleration"}:
		obj := &openxv1.Toleration{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Tolerations(namespace).Create(ctx, obj, createOpts)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (r *Resources) delete(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	delOpts := v1.DeleteOptions{}
	switch *gvk {
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}:
		obj := &appsv1.DaemonSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.AppsV1().DaemonSets(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}:
		obj := &appsv1.Deployment{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.AppsV1().Deployments(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}:
		obj := &appsv1.StatefulSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.AppsV1().StatefulSets(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: autoscalingv2.GroupName, Version: autoscalingv2.SchemeGroupVersion.Version, Kind: "HorizontalPodAutoscaler"}:
		obj := &autoscalingv2.HorizontalPodAutoscaler{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ConfigMap"}:
		obj := &corev1.ConfigMap{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().ConfigMaps(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Endpoints"}:
		obj := &corev1.Endpoints{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Endpoints(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &corev1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Events(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Namespace"}:
		obj := &corev1.Namespace{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Namespaces().Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Node"}:
		obj := &corev1.Node{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Nodes().Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolume"}:
		obj := &corev1.PersistentVolume{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().PersistentVolumes().Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolumeClaim"}:
		obj := &corev1.PersistentVolumeClaim{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}:
		obj := &corev1.Pod{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Pods(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Secret"}:
		obj := &corev1.Secret{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Secrets(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Service"}:
		obj := &corev1.Service{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().Services(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ServiceAccount"}:
		obj := &corev1.ServiceAccount{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.CoreV1().ServiceAccounts(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRole"}:
		obj := &rbacv1.ClusterRole{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.RbacV1().ClusterRoles().Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRoleBinding"}:
		obj := &rbacv1.ClusterRoleBinding{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.nativeClientSet.RbacV1().ClusterRoleBindings().Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Affinity"}:
		obj := &openxv1.Affinity{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Affinities(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunAccessControl"}:
		obj := &openxv1.AliyunAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().AliyunAccessControls(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcAccessControl"}:
		obj := &openxv1.VolcAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().VolcAccessControls(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunLoadBalancer"}:
		obj := &openxv1.AliyunLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().AliyunLoadBalancers(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcLoadBalancer"}:
		obj := &openxv1.VolcLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().VolcLoadBalancers(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Etcd"}:
		obj := &openxv1.Etcd{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Etcds(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Mysql"}:
		obj := &openxv1.Mysql{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Mysqls(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "NodeSelector"}:
		obj := &openxv1.NodeSelector{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().NodeSelectors(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Redis"}:
		obj := &openxv1.Redis{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Redises(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Openx"}:
		obj := &openxv1.Openx{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Openxes(namespace).Delete(ctx, obj.Name, delOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Toleration"}:
		obj := &openxv1.Toleration{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		err = r.openxClientSet.OpenxV1().Tolerations(namespace).Delete(ctx, obj.Name, delOpts)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (r *Resources) update(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	updateOpts := v1.UpdateOptions{}
	switch *gvk {
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}:
		obj := &appsv1.DaemonSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().DaemonSets(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}:
		obj := &appsv1.Deployment{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().Deployments(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}:
		obj := &appsv1.StatefulSet{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AppsV1().StatefulSets(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: autoscalingv2.GroupName, Version: autoscalingv2.SchemeGroupVersion.Version, Kind: "HorizontalPodAutoscaler"}:
		obj := &autoscalingv2.HorizontalPodAutoscaler{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ConfigMap"}:
		obj := &corev1.ConfigMap{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().ConfigMaps(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Endpoints"}:
		obj := &corev1.Endpoints{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Endpoints(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &corev1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Events(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Namespace"}:
		obj := &corev1.Namespace{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Namespaces().Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Node"}:
		obj := &corev1.Node{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Nodes().Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolume"}:
		obj := &corev1.PersistentVolume{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().PersistentVolumes().Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolumeClaim"}:
		obj := &corev1.PersistentVolumeClaim{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}:
		obj := &corev1.Pod{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Pods(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Secret"}:
		obj := &corev1.Secret{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Secrets(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Service"}:
		obj := &corev1.Service{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().Services(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ServiceAccount"}:
		obj := &corev1.ServiceAccount{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.CoreV1().ServiceAccounts(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRole"}:
		obj := &rbacv1.ClusterRole{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.RbacV1().ClusterRoles().Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRoleBinding"}:
		obj := &rbacv1.ClusterRoleBinding{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.nativeClientSet.RbacV1().ClusterRoleBindings().Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Affinity"}:
		obj := &openxv1.Affinity{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Affinities(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunAccessControl"}:
		obj := &openxv1.AliyunAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().AliyunAccessControls(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcAccessControl"}:
		obj := &openxv1.VolcAccessControl{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().VolcAccessControls(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunLoadBalancer"}:
		obj := &openxv1.AliyunLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().AliyunLoadBalancers(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcLoadBalancer"}:
		obj := &openxv1.VolcLoadBalancer{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().VolcLoadBalancers(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Etcd"}:
		obj := &openxv1.Etcd{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Etcds(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Mysql"}:
		obj := &openxv1.Mysql{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Mysqls(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "NodeSelector"}:
		obj := &openxv1.NodeSelector{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().NodeSelectors(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Redis"}:
		obj := &openxv1.Redis{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Redises(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Openx"}:
		obj := &openxv1.Openx{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		if err = r.prepareUpdateOpenx(obj); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Openxes(namespace).Update(ctx, obj, updateOpts)
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Toleration"}:
		obj := &openxv1.Toleration{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		_, err = r.openxClientSet.OpenxV1().Tolerations(namespace).Update(ctx, obj, updateOpts)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (r *Resources) list(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	selector, err := v1.LabelSelectorAsSelector(&v1.LabelSelector{
		MatchLabels:      nil,
		MatchExpressions: nil,
	})
	if err != nil {
		return 1, nil, err
	}
	var objList Object
	switch *gvk {
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Namespace"}:
		obj, err := r.InformerFactory.Core().V1().Namespaces().Lister().List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.NamespaceList{Items: make([]corev1.Namespace, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Node"}:
		obj, err := r.InformerFactory.Core().V1().Nodes().Lister().List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.NodeList{Items: make([]corev1.Node, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolume"}:
		obj, err := r.InformerFactory.Core().V1().PersistentVolumes().Lister().List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.PersistentVolumeList{Items: make([]corev1.PersistentVolume, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRole"}:
		obj, err := r.InformerFactory.Rbac().V1().ClusterRoles().Lister().List(selector)
		if err != nil {
			break
		}
		listItem := &rbacv1.ClusterRoleList{Items: make([]rbacv1.ClusterRole, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRoleBinding"}:
		obj, err := r.InformerFactory.Rbac().V1().ClusterRoleBindings().Lister().List(selector)
		if err != nil {
			break
		}
		listItem := &rbacv1.ClusterRoleBindingList{Items: make([]rbacv1.ClusterRoleBinding, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}:
		obj, err := r.InformerFactory.Apps().V1().DaemonSets().Lister().DaemonSets(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &appsv1.DaemonSetList{Items: make([]appsv1.DaemonSet, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}:
		obj, err := r.InformerFactory.Apps().V1().Deployments().Lister().Deployments(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &appsv1.DeploymentList{Items: make([]appsv1.Deployment, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}:
		obj, err := r.InformerFactory.Apps().V1().StatefulSets().Lister().StatefulSets(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &appsv1.StatefulSetList{Items: make([]appsv1.StatefulSet, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: autoscalingv2.GroupName, Version: autoscalingv2.SchemeGroupVersion.Version, Kind: "HorizontalPodAutoscaler"}:
		obj, err := r.InformerFactory.Autoscaling().V2().HorizontalPodAutoscalers().Lister().HorizontalPodAutoscalers(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &autoscalingv2.HorizontalPodAutoscalerList{Items: make([]autoscalingv2.HorizontalPodAutoscaler, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ConfigMap"}:
		obj, err := r.InformerFactory.Core().V1().ConfigMaps().Lister().ConfigMaps(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.ConfigMapList{Items: make([]corev1.ConfigMap, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Endpoints"}:
		obj, err := r.InformerFactory.Core().V1().Endpoints().Lister().Endpoints(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.EndpointsList{Items: make([]corev1.Endpoints, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj, err := r.InformerFactory.Core().V1().Events().Lister().Events(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.EventList{Items: make([]corev1.Event, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolumeClaim"}:
		obj, err := r.InformerFactory.Core().V1().PersistentVolumeClaims().Lister().PersistentVolumeClaims(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.PersistentVolumeClaimList{Items: make([]corev1.PersistentVolumeClaim, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}:
		obj, err := r.InformerFactory.Core().V1().Pods().Lister().Pods(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.PodList{Items: make([]corev1.Pod, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Secret"}:
		obj, err := r.InformerFactory.Core().V1().Secrets().Lister().Secrets(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.SecretList{Items: make([]corev1.Secret, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Service"}:
		obj, err := r.InformerFactory.Core().V1().Services().Lister().Services(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.ServiceList{Items: make([]corev1.Service, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ServiceAccount"}:
		obj, err := r.InformerFactory.Core().V1().ServiceAccounts().Lister().ServiceAccounts(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &corev1.ServiceAccountList{Items: make([]corev1.ServiceAccount, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Affinity"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Affinities().Lister().Affinities(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.AffinityList{Items: make([]openxv1.Affinity, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunAccessControl"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().AliyunAccessControls().Lister().AliyunAccessControls(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.AliyunAccessControlList{Items: make([]openxv1.AliyunAccessControl, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcAccessControl"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().VolcAccessControls().Lister().VolcAccessControls(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.VolcAccessControlList{Items: make([]openxv1.VolcAccessControl, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunLoadBalancer"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().AliyunLoadBalancers().Lister().AliyunLoadBalancers(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.AliyunLoadBalancerList{Items: make([]openxv1.AliyunLoadBalancer, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcLoadBalancer"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().VolcLoadBalancers().Lister().VolcLoadBalancers(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.VolcLoadBalancerList{Items: make([]openxv1.VolcLoadBalancer, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Etcd"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Etcds().Lister().Etcds(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.EtcdList{Items: make([]openxv1.Etcd, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Mysql"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Mysqls().Lister().Mysqls(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.MysqlList{Items: make([]openxv1.Mysql, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "NodeSelector"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().NodeSelectors().Lister().NodeSelectors(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.NodeSelectorList{Items: make([]openxv1.NodeSelector, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Redis"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Redises().Lister().Redises(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.RedisList{Items: make([]openxv1.Redis, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Openx"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Openxes().Lister().Openxes(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.OpenxList{Items: make([]openxv1.Openx, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	case schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Toleration"}:
		obj, err := r.OpenXInformerFactory.Openx().V1().Tolerations().Lister().Tolerations(namespace).List(selector)
		if err != nil {
			break
		}
		listItem := &openxv1.TolerationList{Items: make([]openxv1.Toleration, 0)}
		for _, v := range obj {
			listItem.Items = append(listItem.Items, *v)
		}
		objList = listItem
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		return 1, nil, err
	}
	res, err = objList.Marshal()
	if err != nil {
		return 1, nil, err
	}
	return code, res, err
}

func (r *Resources) convertGroupVersionKind(object interface{}) schema.GroupVersionKind {
	switch reflect.TypeOf(object) {
	case reflect.TypeOf(&appsv1.DaemonSet{}):
		return schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}
	case reflect.TypeOf(&appsv1.Deployment{}):
		return schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}
	case reflect.TypeOf(&appsv1.StatefulSet{}):
		return schema.GroupVersionKind{Group: appsv1.GroupName, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}
	case reflect.TypeOf(&autoscalingv2.HorizontalPodAutoscaler{}):
		return schema.GroupVersionKind{Group: autoscalingv2.GroupName, Version: autoscalingv2.SchemeGroupVersion.Version, Kind: "HorizontalPodAutoscaler"}
	case reflect.TypeOf(&corev1.ConfigMap{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ConfigMap"}
	case reflect.TypeOf(&corev1.Endpoints{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Endpoints"}
	case reflect.TypeOf(&corev1.Event{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Event"}
	case reflect.TypeOf(&corev1.Namespace{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Namespace"}
	case reflect.TypeOf(&corev1.Node{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Node"}
	case reflect.TypeOf(&corev1.PersistentVolume{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolume"}
	case reflect.TypeOf(&corev1.PersistentVolumeClaim{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "PersistentVolumeClaim"}
	case reflect.TypeOf(&corev1.Pod{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}
	case reflect.TypeOf(&corev1.Secret{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Secret"}
	case reflect.TypeOf(&corev1.Service{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "Service"}
	case reflect.TypeOf(&corev1.ServiceAccount{}):
		return schema.GroupVersionKind{Group: corev1.GroupName, Version: corev1.SchemeGroupVersion.Version, Kind: "ServiceAccount"}
	case reflect.TypeOf(&rbacv1.ClusterRole{}):
		return schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRole"}
	case reflect.TypeOf(&rbacv1.ClusterRoleBinding{}):
		return schema.GroupVersionKind{Group: rbacv1.GroupName, Version: rbacv1.SchemeGroupVersion.Version, Kind: "ClusterRoleBinding"}
	case reflect.TypeOf(&openxv1.Affinity{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Affinity"}
	case reflect.TypeOf(&openxv1.AliyunAccessControl{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunAccessControl"}
	case reflect.TypeOf(&openxv1.VolcAccessControl{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcAccessControl"}
	case reflect.TypeOf(&openxv1.AliyunLoadBalancer{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "AliyunLoadBalancer"}
	case reflect.TypeOf(&openxv1.VolcLoadBalancer{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "VolcLoadBalancer"}
	case reflect.TypeOf(&openxv1.Etcd{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Etcd"}
	case reflect.TypeOf(&openxv1.Mysql{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Mysql"}
	case reflect.TypeOf(&openxv1.NodeSelector{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "NodeSelector"}
	case reflect.TypeOf(&openxv1.Redis{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Redis"}
	case reflect.TypeOf(&openxv1.Openx{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Openx"}
	case reflect.TypeOf(&openxv1.Toleration{}):
		return schema.GroupVersionKind{Group: openxv1.GroupName, Version: openxv1.SchemeGroupVersion.Version, Kind: "Toleration"}
	case reflect.TypeOf(&metricsv1beta1.NodeMetrics{}):
		return schema.GroupVersionKind{Group: metricsv1beta1.GroupName, Version: metricsv1beta1.SchemeGroupVersion.Version, Kind: "NodeMetrics"}
	case reflect.TypeOf(&metricsv1beta1.PodMetrics{}):
		return schema.GroupVersionKind{Group: metricsv1beta1.GroupName, Version: metricsv1beta1.SchemeGroupVersion.Version, Kind: "PodMetrics"}
	default:
		zaplogger.Sugar().Errorf(ErrGVKNotExist, object)
		return schema.GroupVersionKind{}
	}
}
