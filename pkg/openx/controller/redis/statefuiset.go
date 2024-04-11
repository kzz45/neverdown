package redis

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func validateReplicas(role openxv1.ClusterRole, replicas *int32) *int32 {
	switch role {
	case openxv1.ClusterRoleMaster:
		if *replicas > 1 {
			var a int32 = 1
			return &a
		}
	case openxv1.ClusterRoleSlave:
		return replicas
	}
	return replicas
}

func (rc *RedisController) newStatefulSet(m *openxv1.Redis, in openxv1.RedisConfig) *appsv1.StatefulSet {
	var labels = map[string]string{
		controller.LabelKind:        controllerKind.Kind,
		controller.LabelController:  m.Name,
		controller.LabelClusterRole: string(in.Role),
	}
	objectName := controller.ClusterRoleStatefulSetName(in.Role, m.Name)
	var (
		volumes []corev1.Volume
	)
	if in.PersistentStorage.StorageVolumePath != "" {
		t := corev1.HostPathDirectoryOrCreate
		hostPath := &corev1.HostPathVolumeSource{
			Type: &t,
			Path: fmt.Sprintf("%s/%s/%s/%s", in.PersistentStorage.StorageVolumePath, m.Namespace, controllerName, objectName),
		}
		volumes = []corev1.Volume{
			{
				Name: controller.NasVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: hostPath,
				},
			},
		}
	}
	if in.Pod.Spec.NodeSelector == nil {
		in.Pod.Spec.NodeSelector = make(map[string]string)
	}
	if in.Pod.Spec.Affinity == nil {
		in.Pod.Spec.Affinity = &corev1.Affinity{}
	}
	if in.Pod.Spec.Tolerations == nil {
		in.Pod.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			//Annotations: rc.cloudConfig.Annotation(in.Service.Spec.Type, in.CloudNetworkConfig.ServiceWhiteList),
			Name:      objectName,
			Namespace: m.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m, controllerKind),
			},
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: validateReplicas(in.Role, in.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: in.Pod.Annotations,
					Labels:      labels,
				},
				Spec: corev1.PodSpec{
					Volumes:            volumes,
					Containers:         []corev1.Container{},
					NodeSelector:       in.Pod.Spec.NodeSelector,
					ServiceAccountName: in.Pod.Spec.ServiceAccountName,
					ImagePullSecrets:   in.Pod.Spec.ImagePullSecrets,
					Affinity:           in.Pod.Spec.Affinity,
					Tolerations:        in.Pod.Spec.Tolerations,
					SecurityContext:    in.Pod.Spec.SecurityContext,
				},
			},
		},
	}
	if len(in.Pod.Spec.Containers) == 0 || len(in.Pod.Spec.Containers) > 1 {
		return sts
	}
	sts.Spec.Template.Spec.Containers = make([]corev1.Container, 0)
	for _, container := range in.Pod.Spec.Containers {
		ports := []corev1.ContainerPort{
			{
				ContainerPort: RedisDefaultPort,
			},
		}
		if len(container.Ports) > 0 {
			ports = container.Ports
		}
		port := strconv.Itoa(int(ports[0].ContainerPort))
		envs := []corev1.EnvVar{
			{
				Name:  EnvRedisConf,
				Value: fmt.Sprintf(EnvRedisConfTemplate, objectName),
			},
			{
				Name:  EnvRedisDir,
				Value: "",
			},
			{
				Name:  EnvRedisDbFileName,
				Value: fmt.Sprintf(EnvRedisDbFileNameTemplate, objectName),
			},
			{
				Name:  EnvRedisPort,
				Value: port,
			},
		}
		if in.Role == openxv1.ClusterRoleSlave {
			masterPort := strconv.Itoa(RedisDefaultPort)
			if m.Spec.Master.Service != nil && len(m.Spec.Master.Service.Spec.Ports) > 0 {
				masterPort = strconv.Itoa(int(m.Spec.Master.Service.Spec.Ports[0].Port))
			}

			envs = append(envs, []corev1.EnvVar{
				{
					Name:  "GET_HOSTS_FROM",
					Value: "dns",
				},
				{
					Name:  EnvRedisMasterHost,
					Value: controller.ClusterRoleServiceName(openxv1.ClusterRoleMaster, m.Name),
				},
				{
					Name:  EnvRedisMasterPort,
					Value: masterPort,
				},
			}...)
		}
		c := corev1.Container{
			Name:      objectName,
			Image:     container.Image,
			Command:   container.Command,
			Args:      container.Args,
			Ports:     ports,
			Env:       envs,
			Resources: container.Resources,
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: controller.NasMountPath,
					Name:      controller.NasVolumeName,
				},
			},
			ImagePullPolicy: corev1.PullAlways,
			SecurityContext: container.SecurityContext,
		}
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, c)
	}
	return sts
}

