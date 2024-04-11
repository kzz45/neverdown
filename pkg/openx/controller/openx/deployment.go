package openx

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

func convertImageToLabels(labels map[string]string, image string) map[string]string {
	t1 := strings.Split(image, "/")
	if len(t1) != 3 {
		return labels
	}
	t2 := strings.Split(t1[2], ":")
	if len(t2) != 2 {
		return labels
	}
	_, project, repository, tag := t1[0], t1[1], t2[0], t2[1]
	labels[controller.LabelJingxProject] = project
	labels[controller.LabelJingxRepository] = repository
	labels[controller.LabelJingxTag] = tag
	return labels
}

func appendVolumes(v []corev1.Volume) []corev1.Volume {
	var mode int32 = 420
	res := make([]corev1.Volume, 0)
	for _, v := range v {
		if v.ConfigMap != nil {
			if v.ConfigMap.DefaultMode == nil {
				v.ConfigMap.DefaultMode = &mode
			}
		}
		res = append(res, v)
	}
	return res
}

func (tc *OpenxController) newDeployment(m *openxv1.Openx, app openxv1.App, ohr *realHpaReplicas) *appsv1.Deployment {
	if len(app.Pod.Spec.Containers) == 0 || len(app.Pod.Spec.Containers) > 1 {
		return nil
	}
	var labels = map[string]string{
		controller.LabelKind:        controllerKind.Kind,
		controller.LabelController:  m.Name,
		controller.LabelAppName:     app.AppName,
		controller.LabelWatchPolicy: string(app.WatchPolicy),
	}
	labels = convertImageToLabels(labels, app.Pod.Spec.Containers[0].Image)
	for k, v := range app.Pod.Labels {
		labels[k] = v
	}
	objectName := controller.OpenxAppName(m.Name, app.AppName)
	zaplogger.Sugar().Debugf("newDeployment app.Name:%s sc:%#v", objectName, app.Pod.Spec.SecurityContext)
	if app.Pod.Spec.Volumes == nil {
		app.Pod.Spec.Volumes = make([]corev1.Volume, 0)
	}
	app.Pod.Spec.Volumes = appendVolumes(app.Pod.Spec.Volumes)
	if app.PersistentStorage.StorageVolumePath != "" {
		t := corev1.HostPathDirectoryOrCreate
		hostPath := &corev1.HostPathVolumeSource{
			Type: &t,
			Path: fmt.Sprintf("%s/%s/%s/%s", app.PersistentStorage.StorageVolumePath, m.Namespace, controllerName, objectName),
		}
		app.Pod.Spec.Volumes = append(app.Pod.Spec.Volumes,
			corev1.Volume{
				Name: controller.NasVolumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: hostPath,
				},
			},
		)
	}
	if app.Pod.Spec.NodeSelector == nil {
		app.Pod.Spec.NodeSelector = make(map[string]string)
	}
	if app.Pod.Spec.Affinity == nil {
		app.Pod.Spec.Affinity = &corev1.Affinity{}
	}
	if app.Pod.Spec.Tolerations == nil {
		app.Pod.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	replicas, ok := ohr.realHpaDesiredReplicas(objectName)
	if !ok || *replicas == 0 {
		replicas = app.Replicas
	}
	dp := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      objectName,
			Namespace: m.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m, controllerKind),
			},
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: app.Pod.Annotations,
					Labels:      labels,
				},
				Spec: corev1.PodSpec{
					Volumes:            app.Pod.Spec.Volumes,
					Containers:         []corev1.Container{},
					NodeSelector:       app.Pod.Spec.NodeSelector,
					ServiceAccountName: app.Pod.Spec.ServiceAccountName,
					ImagePullSecrets:   app.Pod.Spec.ImagePullSecrets,
					Affinity:           app.Pod.Spec.Affinity,
					Tolerations:        app.Pod.Spec.Tolerations,
					SecurityContext:    app.Pod.Spec.SecurityContext,
				},
			},
		},
	}
	dp.Spec.Template.Spec.Containers = make([]corev1.Container, 0)
	for _, container := range app.Pod.Spec.Containers {
		if container.VolumeMounts == nil {
			container.VolumeMounts = make([]corev1.VolumeMount, 0)
		}
		if app.PersistentStorage.StorageVolumePath != "" {
			container.VolumeMounts = append(container.VolumeMounts,
				corev1.VolumeMount{
					MountPath: controller.NasMountPath,
					Name:      controller.NasVolumeName,
				})
		}
		c := corev1.Container{
			Name:            objectName,
			Image:           container.Image,
			Command:         container.Command,
			Args:            container.Args,
			Ports:           container.Ports,
			Env:             ExposePodInformationByEnvs(container.Env),
			Resources:       container.Resources,
			VolumeMounts:    container.VolumeMounts,
			ImagePullPolicy: corev1.PullAlways,
			SecurityContext: container.SecurityContext,
		}
		dp.Spec.Template.Spec.Containers = append(dp.Spec.Template.Spec.Containers, c)
	}
	return dp
}

