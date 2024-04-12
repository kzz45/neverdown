package syncer

import (
	"context"
	"fmt"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/env"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/proto"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

type Syncer struct {
	ctx   context.Context
	cfg   kubernetes.Interface
	rules []string

	queues chan *queueElement
}

type queueElement struct {
	eventType proto.EventType
	obj       runtime.Object
	gvk       rbacv1.GroupVersionKind
}

func New(ctx context.Context, cfg kubernetes.Interface) *Syncer {
	s := &Syncer{
		ctx:    ctx,
		cfg:    cfg,
		rules:  env.DomainRefactorRules(),
		queues: make(chan *queueElement, 40960),
	}
	go s.consumer()
	return s
}

func errNotNil(err error) bool {
	if err != nil {
		return true
	}
	return false
}

func (s *Syncer) consumer() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg, ok := <-s.queues:
			if !ok {
				return
			}
			s.retry(msg)
		}
	}
}

func (s *Syncer) retry(event *queueElement) {
	zaplogger.Sugar().Infow("retry", zap.Any("gvk", event.gvk), zap.Any("type", event.eventType), zap.Any("name", event.obj.(metav1.Object).GetName()))
	defaultBackoff := wait.Backoff{
		Steps:    10,
		Duration: 200 * time.Millisecond,
		Factor:   5.0,
		Jitter:   0.1,
	}
	err := retry.OnError(defaultBackoff, errNotNil, func() error {
		return s.sync(event)
	})
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
}

