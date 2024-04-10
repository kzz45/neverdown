package app

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/authx/validation"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/utils/random"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ServiceAccount interface {
	AppId() string
	List() (*rbacv1.AppServiceAccountList, error)
	Create(ctx context.Context, sa *rbacv1.AppServiceAccount) error
	Get(name string) (*rbacv1.AppServiceAccount, error)
	Update(ctx context.Context, sa *rbacv1.AppServiceAccount) error
	Delete(ctx context.Context, name string) error
	Validate(username, password string) bool
}

type ServiceAccounts struct {
	ctx       context.Context
	clientSet kubernetes.Interface
	appId     string
	Role      Role
	items     cmap.ConcurrentMap
}

func NewAppServiceAccounts(ctx context.Context, cfg kubernetes.Interface, role Role) *ServiceAccounts {
	asa := &ServiceAccounts{
		ctx:       ctx,
		clientSet: cfg,
		appId:     role.AppId(),
		Role:      role,
		items:     cmap.New(),
	}
	go asa.watch()
	return asa
}

func (asa *ServiceAccounts) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := asa.clientSet.RbacV1().AppServiceAccounts(asa.appId).Watch(asa.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-asa.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = asa.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (asa *ServiceAccounts) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Added:
		obj := e.Object.(*rbacv1.AppServiceAccount)
		err = asa.sync(obj)
	case watch.Modified:
		obj := e.Object.(*rbacv1.AppServiceAccount)
		err = asa.sync(obj)
	case watch.Deleted:
		obj := e.Object.(*rbacv1.AppServiceAccount)
		err = asa.delete(obj)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (asa *ServiceAccounts) sync(in *rbacv1.AppServiceAccount) error {
	t, ok := asa.items.Get(in.Name)
	if ok {
		sa := t.(*rbacv1.AppServiceAccount)
		if sa.UID == in.UID {
			if sa.ResourceVersion == in.ResourceVersion {
				return nil
			}
		}
	}
	asa.items.Set(in.Name, in)
	return nil
}

func (asa *ServiceAccounts) delete(in *rbacv1.AppServiceAccount) error {
	asa.items.Remove(in.Name)
	return nil
}

func (asa *ServiceAccounts) AppId() string {
	return asa.appId
}

func (asa *ServiceAccounts) List() (*rbacv1.AppServiceAccountList, error) {
	res := &rbacv1.AppServiceAccountList{
		Items: make([]rbacv1.AppServiceAccount, 0),
	}
	for _, v := range asa.items.Items() {
		cr := v.(*rbacv1.AppServiceAccount)
		res.Items = append(res.Items, *cr)
	}
	return res, nil
}

func (asa *ServiceAccounts) validateRoleRef(ref rbacv1.RoleRef) error {
	// todo is it necessary to specify RoleRef
	if ref.ClusterRoleName == "" {
		return nil
	}
	_, err := asa.Role.Get(ref.ClusterRoleName)
	return err
}

func (asa *ServiceAccounts) Create(ctx context.Context, sa *rbacv1.AppServiceAccount) error {
	sa.Namespace = asa.appId
	if err := asa.validateRoleRef(sa.Spec.RoleRef); err != nil {
		return err
	}
	if err := validation.Account(sa.Spec.AccountMeta.Username); err != nil {
		return err
	}
	sa.Name = sa.Spec.AccountMeta.Username
	sa.Spec.AccountMeta.Password = random.GenRandomString(32)
	_, err := asa.clientSet.RbacV1().AppServiceAccounts(asa.appId).Create(ctx, sa, metav1.CreateOptions{})
	return err
}

func (asa *ServiceAccounts) Get(name string) (*rbacv1.AppServiceAccount, error) {
	t, ok := asa.items.Get(name)
	if !ok {
		return nil, errors.NewNotFound(rbacv1.Resource("AppServiceAccount"), name)
	}
	cr := t.(*rbacv1.AppServiceAccount)
	return cr, nil
}

func (asa *ServiceAccounts) Update(ctx context.Context, sa *rbacv1.AppServiceAccount) error {
	sa.Namespace = asa.appId
	if err := asa.validateRoleRef(sa.Spec.RoleRef); err != nil {
		return err
	}
	if err := validation.Account(sa.Spec.AccountMeta.Username); err != nil {
		return err
	}
	sa.Name = sa.Spec.AccountMeta.Username
	t, ok := asa.items.Get(sa.Name)
	if !ok {
		return errors.NewNotFound(rbacv1.Resource("AppServiceAccount"), sa.Name)
	}
	saCache := t.(*rbacv1.AppServiceAccount)
	saCache.Spec.Desc = sa.Spec.Desc
	saCache.Spec.RoleRef = sa.Spec.RoleRef
	saCache.Spec.AccountMeta.Nickname = sa.Spec.AccountMeta.Nickname
	_, err := asa.clientSet.RbacV1().AppServiceAccounts(asa.appId).Update(ctx, saCache, metav1.UpdateOptions{})
	return err
}

func (asa *ServiceAccounts) Delete(ctx context.Context, name string) error {
	if _, ok := asa.items.Get(name); !ok {
		return errors.NewNotFound(rbacv1.Resource("AppServiceAccount"), name)
	}
	return asa.clientSet.RbacV1().AppServiceAccounts(asa.appId).Delete(ctx, name, metav1.DeleteOptions{})
}

func (asa *ServiceAccounts) Validate(username, password string) bool {
	if err := validation.Account(username); err != nil {
		return false
	}
	t, ok := asa.items.Get(username)
	if !ok {
		return false
	}
	sa := t.(*rbacv1.AppServiceAccount)
	return password == sa.Spec.AccountMeta.Password
}
