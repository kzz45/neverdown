package openx

import (
	"context"
	"reflect"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (tc *OpenxController) syncStatusOnly(m *openxv1.Openx, dpList []*appsv1.Deployment, hpaList []*autoscalingv2.HorizontalPodAutoscaler) error {
	status := openxv1.OpenxStatus{
		Items: make(map[string]openxv1.AppStatus, 0),
	}
	for _, dp := range dpList {
		status.Items[dp.Name] = openxv1.AppStatus{
			DeploymentStatus:              dp.Status,
			HorizontalPodAutoscalerStatus: autoscalingv2.HorizontalPodAutoscalerStatus{},
		}
	}
	for _, hpa := range hpaList {
		if t, ok := status.Items[hpa.Name]; ok {
			t.HorizontalPodAutoscalerStatus = hpa.Status
			status.Items[hpa.Name] = t
		} else {
			status.Items[hpa.Name] = openxv1.AppStatus{
				DeploymentStatus:              appsv1.DeploymentStatus{},
				HorizontalPodAutoscalerStatus: hpa.Status,
			}
		}
	}
	return tc.syncOpenxStatus(status, m)
}

func (tc *OpenxController) syncOpenxStatus(status openxv1.OpenxStatus, t *openxv1.Openx) error {
	if reflect.DeepEqual(t.Status, status) {
		return nil
	}
	newOpenx := t
	newOpenx.Status = status
	_, err := tc.openx.OpenxV1().Openxes(newOpenx.Namespace).UpdateStatus(context.TODO(), newOpenx, metav1.UpdateOptions{})
	return err
}

func (tc *OpenxController) forceSyncOpenx(t *openxv1.Openx) error {
	_, err := tc.openx.OpenxV1().Openxes(t.Namespace).Update(context.TODO(), t, metav1.UpdateOptions{})
	return err
}