func (s *Syncer) sync(event *queueElement) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*3)
	defer cancel()
	var err error
	switch event.gvk {
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		obj := event.obj.(*jingxv1.Project)
		switch event.eventType {
		case proto.EventAdded:
			_, err = s.cfg.JingxV1().Projects(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Projects(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		case proto.EventDeleted:
			err = s.cfg.JingxV1().Projects(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
		case proto.EventModified:
			_, err = s.cfg.JingxV1().Projects(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Projects(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		}
		if err != nil {
			zaplogger.Sugar().Error(err)
			return err
		}
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		obj := event.obj.(*jingxv1.Repository)
		switch event.eventType {
		case proto.EventAdded:
			_, err = s.cfg.JingxV1().Repositories(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Repositories(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		case proto.EventDeleted:
			err = s.cfg.JingxV1().Repositories(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
		case proto.EventModified:
			_, err = s.cfg.JingxV1().Repositories(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Repositories(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		}
		if err != nil {
			zaplogger.Sugar().Error(err)
			return err
		}
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		obj := event.obj.(*jingxv1.Tag)
		switch event.eventType {
		case proto.EventAdded:
			_, err = s.cfg.JingxV1().Tags(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Tags(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		case proto.EventDeleted:
			err = s.cfg.JingxV1().Tags(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
		case proto.EventModified:
			_, err = s.cfg.JingxV1().Tags(obj.Namespace).Create(ctx, obj, metav1.CreateOptions{})
			if errors.IsAlreadyExists(err) {
				_, err = s.cfg.JingxV1().Tags(obj.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
			}
		}
		if err != nil {
			zaplogger.Sugar().Error(err)
			return err
		}
	default:
		err = fmt.Errorf(aggregator.ErrGVKNotExist, event.gvk)
	}
	return nil
}

func (s *Syncer) Handler(in *proto.Response) error {
	if in.Code != 0 {
		return fmt.Errorf("response error code:%d err:%s", in.Code, string(in.Raw))
	}
	if in.Verb == proto.VerbPing {
		return nil
	}
	zaplogger.Sugar().Infow("handle", zap.Any("gvk", in.GroupVersionKind), zap.Any("verb", in.Verb))
	var err error
	switch in.GroupVersionKind {
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		switch in.Verb {
		case proto.VerbList:
			objList := &jingxv1.ProjectList{}
			if err = objList.Unmarshal(in.Raw); err != nil {
				return err
			}
			for _, obj := range objList.Items {
				time.Sleep(time.Millisecond * 200)
				cp := obj.DeepCopy()
				obj.ObjectMeta = metav1.ObjectMeta{
					Name:      cp.Name,
					Namespace: cp.Namespace,
				}
				obj.Spec.Domains = refactorDomains(obj.Spec.Domains, s.rules)
				s.queues <- &queueElement{
					eventType: proto.EventAdded,
					obj:       obj.DeepCopy(),
					gvk:       in.GroupVersionKind,
				}
			}
		case proto.VerbWatch:
			event := &proto.WatchEvent{}
			if err = event.Unmarshal(in.Raw); err != nil {
				return err
			}
			obj := &jingxv1.Project{}
			if err = obj.Unmarshal(event.Raw); err != nil {
				return err
			}
			cp := obj.DeepCopy()
			obj.ObjectMeta = metav1.ObjectMeta{
				Name:      cp.Name,
				Namespace: cp.Namespace,
			}
			obj.Spec.Domains = refactorDomains(obj.Spec.Domains, s.rules)
			s.queues <- &queueElement{
				eventType: proto.EventType(event.Type),
				obj:       obj.DeepCopy(),
				gvk:       in.GroupVersionKind,
			}
		default:
			err = fmt.Errorf("invaild projects handle verb:%s", in.Verb)
		}
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		switch in.Verb {
		case proto.VerbList:
			objList := &jingxv1.RepositoryList{}
			if err = objList.Unmarshal(in.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			for _, obj := range objList.Items {
				time.Sleep(time.Millisecond * 200)
				cp := obj.DeepCopy()
				obj.ObjectMeta = metav1.ObjectMeta{
					Name:      cp.Name,
					Namespace: cp.Namespace,
				}
				s.queues <- &queueElement{
					eventType: proto.EventAdded,
					obj:       obj.DeepCopy(),
					gvk:       in.GroupVersionKind,
				}
			}
		case proto.VerbWatch:
			event := &proto.WatchEvent{}
			if err = event.Unmarshal(in.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			obj := &jingxv1.Repository{}
			if err = obj.Unmarshal(event.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			cp := obj.DeepCopy()
			obj.ObjectMeta = metav1.ObjectMeta{
				Name:      cp.Name,
				Namespace: cp.Namespace,
			}
			s.queues <- &queueElement{
				eventType: proto.EventType(event.Type),
				obj:       obj.DeepCopy(),
				gvk:       in.GroupVersionKind,
			}
		default:
			err = fmt.Errorf("invaild repository handle verb:%s", in.Verb)
		}
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		switch in.Verb {
		case proto.VerbList:
			objList := &jingxv1.TagList{}
			if err = objList.Unmarshal(in.Raw); err != nil {
				return err
			}
			for _, obj := range objList.Items {
				time.Sleep(time.Millisecond * 200)
				cp := obj.DeepCopy()
				obj.ObjectMeta = metav1.ObjectMeta{
					Name:      cp.Name,
					Namespace: cp.Namespace,
				}
				s.queues <- &queueElement{
					eventType: proto.EventAdded,
					obj:       obj.DeepCopy(),
					gvk:       in.GroupVersionKind,
				}
			}
		case proto.VerbWatch:
			event := &proto.WatchEvent{}
			if err = event.Unmarshal(in.Raw); err != nil {
				return err
			}
			obj := &jingxv1.Tag{}
			if err = obj.Unmarshal(event.Raw); err != nil {
				return err
			}
			cp := obj.DeepCopy()
			obj.ObjectMeta = metav1.ObjectMeta{
				Name:      cp.Name,
				Namespace: cp.Namespace,
			}
			s.queues <- &queueElement{
				eventType: proto.EventType(event.Type),
				obj:       obj.DeepCopy(),
				gvk:       in.GroupVersionKind,
			}
		default:
			err = fmt.Errorf("invaild tag handle verb:%s", in.Verb)
		}
	default:
		err = fmt.Errorf(aggregator.ErrGVKNotExist, in.GroupVersionKind)
	}
	return nil
}