func (tc *OpenxController) syncDeployments(m *openxv1.Openx, dpList []*appsv1.Deployment, ohr *realHpaReplicas) error {
	cur := make([]*appsv1.Deployment, 0)
	for _, app := range m.Spec.Applications {
		if dp := tc.newDeployment(m, app, ohr); dp != nil {
			cur = append(cur, dp)
			zaplogger.Sugar().Debugf("OpenxController dp.Name:%s sc:%#v", dp.Name, dp.Spec.Template.Spec.SecurityContext)
		}
	}
	var errlist []error
	add, update, del := diffDeployments(cur, dpList)
	zaplogger.Sugar().Debugf("OpenxController syncDeployments cur:%d dplist:%d add:%d update:%d del:%d", len(cur), len(dpList), len(add), len(update), len(del))
	for _, v := range del {
		if err := tc.dpControl.Delete(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range add {
		if _, err := tc.dpControl.Create(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range update {
		if _, err := tc.dpControl.Update(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	return utilerrors.NewAggregate(errlist)
}

func compareDeployment(s1, s2 *appsv1.Deployment) (int, bool) {
	if !reflect.DeepEqual(s1.Name, s2.Name) {
		return 1, true
	}
	if !reflect.DeepEqual(s1.Labels, s2.Labels) {
		return 2, true
	}
	var replicas int32 = 0
	if s1.Spec.Replicas == nil {
		s1.Spec.Replicas = &replicas
	}
	if s2.Spec.Replicas == nil {
		s2.Spec.Replicas = &replicas
	}
	if !reflect.DeepEqual(*s1.Spec.Replicas, *s2.Spec.Replicas) {
		return 3, true
	}
	if s1.Spec.Template.Spec.Volumes == nil {
		s1.Spec.Template.Spec.Volumes = []corev1.Volume{}
	}
	if s2.Spec.Template.Spec.Volumes == nil {
		s2.Spec.Template.Spec.Volumes = []corev1.Volume{}
	}
	if len(s1.Spec.Template.Spec.Volumes) != len(s2.Spec.Template.Spec.Volumes) {
		return 4, true
	}
	for _, volume := range s1.Spec.Template.Spec.Volumes {
		exist := false
		for _, volumeReality := range s2.Spec.Template.Spec.Volumes {
			if volume.Name == volumeReality.Name {
				if !apiequality.Semantic.DeepEqual(volume, volumeReality) {
					return 5, true
				} else {
					exist = true
				}
			}
		}
		if !exist {
			return 6, true
		}
	}
	if s1.Spec.Template.Spec.NodeSelector == nil {
		s1.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	if s2.Spec.Template.Spec.NodeSelector == nil {
		s2.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.NodeSelector, s2.Spec.Template.Spec.NodeSelector) {
		return 7, true
	}
	if s1.Spec.Template.Spec.Affinity == nil {
		s1.Spec.Template.Spec.Affinity = &corev1.Affinity{}
	}
	if s2.Spec.Template.Spec.Affinity == nil {
		s2.Spec.Template.Spec.Affinity = &corev1.Affinity{}
	}
	if !apiequality.Semantic.DeepEqual(s1.Spec.Template.Spec.Affinity, s2.Spec.Template.Spec.Affinity) {
		return 8, true
	}
	if s1.Spec.Template.Spec.Tolerations == nil {
		s1.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	if s2.Spec.Template.Spec.Tolerations == nil {
		s2.Spec.Template.Spec.Tolerations = make([]corev1.Toleration, 0)
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.Tolerations, s2.Spec.Template.Spec.Tolerations) {
		return 9, true
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.ServiceAccountName, s2.Spec.Template.Spec.ServiceAccountName) {
		return 10, true
	}
	if s1.Spec.Template.Spec.SecurityContext == nil {
		s1.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	if s2.Spec.Template.Spec.SecurityContext == nil {
		s2.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
	if !reflect.DeepEqual(s1.Spec.Template.Spec.SecurityContext, s2.Spec.Template.Spec.SecurityContext) {
		zaplogger.Sugar().Infof("s1 name:%s sc:%#v", s1.Name, s1.Spec.Template.Spec.SecurityContext)
		zaplogger.Sugar().Infof("s2 name:%s sc:%#v", s2.Name, s2.Spec.Template.Spec.SecurityContext)
		return 11, true
	}
	// container
	if len(s1.Spec.Template.Spec.Containers) != len(s2.Spec.Template.Spec.Containers) {
		return 12, true
	}
	if len(s1.Spec.Template.Spec.Containers) == 0 {
		return 13, false
	}
	p1, p2 := s1.Spec.Template.Spec.Containers[0], s2.Spec.Template.Spec.Containers[0]
	if !reflect.DeepEqual(p1.Image, p2.Image) {
		return 14, true
	}
	if !apiequality.Semantic.DeepEqual(p1.Ports, p2.Ports) {
		return 15, true
	}
	if !reflect.DeepEqual(FilterExposePodInformation(p1.Env), FilterExposePodInformation(p2.Env)) {
		return 16, true
	}
	if !apiequality.Semantic.DeepEqual(p1.Resources, p2.Resources) {
		return 17, true
	}
	if p1.VolumeMounts == nil {
		p1.VolumeMounts = []corev1.VolumeMount{}
	}
	if p2.VolumeMounts == nil {
		p2.VolumeMounts = []corev1.VolumeMount{}
	}
	if len(p1.VolumeMounts) != len(p2.VolumeMounts) {
		return 18, true
	}
	for _, volume := range p1.VolumeMounts {
		exist := false
		for _, volumeReality := range p2.VolumeMounts {
			if volume.Name == volumeReality.Name {
				if !reflect.DeepEqual(volume, volumeReality) {
					return 19, true
				} else {
					exist = true
				}
			}
		}
		if !exist {
			return 20, true
		}
	}
	if !reflect.DeepEqual(p1.Command, p2.Command) {
		return 21, true
	}
	if !reflect.DeepEqual(p1.Args, p2.Args) {
		return 22, true
	}
	if !apiequality.Semantic.DeepEqual(p1.SecurityContext, p2.SecurityContext) {
		zaplogger.Sugar().Infof("container s1 name:%s sc:%#v", s1.Name, p1.SecurityContext)
		zaplogger.Sugar().Infof("container s2 name:%s sc:%#v", s2.Name, p2.SecurityContext)
		return 23, true
	}
	if _, ok := controller.ComparePodTemplateSpecAnnotations(s1.Spec.Template, s2.Spec.Template); !ok {
		return 24, true
	}
	return 25, false
}

func diffDeployments(cur, ori []*appsv1.Deployment) ([]*appsv1.Deployment, []*appsv1.Deployment, []*appsv1.Deployment) {
	now := make(map[string]*appsv1.Deployment)
	old := make(map[string]*appsv1.Deployment)
	for _, v := range cur {
		now[v.Name] = v.DeepCopy()
	}
	for _, v := range ori {
		old[v.Name] = v.DeepCopy()
	}
	add := make([]*appsv1.Deployment, 0)
	up := make([]*appsv1.Deployment, 0)
	del := make([]*appsv1.Deployment, 0)
	// add & update
	for k, v := range now {
		if t, ok := old[k]; !ok {
			add = append(add, v.DeepCopy())
			zaplogger.Sugar().Infow("diffDeployments add",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		} else {
			// compare change
			code, changed := compareDeployment(v, t)
			if code == 2 {
				// delete old
				del = append(del, t.DeepCopy())
				// add new
				add = append(add, v.DeepCopy())
			} else {
				// only update
				if changed {
					up = append(up, v.DeepCopy())
					zaplogger.Sugar().Infow("diffDeployments update",
						zap.String("controller", controllerName),
						zap.String("name", k), zap.String("namespace", v.Namespace), zap.Int("code", code))
				}
			}
		}
	}
	// del
	for k, v := range old {
		if _, ok := now[k]; !ok {
			del = append(del, v.DeepCopy())
			zaplogger.Sugar().Infow("diffDeployments delete",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		}
	}
	return add, up, del
}
