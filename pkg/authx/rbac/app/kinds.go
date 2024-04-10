package app

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/authx/validation"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Kind interface {
	AppId() string
	List() (*rbacv1.GroupVersionKindRuleList, error)
	Create(ctx context.Context, rule *rbacv1.GroupVersionKindRule) error
	Update(ctx context.Context, rule *rbacv1.GroupVersionKindRule) error
	Delete(ctx context.Context, name string) error
	Validate(rule rbacv1.GroupVersionKind, verb string) bool
}

type GroupVersionKinds struct {
	ctx       context.Context
	clientSet kubernetes.Interface
	appId     string
	items     cmap.ConcurrentMap
	verbs     cmap.ConcurrentMap
}

func NewGroupVersionKinds(ctx context.Context, cfg kubernetes.Interface, appId string) *GroupVersionKinds {
	r := &GroupVersionKinds{
		ctx:       ctx,
		clientSet: cfg,
		appId:     appId,
		items:     cmap.New(),
		verbs:     cmap.New(),
	}
	go r.watch()
	return r
}

func GroupVersionKindToString(gvk rbacv1.GroupVersionKind) string {
	return gvk.Group + "." + gvk.Version + ".Kind=" + gvk.Kind
}

func GroupVersionKindWithVerbToString(gvk rbacv1.GroupVersionKind, verb string) string {
	return gvk.Group + "." + gvk.Version + ".Kind=" + gvk.Kind + ".Verb=" + verb
}

func (gvk *GroupVersionKinds) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := gvk.clientSet.RbacV1().GroupVersionKindRules(gvk.appId).Watch(gvk.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-gvk.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = gvk.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (gvk *GroupVersionKinds) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Added:
		obj := e.Object.(*rbacv1.GroupVersionKindRule)
		err = gvk.sync(obj)
	case watch.Modified:
		obj := e.Object.(*rbacv1.GroupVersionKindRule)
		err = gvk.sync(obj)
	case watch.Deleted:
		obj := e.Object.(*rbacv1.GroupVersionKindRule)
		err = gvk.delete(obj)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

// sync supports add and modify
func (gvk *GroupVersionKinds) sync(in *rbacv1.GroupVersionKindRule) error {
	var empty, add, del []string
	t, ok := gvk.items.Get(GroupVersionKindToString(in.Spec.GroupVersionKind))
	if ok {
		ori := t.(*rbacv1.GroupVersionKindRule)
		if ori.UID == in.UID {
			if ori.ResourceVersion == in.ResourceVersion {
				return nil
			}
		}
		add, del = diffVerbs(ori.Spec.Verbs, in.Spec.Verbs)
	} else {
		add, del = diffVerbs(empty, in.Spec.Verbs)
	}
	gvk.diff(in.Spec.GroupVersionKind, add, del)
	gvk.items.Set(GroupVersionKindToString(in.Spec.GroupVersionKind), in)
	return nil
}

func (gvk *GroupVersionKinds) delete(in *rbacv1.GroupVersionKindRule) error {
	var empty, add, del []string
	t, ok := gvk.items.Get(GroupVersionKindToString(in.Spec.GroupVersionKind))
	if ok {
		ori := t.(*rbacv1.GroupVersionKindRule)
		ori.Spec.Verbs = append(ori.Spec.Verbs, in.Spec.Verbs...)
		add, del = diffVerbs(ori.Spec.Verbs, empty)
	}
	gvk.diff(in.Spec.GroupVersionKind, add, del)
	gvk.items.Remove(GroupVersionKindToString(in.Spec.GroupVersionKind))
	return nil
}

func (gvk *GroupVersionKinds) diff(in rbacv1.GroupVersionKind, add []string, del []string) {
	for _, v := range add {
		gvk.verbs.Set(GroupVersionKindWithVerbToString(in, v), 1)
	}
	for _, v := range del {
		gvk.verbs.Remove(GroupVersionKindWithVerbToString(in, v))
	}
}

func diffVerbs(ori []string, ups []string) ([]string, []string) {
	add := make([]string, 0)
	del := make([]string, 0)
	for _, v1 := range ori {
		exist := false
		for _, v2 := range ups {
			if v2 == v1 {
				exist = true
				break
			}
		}
		if !exist {
			del = append(del, v1)
		}
	}
	for _, v1 := range ups {
		exist := false
		for _, v2 := range ori {
			if v2 == v1 {
				exist = true
				break
			}
		}
		if !exist {
			add = append(add, v1)
		}
	}
	return add, del
}

func (gvk *GroupVersionKinds) AppId() string {
	return gvk.appId
}

func (gvk *GroupVersionKinds) List() (*rbacv1.GroupVersionKindRuleList, error) {
	res := &rbacv1.GroupVersionKindRuleList{
		Items: make([]rbacv1.GroupVersionKindRule, 0),
	}
	for _, v := range gvk.items.Items() {
		tmp := v.(*rbacv1.GroupVersionKindRule)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func checkConflict(verbs []string) error {
	exist := make(map[string]bool, 1)
	for _, v := range verbs {
		if _, ok := exist[v]; ok {
			return fmt.Errorf("repeated verb:%s", v)
		}
		exist[v] = true
	}
	return nil
}

func (gvk *GroupVersionKinds) Create(ctx context.Context, rule *rbacv1.GroupVersionKindRule) error {
	rule.Namespace = gvk.appId
	rule.Name = GroupVersionKindToString(rule.Spec.GroupVersionKind)
	if err := checkConflict(rule.Spec.Verbs); err != nil {
		return err
	}
	for _, v := range rule.Spec.Verbs {
		if err := validation.Verb(v); err != nil {
			return err
		}
	}
	_, err := gvk.clientSet.RbacV1().GroupVersionKindRules(rule.Namespace).Create(ctx, rule, metav1.CreateOptions{})
	return err
}

func (gvk *GroupVersionKinds) Update(ctx context.Context, rule *rbacv1.GroupVersionKindRule) error {
	rule.Namespace = gvk.appId
	rule.Name = GroupVersionKindToString(rule.Spec.GroupVersionKind)
	if err := checkConflict(rule.Spec.Verbs); err != nil {
		return err
	}
	if _, ok := gvk.items.Get(rule.Name); !ok {
		return errors.NewNotFound(rbacv1.Resource("GroupVersionKindRule"), rule.Name)
	}
	for _, v := range rule.Spec.Verbs {
		if err := validation.Verb(v); err != nil {
			return err
		}
	}
	_, err := gvk.clientSet.RbacV1().GroupVersionKindRules(rule.Namespace).Update(ctx, rule, metav1.UpdateOptions{})
	return err
}

func (gvk *GroupVersionKinds) Delete(ctx context.Context, name string) error {
	if _, ok := gvk.items.Get(name); !ok {
		return errors.NewNotFound(rbacv1.Resource("GroupVersionKindRule"), name)
	}
	return gvk.clientSet.RbacV1().GroupVersionKindRules(gvk.appId).Delete(ctx, name, metav1.DeleteOptions{})
}

// Validate will check whether the request has the permission to take the specified action
func (gvk *GroupVersionKinds) Validate(rule rbacv1.GroupVersionKind, verb string) bool {
	_, ok := gvk.verbs.Get(GroupVersionKindWithVerbToString(rule, verb))
	return ok
}
