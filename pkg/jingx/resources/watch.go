package resources

import (
	"github.com/kzz45/neverdown/pkg/zaplogger"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/proto"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	DefaultNamespace = ""
)

type WatchEvent struct {
	Namespace        string
	GroupVersionKind rbacv1.GroupVersionKind
	Verb             string
	Raw              *proto.Response
}

func (r *Resource) Watch(eventsChan chan<- *proto.Response) error {
	for gvk, provider := range r.api.Providers() {
		go r.watch(gvk, provider, eventsChan)
	}
	return nil
}

func (r *Resource) watch(gvk schema.GroupVersionKind, provider aggregator.WatchProvider, eventsChan chan<- *proto.Response) {
retry:
	for {
		select {
		case e, isClosed := <-provider.Watcher():
			if !isClosed {
				zaplogger.Sugar().Infof("resource watch resourceType:%#v closed", provider)
				goto retry
			}
			switch e.Type {
			case watch.Modified, watch.Added, watch.Deleted:
			case watch.Error:
				zaplogger.Sugar().Errorw("watch receive ERROR event", "obj", e.Object)
				continue
			}
			res, err := r.watchResponse(gvk, e)
			if err != nil {
				zaplogger.Sugar().Fatal(err)
				continue
			}
			eventsChan <- res
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Resource) watchResponse(gvk schema.GroupVersionKind, e watch.Event) (*proto.Response, error) {
	ns, raw, err := r.api.ConvertObjectToRaw(&gvk, e.Object)
	if err != nil {
		return nil, err
	}
	embed := &proto.WatchEvent{
		Type: string(e.Type),
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
		Namespace: ns,
		Raw:       embeddedRaw,
	}
	return res, nil
}
