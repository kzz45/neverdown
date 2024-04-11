package controller

import (
	"fmt"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
)

const (
	LabelKind       = "kind.openx.neverdown.io"
	LabelController = "controller.openx.neverdown.io"
	LabelAppName    = "app.openx.openx.neverdown.io"

	// LabelClusterRole was only used for mysql or redis
	LabelClusterRole = "clusterRole.openx.neverdown.io"

	LabelName = "name"
)

func ClusterRoleServiceName(role openxv1.ClusterRole, name string) string {
	return fmt.Sprintf("%s-%s", name, role)
}

func ClusterRoleStatefulSetName(role openxv1.ClusterRole, name string) string {
	return fmt.Sprintf("%s-%s", name, role)
}

func OpenxAppName(objectName string, appName string) string {
	return fmt.Sprintf("%s-%s", objectName, appName)
}

func OpenxServiceName(objectName string, appName string) string {
	return fmt.Sprintf("%s-%s", objectName, appName)
}

func OpenxExtensionServiceName(objectName string, appName string) string {
	return fmt.Sprintf("%s-%s-ext", objectName, appName)
}

func OpenxHpaName(objectName string, appName string) string {
	return fmt.Sprintf("%s-%s", objectName, appName)
}

func EtcdServiceName(objectName string) string {
	return fmt.Sprintf("%s-headless", objectName)
}

const (
	LabelJingxProject    = "jingx-project.openx.neverdown.io"
	LabelJingxRepository = "jingx-repository.openx.neverdown.io"
	LabelJingxTag        = "jingx-tag.openx.neverdown.io"
	LabelWatchPolicy     = "watch-policy.openx.neverdown.io"
)
