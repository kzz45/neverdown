package app

import "k8s.io/apimachinery/pkg/runtime/schema"

func (c *ControllerContext) runInformer(stop <-chan struct{}) {
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "DaemonSets",
	}, c.InformerFactory.Apps().V1().DaemonSets().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployments",
	}, c.InformerFactory.Apps().V1().Deployments().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "StatefulSets",
	}, c.InformerFactory.Apps().V1().StatefulSets().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "autoscaling",
		Version: "v2",
		Kind:    "HorizontalPodAutoscalers",
	}, c.InformerFactory.Autoscaling().V2().HorizontalPodAutoscalers().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "ConfigMaps",
	}, c.InformerFactory.Core().V1().ConfigMaps().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Endpoints",
	}, c.InformerFactory.Core().V1().Endpoints().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Events",
	}, c.InformerFactory.Core().V1().Events().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Namespaces",
	}, c.InformerFactory.Core().V1().Namespaces().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Nodes",
	}, c.InformerFactory.Core().V1().Nodes().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "PersistentVolumes",
	}, c.InformerFactory.Core().V1().PersistentVolumes().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "PersistentVolumeClaims",
	}, c.InformerFactory.Core().V1().PersistentVolumeClaims().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Pods",
	}, c.InformerFactory.Core().V1().Pods().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Secrets",
	}, c.InformerFactory.Core().V1().Secrets().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Services",
	}, c.InformerFactory.Core().V1().Services().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "ServiceAccounts",
	}, c.InformerFactory.Core().V1().ServiceAccounts().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "rbac",
		Version: "v1",
		Kind:    "ClusterRoles",
	}, c.InformerFactory.Rbac().V1().ClusterRoles().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "rbac",
		Version: "v1",
		Kind:    "ClusterRoleBindings",
	}, c.InformerFactory.Rbac().V1().ClusterRoleBindings().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Affinities",
	}, c.OpenXInformerFactory.Openx().V1().Affinities().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "AliyunAccessControls",
	}, c.OpenXInformerFactory.Openx().V1().AliyunAccessControls().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "VolcAccessControls",
	}, c.OpenXInformerFactory.Openx().V1().VolcAccessControls().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "AliyunLoadBalancers",
	}, c.OpenXInformerFactory.Openx().V1().AliyunLoadBalancers().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "VolcLoadBalancers",
	}, c.OpenXInformerFactory.Openx().V1().VolcLoadBalancers().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Etcds",
	}, c.OpenXInformerFactory.Openx().V1().Etcds().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Mysqls",
	}, c.OpenXInformerFactory.Openx().V1().Mysqls().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "NodeSelectors",
	}, c.OpenXInformerFactory.Openx().V1().NodeSelectors().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Redises",
	}, c.OpenXInformerFactory.Openx().V1().Redises().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Openxes",
	}, c.OpenXInformerFactory.Openx().V1().Openxes().Informer(), stop)
	go c.runSharedIndexInformer(schema.GroupVersionKind{
		Group:   "openx",
		Version: "v1",
		Kind:    "Tolerations",
	}, c.OpenXInformerFactory.Openx().V1().Tolerations().Informer(), stop)
}
