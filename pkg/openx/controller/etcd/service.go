package etcd

import (
	"reflect"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func (ec *EtcdController) newService(e *openxv1.Etcd) (*corev1.Service, error) {
	var labels = map[string]string{
		controller.LabelKind:       controllerKind.Kind,
		controller.LabelController: e.Name,
	}
	ports := []corev1.ServicePort{
		{
			Port:     ClientPort,
			Name:     "client",
			Protocol: corev1.ProtocolTCP,
		},
		{
			Port:     PeersPort,
			Name:     "peer",
			Protocol: corev1.ProtocolTCP,
		},
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controller.EtcdServiceName(e.Name),
			Namespace: e.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(e, controllerKind),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Type:      corev1.ServiceTypeClusterIP,
			Ports:     ports,
			Selector:  labels,
		},
	}, nil
}

func (ec *EtcdController) syncServices(e *openxv1.Etcd, svcList []*corev1.Service) error {
	cur := make([]*corev1.Service, 0)
	svc, err := ec.newService(e)
	if err != nil {
		return err
	}
	cur = append(cur, svc)
	var errlist []error
	add, update, del := diffServices(cur, svcList)
	for _, v := range add {
		if _, err := ec.svcControl.Create(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range update {
		if _, err := ec.svcControl.Update(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	for _, v := range del {
		if err := ec.svcControl.Delete(v); err != nil {
			errlist = append(errlist, err)
		}
	}
	return utilerrors.NewAggregate(errlist)
}

func compareService(s1, s2 *corev1.Service) (int, bool) {
	// name
	if !reflect.DeepEqual(s1.Name, s2.Name) {
		return 1, true
	}
	// annotation
	if s1.Annotations == nil {
		s1.Annotations = make(map[string]string)
	}
	if s2.Annotations == nil {
		s2.Annotations = make(map[string]string)
	}
	if !reflect.DeepEqual(s1.Annotations, s2.Annotations) {
		return 2, true
	}
	if s1.Spec.Type != s2.Spec.Type {
		return 3, true
	}
	if len(s1.Spec.Ports) != len(s2.Spec.Ports) {
		return 4, true
	}
	for _, v := range s1.Spec.Ports {
		exist := false
		for _, v2 := range s2.Spec.Ports {
			if v.Port == v2.Port {
				exist = true
				if v.Name != v2.Name {
					return 5, true
				}
				if v.NodePort != v2.NodePort && v2.NodePort != 0 {
					return 6, true
				}
				if v.Protocol != v2.Protocol {
					return 7, true
				}
				if v.TargetPort.Type != v2.TargetPort.Type {
					return 8, true
				}
				if v.TargetPort.IntVal != v2.TargetPort.IntVal {
					return 9, true
				}
				if v.TargetPort.StrVal != v2.TargetPort.StrVal {
					return 10, true
				}
				break
			}
		}
		if !exist {
			return 11, true
		}
	}
	return 12, false
}

func diffServices(cur, ori []*corev1.Service) ([]*corev1.Service, []*corev1.Service, []*corev1.Service) {
	now := make(map[string]*corev1.Service)
	old := make(map[string]*corev1.Service)
	for _, v := range cur {
		now[v.Name] = v.DeepCopy()
	}
	for _, v := range ori {
		old[v.Name] = v.DeepCopy()
	}
	add := make([]*corev1.Service, 0)
	up := make([]*corev1.Service, 0)
	del := make([]*corev1.Service, 0)
	// add & update
	for k, v := range now {
		if t, ok := old[k]; !ok {
			add = append(add, v.DeepCopy())
			zaplogger.Sugar().Infow("diffServices add",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		} else {
			// compare change
			if code, changed := compareService(v, t); changed {
				pre := v.DeepCopy()
				pre.ResourceVersion = t.ResourceVersion
				up = append(up, pre)
				zaplogger.Sugar().Infow("diffServices update",
					zap.String("controller", controllerName),
					zap.String("name", k), zap.String("namespace", v.Namespace), zap.Int("code", code))

			}
		}
	}
	// del
	for k, v := range old {
		if _, ok := now[k]; !ok {
			del = append(del, v.DeepCopy())
			zaplogger.Sugar().Infow("diffServices delete",
				zap.String("controller", controllerName),
				zap.String("name", k), zap.String("namespace", v.Namespace))
		}
	}
	return add, up, del
}