func (rc *RedisController) syncStatefulSets(m *openxv1.Redis, stsList []*appsv1.StatefulSet) error {
	cur := make([]*appsv1.StatefulSet, 0)
	if m.Spec.Master.Pod != nil {
		m.Spec.Master.Role = openxv1.ClusterRoleMaster
		cur = append(cur, rc.newStatefulSet(m, m.Spec.Master))
	}
	if m.Spec.Slave.Pod != nil {
		m.Spec.Slave.Role = openxv1.ClusterRoleSlave
		cur = append(cur, rc.newStatefulSet(m, m.Spec.Slave))
	}
	var errlist []error
	add, update, del := diffStatefulSets(cur, stsList)
	zaplogger.Sugar().Debugf("syncStatefulSets cur:%d stslist:%d add:%d update:%d del:%d", len(cur), len(stsList), len(add), len(update), len(del))
	for _, v := range add {
		if _, err := rc.stsControl.Create(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range update {
		if _, err := rc.stsControl.Update(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range del {
		if err := rc.stsControl.Delete(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	return utilerrors.NewAggregate(errlist)
}

func compareStatefulSet(s1, s2 *appsv1.StatefulSet) (int, bool) {
	if !reflect.DeepEqual(s1.Name, s2.Name) {
		return 1, true
	}
	var replicas int32 = 0
	if s1.Spec.Replicas == nil {
		s1.Spec.Replicas = &replicas
	}
	if s2.Spec.Replicas == nil {
		s2.Spec.Replicas = &replicas
	}
	if !reflect.DeepEqual(*s1.Spec.Replicas, *s2.Spec.Replicas) {
		return 2, true
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.Volumes, s2.Spec.Template.Spec.Volumes) {
		return 3, true
	}
	if s1.Spec.Template.Spec.NodeSelector == nil {
		s1.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	if s2.Spec.Template.Spec.NodeSelector == nil {
		s2.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.NodeSelector, s2.Spec.Template.Spec.NodeSelector) {
		return 4, true
	}
	if s1.Spec.Template.Spec.Affinity == nil {
		s1.Spec.Template.Spec.Affinity = &corev1.Affinity{}
	}
	if s2.Spec.Template.Spec.Affinity == nil {
		s2.Spec.Template.Spec.Affinity = &corev1.Affinity{}
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.Affinity, s2.Spec.Template.Spec.Affinity) {
		return 5, true
	}
	if s1.Spec.Template.Spec.Tolerations == nil {
		s1.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	if s2.Spec.Template.Spec.Tolerations == nil {
		s2.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.Tolerations, s2.Spec.Template.Spec.Tolerations) {
		return 6, true
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.ServiceAccountName, s2.Spec.Template.Spec.ServiceAccountName) {
		return 7, true
	}
	if s1.Spec.Template.Spec.SecurityContext == nil {
		s1.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	if s2.Spec.Template.Spec.SecurityContext == nil {
		s2.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	if !apiequality.Semantic.DeepEqual(s1.Spec.Template.Spec.SecurityContext, s2.Spec.Template.Spec.SecurityContext) {
		return 8, true
	}
	// container
	if len(s1.Spec.Template.Spec.Containers) != len(s2.Spec.Template.Spec.Containers) {
		return 9, true
	}
	if len(s1.Spec.Template.Spec.Containers) == 0 {
		return 10, false
	}
	p1, p2 := s1.Spec.Template.Spec.Containers[0], s2.Spec.Template.Spec.Containers[0]
	if !reflect.DeepEqual(p1.Image, p2.Image) {
		return 11, true
	}
	if !reflect.DeepEqual(p1.Ports, p2.Ports) {
		return 12, true
	}
	if !reflect.DeepEqual(p1.Env, p2.Env) {
		return 13, true
	}
	if !reflect.DeepEqual(p1.Resources, p2.Resources) {
		return 14, true
	}
	if !reflect.DeepEqual(p1.VolumeMounts, p2.VolumeMounts) {
		return 15, true
	}
	if !apiequality.Semantic.DeepEqual(p1.SecurityContext, p2.SecurityContext) {
		return 16, true
	}
	if _, ok := controller.ComparePodTemplateSpecAnnotations(s1.Spec.Template, s2.Spec.Template); !ok {
		return 17, true
	}
	return 18, false
}

func diffStatefulSets(cur, ori []*appsv1.StatefulSet) ([]*appsv1.StatefulSet, []*appsv1.StatefulSet, []*appsv1.StatefulSet) {
	now := make(map[string]*appsv1.StatefulSet)
	old := make(map[string]*appsv1.StatefulSet)
	for _, v := range cur {
		now[v.Name] = v.DeepCopy()
	}
	for _, v := range ori {
		old[v.Name] = v.DeepCopy()
	}
	add := make([]*appsv1.StatefulSet, 0)
	up := make([]*appsv1.StatefulSet, 0)
	del := make([]*appsv1.StatefulSet, 0)
	// add & update
	for k, v := range now {
		if t, ok := old[k]; !ok {
			add = append(add, v.DeepCopy())
			zaplogger.Sugar().Infow("diffStatefulSets add",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		} else {
			// compare change
			if code, changed := compareStatefulSet(v, t); changed {
				up = append(up, v.DeepCopy())
				zaplogger.Sugar().Infow("diffStatefulSets update",
					zap.String("controller", controllerName),
					zap.String("name", k), zap.String("namespace", v.Namespace), zap.Int("code", code))
			}
		}
	}
	// del
	for k, v := range old {
		if _, ok := now[k]; !ok {
			del = append(del, v.DeepCopy())
			zaplogger.Sugar().Infow("diffStatefulSets delete",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		}
	}
	return add, up, del
}
