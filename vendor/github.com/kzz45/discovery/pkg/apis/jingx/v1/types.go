package v1

import (
	rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

type Project struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta `protobuf:"bytes,1,opt,name=metadata"`

	Spec ProjectSpec `protobuf:"bytes,2,opt,name=spec"`
}

type ProjectSpec struct {
	GenerateId string   `protobuf:"bytes,1,opt,name=generateId"`
	Domains    []string `protobuf:"bytes,2,rep,name=domains"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProjectList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta `protobuf:"bytes,1,opt,name=metadata"`
	Items           []Project `protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

type Repository struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta `protobuf:"bytes,1,opt,name=metadata"`

	Spec RepositorySpec `protobuf:"bytes,2,opt,name=spec"`
}

type RepositoryMeta struct {
	// ProjectName was the Project's object name
	ProjectName string `protobuf:"bytes,1,opt,name=projectName"`
	// RepositoryName was the repository name, it wasn't the final Repository's object name
	// The final Repository's object name should be ProjectName + RepositoryName
	RepositoryName string `protobuf:"bytes,2,opt,name=repositoryName"`
}

type RepositorySpec struct {
	RepositoryMeta RepositoryMeta `protobuf:"bytes,1,opt,name=repositoryMeta"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RepositoryList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta `protobuf:"bytes,1,opt,name=metadata"`
	Items           []Repository `protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

type Tag struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta `protobuf:"bytes,1,opt,name=metadata"`

	Spec TagSpec `protobuf:"bytes,2,opt,name=spec"`
}

type TagSpec struct {
	RepositoryMeta RepositoryMeta `protobuf:"bytes,1,opt,name=repositoryMeta"`
	Tag            string         `protobuf:"bytes,2,opt,name=tag"`
	GitReference   GitReference   `protobuf:"bytes,3,opt,name=gitReference"`
	DockerImage    DockerImage    `protobuf:"bytes,4,opt,name=dockerImage"`
}

type GitReference struct {
	Git        string `protobuf:"bytes,1,opt,name=git"`
	Branch     string `protobuf:"bytes,2,opt,name=branch"`
	CommitHash string `protobuf:"bytes,3,opt,name=commitHash"`
}

type DockerImage struct {
	Sha256           string `protobuf:"bytes,1,opt,name=sha256"`
	Author           string `protobuf:"bytes,2,opt,name=author"`
	LastModifiedTime int64  `protobuf:"varint,3,opt,name=lastModifiedTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TagList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta `protobuf:"bytes,1,opt,name=metadata"`
	Items           []Tag `protobuf:"bytes,2,rep,name=items"`
}

type Verb string

const (
	VerbCreate Verb = "create"
	VerbUpdate Verb = "update"
	VerbDelete Verb = "delete"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

type Event struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta `protobuf:"bytes,1,opt,name=metadata"`

	Spec EventSpec `protobuf:"bytes,2,opt,name=spec"`
}

type EventSpec struct {
	Author           string                  `protobuf:"bytes,1,opt,name=author"`
	GroupVersionKind rbacv1.GroupVersionKind `protobuf:"bytes,2,opt,name=groupVersionKind"`
	Verb             Verb                    `protobuf:"bytes,3,opt,name=verb,casttype=Verb"`
	Raw              []byte                  `protobuf:"bytes,4,opt,name=raw"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EventList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta `protobuf:"bytes,1,opt,name=metadata"`
	Items           []Event `protobuf:"bytes,2,rep,name=items"`
}
