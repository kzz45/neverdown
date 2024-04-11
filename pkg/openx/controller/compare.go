package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func CompareDeploymentAnnotations(newDeploy, oriDeploy *appsv1.Deployment) (int, bool) {
	// original Deployment maybe contain the annotations below
	metrics := map[string]bool{
		"deployment.kubernetes.io/revision": true,
	}
	for k1, v1 := range newDeploy.Annotations {
		if v2, ok := oriDeploy.Annotations[k1]; !ok {
			return 1, false
		} else {
			if v2 != v1 {
				return 2, false
			}
		}
	}
	oriLength := 0
	for k, _ := range oriDeploy.Annotations {
		if _, ok := metrics[k]; !ok {
			oriLength++
		}
	}
	if oriLength != len(newDeploy.Annotations) {
		return 3, false
	}
	return 0, true
}

func CompareStatefulSetAnnotations(newSts, oriSts *appsv1.StatefulSet) (int, bool) {
	for k1, v1 := range newSts.Annotations {
		if v2, ok := oriSts.Annotations[k1]; !ok {
			return 1, false
		} else {
			if v2 != v1 {
				return 2, false
			}
		}
	}
	if len(oriSts.Annotations) != len(newSts.Annotations) {
		return 3, false
	}
	return 0, true
}

func ComparePodTemplateSpecAnnotations(newSts, oriSts corev1.PodTemplateSpec) (int, bool) {
	for k1, v1 := range newSts.Annotations {
		if v2, ok := oriSts.Annotations[k1]; !ok {
			return 1, false
		} else {
			if v2 != v1 {
				return 2, false
			}
		}
	}
	if len(oriSts.Annotations) != len(newSts.Annotations) {
		return 3, false
	}
	return 0, true
}
