package v1

import (
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Recorder was the record for audit
type Recorder struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec RecorderSpec `protobuf:"bytes,2,opt,name=spec"`
}

type RecorderSpec struct {
	GroupVersionKind rbacv1.GroupVersionKind `protobuf:"bytes,1,opt,name=groupVersionKind"`
	Verb             string                  `protobuf:"bytes,2,opt,name=verb"`
	OriginalData     []byte                  `protobuf:"bytes,3,opt,name=originalData"`
	Data             []byte                  `protobuf:"bytes,4,opt,name=data"`
	Author           string                  `protobuf:"bytes,5,opt,name=author"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RecorderList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Recorder `json:"items" protobuf:"bytes,2,rep,name=items"`
}
