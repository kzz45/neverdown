package mysql

import (
	"context"
	"fmt"
	"reflect"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// syncStatusOnly only updates Mysql Status and doesn't take any mutating actions.
func (mc *MysqlController) syncStatusOnly(m *openxv1.Mysql, stsList []*appsv1.StatefulSet) error {
	status := openxv1.ClusterStatus{
		Master: appsv1.StatefulSetStatus{},
		Slave:  appsv1.StatefulSetStatus{},
	}
	if len(stsList) > 2 {
		return fmt.Errorf("error Mysql stsList more than 2")
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
	return mc.syncMysqlStatus(status, m)
}

// syncMysqlStatus checks if the status is up-to-date and sync it if necessary
func (mc *MysqlController) syncMysqlStatus(status openxv1.ClusterStatus, d *openxv1.Mysql) error {
	if reflect.DeepEqual(d.Status, status) {
		return nil
	}
	newMysql := d
	newMysql.Status = status
	_, err := mc.openx.OpenxV1().Mysqls(newMysql.Namespace).UpdateStatus(context.TODO(), newMysql, metav1.UpdateOptions{})
	return err
}
