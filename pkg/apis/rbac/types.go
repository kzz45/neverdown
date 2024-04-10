package rbac

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// App was an application for rbac
type App struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AppSpec `protobuf:"bytes,2,opt,name=spec"`
}

type AppSpec struct {
	Id     string `protobuf:"bytes,1,opt,name=id"`
	Secret string `protobuf:"bytes,2,opt,name=secret"`
	Desc   string `protobuf:"bytes,3,opt,name=desc"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppList is a list of App resources
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []App `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type AccountMeta struct {
	Username string `protobuf:"bytes,1,opt,name=username"`
	Nickname string `protobuf:"bytes,2,opt,name=nickname"`
	Password string `protobuf:"bytes,3,opt,name=password"`
	Email    string `protobuf:"bytes,4,opt,name=email"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// RbacServiceAccount was the service account which could only be used to
// manage the specified application.
type RbacServiceAccount struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec RbacServiceAccountSpec `protobuf:"bytes,2,opt,name=spec"`
}

type RbacServiceAccountSpec struct {
	AccountMeta AccountMeta `protobuf:"bytes,1,opt,name=accountMeta"`
	Apps        []string    `protobuf:"bytes,2,rep,name=apps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RbacServiceAccountList is a list of RbacServiceAccount resources
type RbacServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []RbacServiceAccount `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Group   string `protobuf:"bytes,1,opt,name=group"`
	Version string `protobuf:"bytes,2,opt,name=version"`
	Kind    string `protobuf:"bytes,3,opt,name=kind"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// GroupVersionKindRule was a single element rule for role policy
type GroupVersionKindRule struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec RuleSpec `protobuf:"bytes,2,opt,name=spec"`
}

type RuleSpec struct {
	GroupVersionKind GroupVersionKind `protobuf:"bytes,1,opt,name=groupVersionKind"`
	Verbs            []string         `protobuf:"bytes,2,rep,name=verbs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupVersionKindRuleList is a list of GroupVersionKindRule resources
type GroupVersionKindRuleList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GroupVersionKindRule `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type PolicyRule struct {
	Namespace        string           `protobuf:"bytes,1,opt,name=namespace"`
	GroupVersionKind GroupVersionKind `protobuf:"bytes,2,opt,name=groupVersionKind"`
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.
	// VerbAll represents all kinds.
	Verbs []string `protobuf:"bytes,3,rep,name=verbs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// ClusterRole was the base role for rbac
type ClusterRole struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec ClusterRoleSpec `protobuf:"bytes,2,opt,name=spec"`
}

type ClusterRoleSpec struct {
	Desc  string       `protobuf:"bytes,1,opt,name=desc"`
	Rules []PolicyRule `protobuf:"bytes,2,rep,name=rules"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRoleList is a list of ClusterRole resources
type ClusterRoleList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterRole `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type RoleRef struct {
	ClusterRoleName string `protobuf:"bytes,1,opt,name=clusterRoleName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// AppServiceAccount was the service account which could only be used inside
// the specified application.
type AppServiceAccount struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec AppServiceAccountSpec `protobuf:"bytes,2,opt,name=spec"`
}

type AppServiceAccountSpec struct {
	Desc        string      `protobuf:"bytes,1,opt,name=desc"`
	RoleRef     RoleRef     `protobuf:"bytes,2,opt,name=roleRef"`
	AccountMeta AccountMeta `protobuf:"bytes,3,opt,name=accountMeta"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppServiceAccountList is a list of AppServiceAccount resources
type AppServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AppServiceAccount `json:"items" protobuf:"bytes,2,rep,name=items"`
}
