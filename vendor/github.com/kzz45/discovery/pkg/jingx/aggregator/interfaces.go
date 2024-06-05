package aggregator

import (
	"github.com/kzz45/discovery/pkg/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

type Object interface {
	Marshal() (dAtA []byte, err error)
	Unmarshal(dAtA []byte) error
}

type WatchProvider interface {
	Watcher() <-chan watch.Event
}

type Api interface {
	ClientSet() kubernetes.Interface
	Create(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error)
	Delete(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error)
	Update(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error)
	List(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error)
	ConvertObjectToRaw(gvk *schema.GroupVersionKind, obj runtime.Object) (namespace string, res []byte, err error)
	Providers() map[schema.GroupVersionKind]WatchProvider
}
