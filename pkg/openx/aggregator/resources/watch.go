package resources

import (
	rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/proto"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func (r *Resources) AddEventHandler() <-chan *proto.Response {
	if r.watchers != nil {
		return r.watchers
	}
	r.watchers = make(chan *proto.Response, 1024)
	return r.watchers
}

func convertNamespace(obj interface{}) string {
	if d, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		ns, _, err := cache.SplitMetaNamespaceKey(d.Key)
		if err != nil {
			zaplogger.Sugar().Fatal(err)
		}
		return ns
	}
	return obj.(metav1.Object).GetNamespace()
}

func (r *Resources) addFunc(object interface{}) {
	if r.watchers == nil {
		return
	}
	res, err := convertEventToProto(r.convertGroupVersionKind(object), convertNamespace(object), object.(runtime.Object), watch.Added)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	r.watchers <- res
}

func (r *Resources) updateFunc(old, cur interface{}) {
	ori := old.(metav1.Object)
	up := cur.(metav1.Object)
	if ori.GetResourceVersion() == up.GetResourceVersion() {
		return
	}
	if r.watchers == nil {
		return
	}
	res, err := convertEventToProto(r.convertGroupVersionKind(cur), convertNamespace(cur), cur.(runtime.Object), watch.Modified)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	r.watchers <- res
}

func (r *Resources) deleteFunc(object interface{}) {
	if r.watchers == nil {
		return
	}
	res, err := convertEventToProto(r.convertGroupVersionKind(object), convertNamespace(object), object.(runtime.Object), watch.Deleted)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	r.watchers <- res
}

func convertEventToProto(gvk schema.GroupVersionKind, namespace string, obj interface{}, et watch.EventType) (*proto.Response, error) {
	mar := obj.(Object)
	raw, err := mar.Marshal()
	if err != nil {
		return nil, err
	}
	embed := &proto.WatchEvent{
		Type: string(et),
		Raw:  raw,
	}
	embeddedRaw, err := embed.Marshal()
	if err != nil {
		return nil, err
	}
	res := &proto.Response{
		Code: 0,
		GroupVersionKind: rbacv1.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		},
		Verb:      proto.VerbWatch,
		Namespace: namespace,
		Raw:       embeddedRaw,
	}
	return res, nil
}
