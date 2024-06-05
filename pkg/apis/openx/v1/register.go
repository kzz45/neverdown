package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "openx.neverdown.io"

var GroupVersion = metav1.GroupVersion{Group: GroupName, Version: "v1"}

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Etcd{}, &EtcdList{},
		&Openx{}, &OpenxList{},
		&Mysql{}, &MysqlList{},
		&Redis{}, &RedisList{},
		&Affinity{}, &AffinityList{},
		&Toleration{}, &TolerationList{},
		&NodeSelector{}, &NodeSelectorList{},
		&LoadBalancer{}, &LoadBalancerList{},
		&AccessControl{}, &AccessControlList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
