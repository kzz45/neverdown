package admin

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/utils/random"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type RbacServiceAccount interface {
	List() (*rbacv1.RbacServiceAccountList, error)
	Create(ctx context.Context, rsa *rbacv1.RbacServiceAccount) error
	Get(name string) (*rbacv1.RbacServiceAccount, error)
	Update(ctx context.Context, rsa *rbacv1.RbacServiceAccount) error
	Delete(ctx context.Context, name string) error
	Validate(username, password string) bool
}

type RbacServiceAccounts struct {
	clientSet kubernetes.Interface
	ctx       context.Context
	namespace string
	apps      App
	items     cmap.ConcurrentMap
}

func NewServiceAccount(ctx context.Context, cfg kubernetes.Interface, apps App) *RbacServiceAccounts {
	sa := &RbacServiceAccounts{
		clientSet: cfg,
		ctx:       ctx,
		namespace: apps.Namespace(),
		apps:      apps,
		items:     cmap.New(),
	}
	go sa.watch()
	return sa
}

func (sa *RbacServiceAccounts) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := sa.clientSet.RbacV1().RbacServiceAccounts(sa.namespace).Watch(sa.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-sa.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = sa.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (sa *RbacServiceAccounts) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Modified, watch.Added:
		obj := e.Object.(*rbacv1.RbacServiceAccount)
		sa.sync(obj)
	case watch.Deleted:
		obj := e.Object.(*rbacv1.RbacServiceAccount)
		sa.del(obj)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (sa *RbacServiceAccounts) sync(obj *rbacv1.RbacServiceAccount) {
	t, ok := sa.items.Get(obj.Name)
	if ok {
		ori := t.(*rbacv1.RbacServiceAccount)
		if ori.UID == obj.UID {
			if ori.ResourceVersion == obj.ResourceVersion {
				return
			}
		}
	}
	sa.items.Set(obj.Name, obj)
}

func (sa *RbacServiceAccounts) del(obj *rbacv1.RbacServiceAccount) {
	sa.items.Remove(obj.Name)
}

func (sa *RbacServiceAccounts) validateApps(apps []string) error {
	for _, v := range apps {
		_, err := sa.apps.Get(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sa *RbacServiceAccounts) List() (*rbacv1.RbacServiceAccountList, error) {
	res := &rbacv1.RbacServiceAccountList{
		ListMeta: metav1.ListMeta{},
		Items:    make([]rbacv1.RbacServiceAccount, 0),
	}
	all := sa.items.Items()
	for _, v := range all {
		tmp := v.(*rbacv1.RbacServiceAccount)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func (sa *RbacServiceAccounts) Create(ctx context.Context, rsa *rbacv1.RbacServiceAccount) error {
	if rsa.Name == "admin" {
		return fmt.Errorf("error couldn't use `admin` as Name")
	}
	rsa.Namespace = sa.namespace
	rsa.Spec.AccountMeta.Username = rsa.Name
	rsa.Spec.AccountMeta.Password = random.GenRandomString(32)
	if err := sa.validateApps(rsa.Spec.Apps); err != nil {
		return err
	}
	_, err := sa.clientSet.RbacV1().RbacServiceAccounts(sa.namespace).Create(ctx, rsa, metav1.CreateOptions{})
	return err
}

func (sa *RbacServiceAccounts) Get(name string) (*rbacv1.RbacServiceAccount, error) {
	t, ok := sa.items.Get(name)
	if !ok {
		return nil, fmt.Errorf("RbacServiceAccount name:%s not found", name)
	}
	return t.(*rbacv1.RbacServiceAccount), nil
}

func (sa *RbacServiceAccounts) Update(ctx context.Context, rsa *rbacv1.RbacServiceAccount) error {
	t, ok := sa.items.Get(rsa.Name)
	if !ok {
		return fmt.Errorf("RbacServiceAccount name:%s not found", rsa.Name)
	}
	ori := t.(*rbacv1.RbacServiceAccount).DeepCopy()
	rsa.Namespace = sa.namespace
	rsa.Spec.AccountMeta.Username = rsa.Name
	rsa.Spec.AccountMeta.Password = ori.Spec.AccountMeta.Password
	if err := sa.validateApps(rsa.Spec.Apps); err != nil {
		return err
	}
	_, err := sa.clientSet.RbacV1().RbacServiceAccounts(sa.namespace).Update(ctx, rsa, metav1.UpdateOptions{})
	return err
}

func (sa *RbacServiceAccounts) Delete(ctx context.Context, name string) error {
	if _, ok := sa.items.Get(name); !ok {
		return fmt.Errorf("RbacServiceAccount name:%s not found", name)
	}
	return sa.clientSet.RbacV1().RbacServiceAccounts(sa.namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (sa *RbacServiceAccounts) Validate(username, password string) bool {
	t, ok := sa.items.Get(username)
	if !ok {
		return false
	}
	p := t.(*rbacv1.RbacServiceAccount)
	return password == p.Spec.AccountMeta.Password
}
