package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"time"
)

var (
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type ResyncPeriodFunc func() time.Duration

// Returns 0 for resyncPeriod in case resyncing is not needed.
func NoResyncPeriodFunc() time.Duration {
	return 0
}

// StaticResyncPeriodFunc returns the resync period specified
func StaticResyncPeriodFunc(resyncPeriod time.Duration) ResyncPeriodFunc {
	return func() time.Duration {
		return resyncPeriod
	}
}

// ServiceControlInterface is an interface that knows how to add or delete
// Services, as well as increment or decrement them. It is used
// by the Service controller to ease testing of actions that it takes.
type ServiceControlInterface interface {
	Create(obj *corev1.Service) (*corev1.Service, error)
	Update(obj *corev1.Service) (*corev1.Service, error)
	Delete(obj *corev1.Service) error
	Patch(namespace, name string, data []byte) error
}

// RealServiceControl is the default implementation of ServiceControllerInterface.
type RealServiceControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

var _ ServiceControlInterface = &RealServiceControl{}

func (r RealServiceControl) Create(obj *corev1.Service) (*corev1.Service, error) {
	return r.KubeClient.CoreV1().Services(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (r RealServiceControl) Update(obj *corev1.Service) (*corev1.Service, error) {
	return r.KubeClient.CoreV1().Services(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (r RealServiceControl) Delete(obj *corev1.Service) error {
	return r.KubeClient.CoreV1().Services(obj.Namespace).Delete(context.TODO(), obj.Name, metav1.DeleteOptions{})
}

func (r RealServiceControl) Patch(namespace, name string, data []byte) error {
	_, err := r.KubeClient.CoreV1().Services(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	return err
}

// DeploymentControlInterface is an interface that knows how to add or delete
// Deployments, as well as increment or decrement them. It is used
// by the Deployment controller to ease testing of actions that it takes.
type DeploymentControlInterface interface {
	Create(obj *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(obj *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(obj *appsv1.Deployment) error
	Patch(namespace, name string, data []byte) error
}

// RealDeploymentControl is the default implementation of DeploymentControllerInterface.
type RealDeploymentControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

var _ DeploymentControlInterface = &RealDeploymentControl{}

func (r RealDeploymentControl) Create(obj *appsv1.Deployment) (*appsv1.Deployment, error) {
	return r.KubeClient.AppsV1().Deployments(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (r RealDeploymentControl) Update(obj *appsv1.Deployment) (*appsv1.Deployment, error) {
	return r.KubeClient.AppsV1().Deployments(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (r RealDeploymentControl) Delete(obj *appsv1.Deployment) error {
	return r.KubeClient.AppsV1().Deployments(obj.Namespace).Delete(context.TODO(), obj.Name, metav1.DeleteOptions{})
}

func (r RealDeploymentControl) Patch(namespace, name string, data []byte) error {
	_, err := r.KubeClient.AppsV1().Deployments(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	return err
}

// StatefulSetControlInterface is an interface that knows how to add or delete
// StatefulSets, as well as increment or decrement them. It is used
// by the StatefulSet controller to ease testing of actions that it takes.
type StatefulSetControlInterface interface {
	Create(obj *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(obj *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(obj *appsv1.StatefulSet) error
	Patch(namespace, name string, data []byte) error
}

// RealStatefulSetControl is the default implementation of StatefulSetControllerInterface.
type RealStatefulSetControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

var _ StatefulSetControlInterface = &RealStatefulSetControl{}

func (r RealStatefulSetControl) Create(obj *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return r.KubeClient.AppsV1().StatefulSets(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (r RealStatefulSetControl) Update(obj *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return r.KubeClient.AppsV1().StatefulSets(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (r RealStatefulSetControl) Delete(obj *appsv1.StatefulSet) error {
	return r.KubeClient.AppsV1().StatefulSets(obj.Namespace).Delete(context.TODO(), obj.Name, metav1.DeleteOptions{})
}

func (r RealStatefulSetControl) Patch(namespace, name string, data []byte) error {
	_, err := r.KubeClient.AppsV1().StatefulSets(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	return err
}

type HorizontalPodAutoscalerControlInterface interface {
	Create(obj *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Update(obj *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error)
	Delete(obj *autoscalingv2.HorizontalPodAutoscaler) error
	Patch(namespace, name string, data []byte) error
}

type RealHorizontalPodAutoscalerControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

var _ HorizontalPodAutoscalerControlInterface = &RealHorizontalPodAutoscalerControl{}

func (r RealHorizontalPodAutoscalerControl) Create(obj *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return r.KubeClient.AutoscalingV2().HorizontalPodAutoscalers(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (r RealHorizontalPodAutoscalerControl) Update(obj *autoscalingv2.HorizontalPodAutoscaler) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return r.KubeClient.AutoscalingV2().HorizontalPodAutoscalers(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (r RealHorizontalPodAutoscalerControl) Delete(obj *autoscalingv2.HorizontalPodAutoscaler) error {
	return r.KubeClient.AutoscalingV2().HorizontalPodAutoscalers(obj.Namespace).Delete(context.TODO(), obj.Name, metav1.DeleteOptions{})
}

func (r RealHorizontalPodAutoscalerControl) Patch(namespace, name string, data []byte) error {
	_, err := r.KubeClient.AutoscalingV2().HorizontalPodAutoscalers(namespace).Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	return err
}
