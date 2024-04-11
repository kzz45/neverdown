package redis

import (
	"context"
	"fmt"
	"reflect"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (rc *RedisController) syncStatusOnly(m *openxv1.Redis, stsList []*appsv1.StatefulSet) error {
	status := openxv1.ClusterStatus{
		Master: appsv1.StatefulSetStatus{},
		Slave:  appsv1.StatefulSetStatus{},
	}
	if len(stsList) > 2 {
		return fmt.Errorf("error Redis stsList more than 2")
	}
	for _, sts := range stsList {
		t, ok := sts.Labels[controller.LabelClusterRole]
		if !ok {
			continue
		}
		switch openxv1.ClusterRole(t) {
		case openxv1.ClusterRoleMaster:
			status.Master = sts.Status
		case openxv1.ClusterRoleSlave:
			status.Slave = sts.Status
		}
	}
	return rc.syncRedisStatus(status, m)
}

func (rc *RedisController) syncRedisStatus(status openxv1.ClusterStatus, d *openxv1.Redis) error {
	if reflect.DeepEqual(d.Status, status) {
		return nil
	}
	newRedis := d
	newRedis.Status = status
	_, err := rc.openx.OpenxV1().Redises(newRedis.Namespace).UpdateStatus(context.TODO(), newRedis, metav1.UpdateOptions{})
	return err
}
