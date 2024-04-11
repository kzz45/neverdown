package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func (m *Mysql) Recorder() *Mysql {
	t := m.DeepCopy()
	t.TypeMeta = metav1.TypeMeta{
		Kind:       "Mysql",
		APIVersion: SchemeGroupVersion.String(),
	}
	t.Spec = MysqlSpec{}
	return t
}

func (m *Redis) Recorder() *Redis {
	t := m.DeepCopy()
	t.TypeMeta = metav1.TypeMeta{
		Kind:       "Redis",
		APIVersion: SchemeGroupVersion.String(),
	}
	t.Spec = RedisSpec{}
	return t
}

func (m *Etcd) Recorder() *Etcd {
	t := m.DeepCopy()
	t.TypeMeta = metav1.TypeMeta{
		Kind:       "Etcd",
		APIVersion: SchemeGroupVersion.String(),
	}
	t.Spec = EtcdSpec{}
	return t
}

func (m *Openx) Recorder() *Openx {
	t := m.DeepCopy()
	t.TypeMeta = metav1.TypeMeta{
		Kind:       "Openx",
		APIVersion: SchemeGroupVersion.String(),
	}
	t.Spec = OpenxSpec{}
	return t
}
