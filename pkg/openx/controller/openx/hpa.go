package openx

import (
	"reflect"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func (tc *OpenxController) newHorizontalPodAutoscaler(m *openxv1.Openx, app openxv1.App) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	var labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: m.Name,
	}
	return &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: nil,
			Name:        controller.OpenxHpaName(m.Name, app.AppName),
			Namespace:   m.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m, controllerKind),
			},
			Labels: labels,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       controller.OpenxAppName(m.Name, app.AppName),
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			MinReplicas: app.HorizontalPodAutoscalerSpec.MinReplicas,
			MaxReplicas: app.HorizontalPodAutoscalerSpec.MaxReplicas,
			Metrics:     app.HorizontalPodAutoscalerSpec.Metrics,
			Behavior:    app.HorizontalPodAutoscalerSpec.Behavior,
		},
	}, nil
}

func compareHorizontalPodAutoscaler(s1, s2 *autoscalingv2.HorizontalPodAutoscaler) (int, bool) {
	// name
	if !reflect.DeepEqual(s1.Name, s2.Name) {
		return 1, true
	}
	if !reflect.DeepEqual(s1.Spec.ScaleTargetRef, s2.Spec.ScaleTargetRef) {
		return 2, true
	}
	if !reflect.DeepEqual(s1.Spec.ScaleTargetRef, s2.Spec.ScaleTargetRef) {
		return 3, true
	}
	if !apiequality.Semantic.DeepEqual(s1.Spec.MinReplicas, s2.Spec.MinReplicas) {
		return 4, true
	}
	if !apiequality.Semantic.DeepEqual(s1.Spec.MaxReplicas, s2.Spec.MaxReplicas) {
		return 5, true
	}
	if len(s1.Spec.Metrics) != len(s2.Spec.Metrics) {
		return 6, true
	}
	for _, m1 := range s1.Spec.Metrics {
		exist := false
		for _, m2 := range s2.Spec.Metrics {
			if !apiequality.Semantic.DeepEqual(m1.Type, m2.Type) {
				continue
			}
			if !apiequality.Semantic.DeepEqual(m1.Object, m2.Object) {
				continue
			}
			if !apiequality.Semantic.DeepEqual(m1.Pods, m2.Pods) {
				continue
			}
			if !apiequality.Semantic.DeepEqual(m1.Resource, m2.Resource) {
				continue
			}
			if !apiequality.Semantic.DeepEqual(m1.ContainerResource, m2.ContainerResource) {
				continue
			}
			if !apiequality.Semantic.DeepEqual(m1.External, m2.External) {
				continue
			}
			exist = true
		}
		if !exist {
			return 7, true
		}
	}
	if !apiequality.Semantic.DeepEqual(s1.Spec.Behavior, s2.Spec.Behavior) {
		return 8, true
	}
	return 9, false
}

func diffHorizontalPodAutoscalers(cur, ori []*autoscalingv2.HorizontalPodAutoscaler) ([]*autoscalingv2.HorizontalPodAutoscaler, []*autoscalingv2.HorizontalPodAutoscaler, []*autoscalingv2.HorizontalPodAutoscaler) {
	now := make(map[string]*autoscalingv2.HorizontalPodAutoscaler)
	old := make(map[string]*autoscalingv2.HorizontalPodAutoscaler)
	for _, v := range cur {
		now[v.Name] = v.DeepCopy()
	}
	for _, v := range ori {
		old[v.Name] = v.DeepCopy()
	}
	add := make([]*autoscalingv2.HorizontalPodAutoscaler, 0)
	up := make([]*autoscalingv2.HorizontalPodAutoscaler, 0)
	del := make([]*autoscalingv2.HorizontalPodAutoscaler, 0)
	// add & update
	for k, v := range now {
		if t, ok := old[k]; !ok {
			add = append(add, v.DeepCopy())
			zaplogger.Sugar().Infow("diffHorizontalPodAutoscalers add",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		} else {
			// compare change
			if code, changed := compareHorizontalPodAutoscaler(v, t); changed {
				pre := v.DeepCopy()
				pre.ResourceVersion = t.ResourceVersion
				up = append(up, pre)
				zaplogger.Sugar().Infow("diffHorizontalPodAutoscalers update",
					zap.String("controller", controllerName),
					zap.String("name", k), zap.String("namespace", v.Namespace), zap.Int("code", code))

			}
		}
	}
	// del
	for k, v := range old {
		if _, ok := now[k]; !ok {
			del = append(del, v.DeepCopy())
			zaplogger.Sugar().Infow("diffHorizontalPodAutoscalers delete",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		}
	}
	return add, up, del
}

func (tc *OpenxController) syncHorizontalPodAutoscalers(m *openxv1.Openx, hpaList []*autoscalingv2.HorizontalPodAutoscaler) error {
	cur := make([]*autoscalingv2.HorizontalPodAutoscaler, 0)
	for _, app := range m.Spec.Applications {
		if app.HorizontalPodAutoscalerSpec == nil {
			continue
		}
		hpa, err := tc.newHorizontalPodAutoscaler(m, app)
		if err != nil {
			return err
		}
		cur = append(cur, hpa)
	}
	var errlist []error
	add, update, del := diffHorizontalPodAutoscalers(cur, hpaList)
	for _, v := range add {
		if _, err := tc.hpaControl.Create(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range update {
		if _, err := tc.hpaControl.Update(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range del {
		if err := tc.hpaControl.Delete(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	return utilerrors.NewAggregate(errlist)
}

type realHpaReplicas struct {
	items map[string]*int32
}

func hpaScaleReplicas(hpaList []*autoscalingv2.HorizontalPodAutoscaler) *realHpaReplicas {
	ohr := &realHpaReplicas{
		items: make(map[string]*int32, 0),
	}
	for _, v := range hpaList {
		a := v.DeepCopy()
		ohr.items[v.Name] = &a.Status.DesiredReplicas
	}
	return ohr
}

func (ohr *realHpaReplicas) realHpaDesiredReplicas(name string) (*int32, bool) {
	if t, ok := ohr.items[name]; ok {
		return t, true
	}
	return nil, false
}
