package etcd

import (
	"fmt"
	"reflect"
	"strings"

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

func (ec *EtcdController) newStatefulSet(e *openxv1.Etcd) *appsv1.StatefulSet {
	var labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: e.Name,
	}
	var (
		volumes []corev1.Volume
	)
	if e.Spec.PersistentStorage.StorageVolumePath != "" {
		t := corev1.HostPathDirectoryOrCreate
		hostPath := &corev1.HostPathVolumeSource{
			Type: &t,
			Path: fmt.Sprintf("%s/%s/%s/%s", e.Spec.PersistentStorage.StorageVolumePath, e.Namespace, controllerName, e.Name),
		}
		volumes = []corev1.Volume{
			{
				Name: controller.NasVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: hostPath,
					//PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					//	ClaimName: fmt.Sprintf(PVCNameTemplate, foo.Spec.MasterSpec.DeploymentName, MasterName),
					//},
				},
			},
		}
	}
	if e.Spec.Pod.Spec.NodeSelector == nil {
		e.Spec.Pod.Spec.NodeSelector = make(map[string]string)
	}
	if e.Spec.Pod.Spec.Affinity == nil {
		e.Spec.Pod.Spec.Affinity = &corev1.Affinity{}
	}
	if e.Spec.Pod.Spec.Tolerations == nil {
		e.Spec.Pod.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	var root int64 = 0
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			//Annotations: ec.cloudConfig.Annotation(in.Service.Spec.Type, in.CloudNetworkConfig.ServiceWhiteList),
			Name:      e.Name,
			Namespace: e.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(e, controllerKind),
			},
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: e.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: e.Spec.Pod.Annotations,
					Labels:      labels,
				},
				Spec: corev1.PodSpec{
					Volumes:            volumes,
					Containers:         []corev1.Container{},
					NodeSelector:       e.Spec.Pod.Spec.NodeSelector,
					ServiceAccountName: e.Spec.Pod.Spec.ServiceAccountName,
					ImagePullSecrets:   e.Spec.Pod.Spec.ImagePullSecrets,
					Affinity:           e.Spec.Pod.Spec.Affinity,
					Tolerations:        e.Spec.Pod.Spec.Tolerations,
					SecurityContext: &corev1.PodSecurityContext{
						SELinuxOptions:      nil,
						WindowsOptions:      nil,
						RunAsUser:           &root,
						RunAsGroup:          nil,
						RunAsNonRoot:        nil,
						SupplementalGroups:  nil,
						FSGroup:             &root,
						Sysctls:             nil,
						FSGroupChangePolicy: nil,
						SeccompProfile:      nil,
					},
				},
			},
			ServiceName: controller.EtcdServiceName(e.Name),
		},
	}
	if len(e.Spec.Pod.Spec.Containers) == 0 || len(e.Spec.Pod.Spec.Containers) > 1 {
		return sts
	}

	peers := make([]string, 0)
	for i := 0; i < int(*e.Spec.Replicas); i++ {
		peers = append(peers, fmt.Sprintf("%s-%d=http://%s-%d.%s:%d", e.Name, i, e.Name, i, controller.EtcdServiceName(e.Name), PeersPort))
	}
	advertiseClientUrl := fmt.Sprintf("http://${HOSTNAME}.%s:%d", controller.EtcdServiceName(e.Name), ClientPort)
	initialAdvertisePeerUrls := fmt.Sprintf("http://${HOSTNAME}.%s:%d", controller.EtcdServiceName(e.Name), PeersPort)
	commands := `
PEERS="%s"
exec etcd --name ${HOSTNAME} \
	--listen-peer-urls http://0.0.0.0:2380 \
	--listen-client-urls http://0.0.0.0:2379 \
	--advertise-client-urls %s \
	--initial-advertise-peer-urls %s \
	--initial-cluster-token etcd-cluster-1 \
	--initial-cluster ${PEERS} \
	--initial-cluster-state new \
	--data-dir /var/run/etcd/${HOSTNAME}.etcd \
	--logger=zap
`
	sts.Spec.Template.Spec.Containers = make([]corev1.Container, 0)
	for _, container := range e.Spec.Pod.Spec.Containers {
		ports := []corev1.ContainerPort{
			{
				Name:          "client",
				ContainerPort: ClientPort,
			},
			{
				Name:          "peer",
				ContainerPort: PeersPort,
			},
		}
		c := corev1.Container{
			Name:  e.Name,
			Image: container.Image,
			Command: []string{
				"/bin/sh",
				"-c",
				fmt.Sprintf(commands, strings.Join(peers, ","), advertiseClientUrl, initialAdvertisePeerUrls),
			},
			Args:  container.Args,
			Ports: ports,
			Env: []corev1.EnvVar{
				{
					Name:  "GET_HOSTS_FROM",
					Value: "dns",
				},
			},
			Resources: container.Resources,
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: "/var/run/etcd",
					Name:      controller.NasVolumeName,
				},
			},
			ImagePullPolicy: corev1.PullAlways,
		}
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, c)
	}
	return sts
}

func (ec *EtcdController) syncStatefulSets(e *openxv1.Etcd, stsList []*appsv1.StatefulSet) error {
	cur := make([]*appsv1.StatefulSet, 0)
	cur = append(cur, ec.newStatefulSet(e))
	var errlist []error
	add, update, del := diffStatefulSets(cur, stsList)
	zaplogger.Sugar().Debugf("syncStatefulSets cur:%d stslist:%d add:%d update:%d del:%d", len(cur), len(stsList), len(add), len(update), len(del))
	for _, v := range add {
		if _, err := ec.stsControl.Create(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range update {
		if _, err := ec.stsControl.Update(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range del {
		if err := ec.stsControl.Delete(v); err != nil {
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
