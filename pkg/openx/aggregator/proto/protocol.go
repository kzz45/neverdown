package proto

import rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"

type Verb string

const (
	VerbCreate Verb = "create"
	VerbDelete Verb = "delete"
	VerbGet    Verb = "get"
	VerbList   Verb = "list"
	VerbPatch  Verb = "patch"
	VerbUpdate Verb = "update"
	VerbWatch  Verb = "watch"
	VerbPing   Verb = "ping"
)

const (
	VerbPodsSSH          Verb = "SSH"
	VerbPodsLogDownload  Verb = "LogDownload"
	VerbPodsLogStreaming Verb = "LogStreaming"
)

type Request struct {
	GroupVersionKind rbacv1.GroupVersionKind `protobuf:"bytes,1,opt,name=groupVersionKind"`
	Namespace        string                  `protobuf:"bytes,2,opt,name=namespace"`
	Verb             Verb                    `protobuf:"bytes,3,opt,name=verb,casttype=Verb"`
	Raw              []byte                  `protobuf:"bytes,4,opt,name=raw"`
}

type Response struct {
	Code             int32                   `protobuf:"varint,1,opt,name=code"`
	GroupVersionKind rbacv1.GroupVersionKind `protobuf:"bytes,2,opt,name=groupVersionKind"`
	Namespace        string                  `protobuf:"bytes,3,opt,name=namespace"`
	Verb             Verb                    `protobuf:"bytes,4,opt,name=verb,casttype=Verb"`
	Raw              []byte                  `protobuf:"bytes,5,opt,name=raw"`
}

// EventType defines the possible types of events.
type EventType string

const (
	EventAdded    EventType = "ADDED"
	EventModified EventType = "MODIFIED"
	EventDeleted  EventType = "DELETED"
	EventBookmark EventType = "BOOKMARK"
	EventError    EventType = "ERROR"
)

type WatchEvent struct {
	Type string `protobuf:"bytes,1,opt,name=type"`
	Raw  []byte `protobuf:"bytes,2,opt,name=raw"`
}

type Context struct {
	Token       string             `protobuf:"bytes,1,opt,name=token"`
	IsAdmin     bool               `protobuf:"varint,2,opt,name=isAdmin"`
	ExpireAt    int32              `protobuf:"varint,3,opt,name=expireAt"`
	ClusterRole rbacv1.ClusterRole `protobuf:"bytes,4,opt,name=clusterRole"`
}
