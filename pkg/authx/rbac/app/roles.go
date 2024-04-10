package app

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Role interface {
	AppId() string
	List() (*rbacv1.ClusterRoleList, error)
	Create(ctx context.Context, cr *rbacv1.ClusterRole) error
	Get(name string) (*rbacv1.ClusterRole, error)
	Update(ctx context.Context, cr *rbacv1.ClusterRole) error
	Delete(ctx context.Context, name string) error
	Validate(name string) (*rbacv1.ClusterRole, bool)
}

type Roles struct {
	ctx       context.Context
	clientSet kubernetes.Interface
	appId     string
	kind      Kind
	items     cmap.ConcurrentMap
}

func NewRoles(ctx context.Context, cfg kubernetes.Interface, kind Kind) *Roles {
	r := &Roles{
		ctx:       ctx,
		clientSet: cfg,
		appId:     kind.AppId(),
		kind:      kind,
		items:     cmap.New(),
	}
	go r.watch()
	return r
}

func (rs *Roles) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := rs.clientSet.RbacV1().ClusterRoles(rs.appId).Watch(rs.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-rs.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = rs.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (rs *Roles) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Added:
		obj := e.Object.(*rbacv1.ClusterRole)
		err = rs.sync(obj)
	case watch.Modified:
		obj := e.Object.(*rbacv1.ClusterRole)
		err = rs.sync(obj)
	case watch.Deleted:
		obj := e.Object.(*rbacv1.ClusterRole)
		err = rs.delete(obj)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (rs *Roles) sync(in *rbacv1.ClusterRole) error {
	t, ok := rs.items.Get(in.Name)
	if ok {
		cr := t.(*rbacv1.ClusterRole)
		if cr.UID == in.UID {
			if cr.ResourceVersion == in.ResourceVersion {
				return nil
			}
		}
	}
	rs.items.Set(in.Name, in)
	return nil
}

func (rs *Roles) delete(in *rbacv1.ClusterRole) error {
	rs.items.Remove(in.Name)
	return nil
}

func (rs *Roles) validateRules(rules []rbacv1.PolicyRule) error {
	// todo is it necessary to bind namespace with rules here?
	// maybe there is no need to do anything here.
	duplicated := make(map[string]bool, 0)
	for _, v := range rules {
		name := GroupVersionKindToString(v.GroupVersionKind) + ".Namespace=" + v.Namespace
		if _, ok := duplicated[name]; ok {
			return fmt.Errorf("GroupVersionKind:%v duplicated in different namespace:%s", v.GroupVersionKind, name)
		}
		duplicated[name] = true
		if err := checkConflict(v.Verbs); err != nil {
			return err
		}
		for _, verb := range v.Verbs {
			ok := rs.kind.Validate(v.GroupVersionKind, verb)
			if !ok {
				return fmt.Errorf("GroupVersionKind:%v invalid verb:%s", v.GroupVersionKind, verb)
			}
		}
	}
	return nil
}

func (rs *Roles) AppId() string {
	return rs.appId
}

func (rs *Roles) List() (*rbacv1.ClusterRoleList, error) {
	res := &rbacv1.ClusterRoleList{
		Items: make([]rbacv1.ClusterRole, 0),
	}
	for _, v := range rs.items.Items() {
		cr := v.(*rbacv1.ClusterRole)
		res.Items = append(res.Items, *cr)
	}
	return res, nil
}

func (rs *Roles) Create(ctx context.Context, cr *rbacv1.ClusterRole) error {
	cr.Namespace = rs.appId
	if err := rs.validateRules(cr.Spec.Rules); err != nil {
		return err
	}
	_, err := rs.clientSet.RbacV1().ClusterRoles(rs.appId).Create(ctx, cr, metav1.CreateOptions{})
	return err
}

func (rs *Roles) Get(name string) (*rbacv1.ClusterRole, error) {
	t, ok := rs.items.Get(name)
	if !ok {
		return nil, errors.NewNotFound(rbacv1.Resource("ClusterRole"), name)
	}
	cr := t.(*rbacv1.ClusterRole)
	return cr, nil
}

func (rs *Roles) Update(ctx context.Context, cr *rbacv1.ClusterRole) error {
	t, ok := rs.items.Get(cr.Name)
	if !ok {
		return errors.NewNotFound(rbacv1.Resource("ClusterRole"), cr.Name)
	}
	obj := t.(*rbacv1.ClusterRole).DeepCopy()
	if err := rs.validateRules(cr.Spec.Rules); err != nil {
		return err
	}
	obj.Spec = cr.Spec
	_, err := rs.clientSet.RbacV1().ClusterRoles(rs.appId).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func (rs *Roles) Delete(ctx context.Context, name string) error {
	if _, ok := rs.items.Get(name); !ok {
		return errors.NewNotFound(rbacv1.Resource("ClusterRole"), name)
	}
	return rs.clientSet.RbacV1().ClusterRoles(rs.appId).Delete(ctx, name, metav1.DeleteOptions{})
}

// Validate will check whether the role name was exist,
// and return the ClusterRole
func (rs *Roles) Validate(name string) (*rbacv1.ClusterRole, bool) {
	t, ok := rs.items.Get(name)
	if !ok {
		return nil, false
	}
	return t.(*rbacv1.ClusterRole), true
}
