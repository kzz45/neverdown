package registry

import (
	"context"
	"strings"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/utils/random"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type MarshalObject interface {
	Marshal() (dAtA []byte, err error)
}

type RecordEventHandler func(ctx context.Context, author string, gvk rbacv1.GroupVersionKind, verb jingxv1.Verb, obj MarshalObject) error

type Event struct {
	clientSet    kubernetes.Interface
	ctx          context.Context
	namespace    string
	watchChannel chan watch.Event
}

func NewEvent(ctx context.Context, cfg kubernetes.Interface, namespace string) *Event {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	e := &Event{
		clientSet: cfg,
		ctx:       ctx,
		namespace: namespace,
	}
	return e
}

func (e *Event) Record(ctx context.Context, author string, gvk rbacv1.GroupVersionKind, verb jingxv1.Verb, obj MarshalObject) error {
	raw, err := obj.Marshal()
	if err != nil {
		return err
	}
	name := time.Now().Format("20060102150405") + "-" + strings.ToLower(random.GenRandomString(6))
	event := &jingxv1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: e.namespace,
			Name:      name,
		},
		Spec: jingxv1.EventSpec{
			Author:           author,
			GroupVersionKind: gvk,
			Verb:             verb,
			Raw:              raw,
		},
	}
	return e.Create(ctx, event)
}

func (e *Event) Create(ctx context.Context, in *jingxv1.Event) error {
	createOpts := metav1.CreateOptions{}
	_, err := e.clientSet.JingxV1().Events(e.namespace).Create(ctx, in, createOpts)
	return err
}

func (e *Event) List() (*jingxv1.EventList, error) {
	opt := metav1.ListOptions{
		Limit: 1000,
	}
	return e.clientSet.JingxV1().Events(e.namespace).List(e.ctx, opt)
}

func (e *Event) Delete(ctx context.Context, in *jingxv1.Event) error {
	opt := metav1.DeleteOptions{}
	return e.clientSet.JingxV1().Events(e.namespace).Delete(ctx, in.Name, opt)
}

func (e *Event) Update(ctx context.Context, in *jingxv1.Event) error {
	opt := metav1.UpdateOptions{}
	_, err := e.clientSet.JingxV1().Events(e.namespace).Update(ctx, in, opt)
	return err
}

func (e *Event) Watcher() <-chan watch.Event {
	if e.watchChannel != nil {
		return e.watchChannel
	}
	e.watchChannel = make(chan watch.Event, 1024)
	go e.watch()
	return e.watchChannel
}

func (e *Event) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := e.clientSet.JingxV1().Events(e.namespace).Watch(e.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-e.ctx.Done():
			res.Stop()
			return
		case msg, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if e.watchChannel != nil {
				e.watchChannel <- msg
			}
		}
	}
}
