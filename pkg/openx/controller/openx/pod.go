package openx

import corev1 "k8s.io/api/core/v1"

// ExposePodInformationByEnvs appends envs to expose more information about the current pod
func ExposePodInformationByEnvs(envs []corev1.EnvVar) []corev1.EnvVar {
	//envs = append(envs, corev1.EnvVar{
	//	Name: NodeName,
	//	ValueFrom: &corev1.EnvVarSource{
	//		FieldRef: &corev1.ObjectFieldSelector{
	//			FieldPath: "spec.nodeName",
	//		},
	//	},
	//})
	//envs = append(envs, corev1.EnvVar{
	//	Name: HostIP,
	//	ValueFrom: &corev1.EnvVarSource{
	//		FieldRef: &corev1.ObjectFieldSelector{
	//			FieldPath: "status.hostIP",
	//		},
	//	},
	//})
	//envs = append(envs, corev1.EnvVar{
	//	Name: PodName,
	//	ValueFrom: &corev1.EnvVarSource{
	//		FieldRef: &corev1.ObjectFieldSelector{
	//			FieldPath: "metadata.name",
	//		},
	//	},
	//})
	envs = append(envs, corev1.EnvVar{
		Name: PodNamespace,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})
	envs = append(envs, corev1.EnvVar{
		Name: PodIP,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	})
	//envs = append(envs, corev1.EnvVar{
	//	Name: PodServiceAccount,
	//	ValueFrom: &corev1.EnvVarSource{
	//		FieldRef: &corev1.ObjectFieldSelector{
	//			FieldPath: "spec.serviceAccountName",
	//		},
	//	},
	//})
	return envs
}

func FilterExposePodInformation(envs []corev1.EnvVar) []corev1.EnvVar {
	res := make([]corev1.EnvVar, 0)
	for _, env := range envs {
		if env.ValueFrom == nil {
			res = append(res, env)
		}
	}
	return res
}
