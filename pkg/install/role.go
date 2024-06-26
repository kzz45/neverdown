package install

import (
	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	appsv1 "k8s.io/api/apps/v1"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *installer) installClusterRoles() {
	b.installRoles()
	b.installClusterRoleBinding()
	b.installServiceAccount()
}

func (b *installer) installRoles() {
	verbs := []string{
		"create",
		"get",
		"list",
		"watch",
		"delete",
		"update",
		"patch",
	}
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openx-controller",
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: verbs,
				APIGroups: []string{
					corev1.GroupName,
				},
				Resources: []string{
					"events",
					"services",
					"endpoints",
					"secrets",
					"namespaces",
					"configmaps",
					"pods",
					"pods/log",
					"pods/exec",
					"nodes",
					"persistentvolumes",
					"persistentvolumeclaims",
					"serviceaccounts",
					"nodes/proxy",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
			{
				Verbs: verbs,
				APIGroups: []string{
					appsv1.GroupName,
				},
				Resources: []string{
					"statefulsets",
					"deployments",
					"daemonsets",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
			{
				Verbs: verbs,
				APIGroups: []string{
					rbacv1.GroupName,
				},
				Resources: []string{
					"clusterroles",
					"clusterrolebindings",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
			{
				Verbs: verbs,
				APIGroups: []string{
					coordinationv1.GroupName,
				},
				Resources: []string{
					"leases",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
			{
				Verbs: verbs,
				APIGroups: []string{
					extensionsv1beta1.GroupName,
				},
				Resources: []string{
					"ingresses",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
			{
				Verbs: verbs,
				APIGroups: []string{
					openxv1.GroupName,
				},
				Resources: []string{
					"loadbalancers",
					"accesscontrols",
					"mysqls",
					"mysqls/status",
					"redises",
					"redises/status",
					"openxs",
					"openxs/status",
					"affinities",
					"tolerations",
					"nodeselectors",
					"etcds",
					"etcds/status",
				},
				ResourceNames:   nil,
				NonResourceURLs: nil,
			},
		},
	}
	if _, err := b.clientSet.RbacV1().ClusterRoles().Create(b.ctx, cr, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
		}
		zaplogger.Sugar().Infof("Already exist ClusterRoles:%s", cr.Name)
	} else {
		zaplogger.Sugar().Infof("Successful create ClusterRoles:%s", cr.Name)
	}
}

func (b *installer) installClusterRoleBinding() {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openx-controller",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      "openx-controller",
				Namespace: "kube-api",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "openx-controller",
		},
	}
	if _, err := b.clientSet.RbacV1().ClusterRoleBindings().Create(b.ctx, crb, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
		}
		zaplogger.Sugar().Infof("Already exist ClusterRoleBindings:%s", crb.Name)
	} else {
		zaplogger.Sugar().Infof("Successful create ClusterRoleBindings:%s", crb.Name)
	}
}

func (b *installer) installServiceAccount() {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openx-controller",
			Namespace: NamespaceKubeApi,
		},
	}
	if _, err := b.clientSet.CoreV1().ServiceAccounts(sa.Namespace).Create(b.ctx, sa, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
		}
		zaplogger.Sugar().Infof("Already exist ServiceAccounts:%s namespace:%s", sa.Name, sa.Namespace)
	} else {
		zaplogger.Sugar().Infof("Successful create ServiceAccounts:%s namespace:%s", sa.Name, sa.Namespace)
	}
}
