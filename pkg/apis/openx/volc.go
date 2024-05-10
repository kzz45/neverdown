package openx

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type VolcLoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              VolcLoadBalancerSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// 火山云的负载均衡
type VolcLoadBalancerSpec struct {
	Instance          CloudLoadBalancerObject `protobuf:"bytes,1,opt,name=instance"`
	OverrideListeners CloudLoadBalancerObject `protobuf:"bytes,2,opt,name=overrideListeners"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VolcLoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VolcLoadBalancer `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type VolcAccessControl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              VolcAccessControlSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// 火山云的访问控制
type VolcAccessControlSpec struct {
	Instance CloudLoadBalancerObject `protobuf:"bytes,1,opt,name=instance"`
	Status   CloudLoadBalancerObject `protobuf:"bytes,2,opt,name=status"`
	Type     CloudLoadBalancerObject `protobuf:"bytes,3,opt,name=type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VolcAccessControlList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VolcAccessControl `json:"items" protobuf:"bytes,2,rep,name=items"`
}
