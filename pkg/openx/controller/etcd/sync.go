package etcd

import (
	"context"
	"fmt"
	"reflect"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// syncStatusOnly only updates Etcd Status and doesn't take any mutating actions.
func (ec *EtcdController) syncStatusOnly(m *openxv1.Etcd, stsList []*appsv1.StatefulSet) error {
	if len(stsList) == 0 {
		return fmt.Errorf("error Etcd stsList was zero")
	}
	if len(stsList) > 1 {
		return fmt.Errorf("error Etcd stsList more than 1")
	}
	return ec.syncEtcdStatus(stsList[0].Status, m)
}

// syncEtcdStatus checks if the status is up-to-date and sync it if necessary
func (ec *EtcdController) syncEtcdStatus(status appsv1.StatefulSetStatus, d *openxv1.Etcd) error {
	if reflect.DeepEqual(d.Status, status) {
		return nil
	}
	newEtcd := d
	newEtcd.Status = status
	_, err := ec.openx.OpenxV1().Etcds(newEtcd.Namespace).UpdateStatus(context.TODO(), newEtcd, metav1.UpdateOptions{})
	return err
}
