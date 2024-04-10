package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	"github.com/kzz45/neverdown/pkg/authx/rbac/app"
	"github.com/kzz45/neverdown/pkg/utils/random"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type App interface {
	Namespace() string
	List() (*rbacv1.AppList, error)
	Create(ctx context.Context, app *rbacv1.App) error
	Get(name string) (*rbacv1.App, error)
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, app *rbacv1.App) error
	Validate(appId, appSecret string) error
	GenericApp(appId string) (*app.GenericApp, error)
}

const (
	DefaultNamespace = "discovery-authx"
	LabelKey         = "discovery-authx-appid"
)

type Apps struct {
	clientSet kubernetes.Interface
	ctx       context.Context
	namespace string
	items     cmap.ConcurrentMap
	accounts  cmap.ConcurrentMap
}

func NewApps(ctx context.Context, cfg kubernetes.Interface, namespace string) *Apps {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	a := &Apps{
		clientSet: cfg,
		ctx:       ctx,
		namespace: namespace,
		items:     cmap.New(),
		accounts:  cmap.New(),
	}
	go a.watch()
	return a
}

func (a *Apps) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := a.clientSet.RbacV1().Apps(a.namespace).Watch(a.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-a.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = a.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (a *Apps) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Modified, watch.Added:
		obj := e.Object.(*rbacv1.App)
		a.sync(obj)
	case watch.Deleted:
		obj := e.Object.(*rbacv1.App)
		a.del(obj)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (a *Apps) sync(obj *rbacv1.App) {
	t, ok := a.items.Get(obj.Name)
	if ok {
		ori := t.(*rbacv1.App)
		if ori.UID == obj.UID {
			if ori.ResourceVersion == obj.ResourceVersion {
				return
			}
		}
	}
	a.items.Set(obj.Name, obj)
	lowerAppid := strings.ToLower(obj.Spec.Id)
	if t2, ok := a.accounts.Get(lowerAppid); ok {
		ga := t2.(*app.GenericApp)
		ga.Sync(obj)
	} else {
		a.accounts.Set(lowerAppid, app.NewGenericApp(a.ctx, a.clientSet, obj))
	}
}

func (a *Apps) del(obj *rbacv1.App) {
	a.items.Remove(obj.Name)
	a.accounts.Remove(strings.ToLower(obj.Spec.Id))
}

func (a *Apps) Namespace() string {
	return a.namespace
}

func (a *Apps) List() (*rbacv1.AppList, error) {
	res := &rbacv1.AppList{
		ListMeta: metav1.ListMeta{},
		Items:    make([]rbacv1.App, 0),
	}
	all := a.items.Items()
	for _, v := range all {
		tmp := v.(*rbacv1.App)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func (a *Apps) Create(ctx context.Context, app *rbacv1.App) error {
	app.Namespace = a.namespace
	for {
		app.Spec.Id = random.GenRandomString(32)
		if _, ok := a.items.Get(strings.ToLower(app.Spec.Id)); ok {
			continue
		}
		break
	}
	app.Spec.Secret = random.GenRandomString(32)
	app.Labels = map[string]string{LabelKey: app.Spec.Id}
	_, err := a.clientSet.RbacV1().Apps(a.namespace).Create(ctx, app, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func CreateRbacV1App(ctx context.Context, clientSet kubernetes.Interface, app *rbacv1.App) (*rbacv1.App, error) {
	app.Namespace = DefaultNamespace
	app.Spec.Id = random.GenRandomString(32)
	app.Spec.Secret = random.GenRandomString(32)
	app.Labels = map[string]string{LabelKey: app.Spec.Id}
	return clientSet.RbacV1().Apps(app.Namespace).Create(ctx, app, metav1.CreateOptions{})
}

func GetRbacV1App(ctx context.Context, clientSet kubernetes.Interface, name string) (*rbacv1.App, error) {
	return clientSet.RbacV1().Apps(DefaultNamespace).Get(ctx, name, metav1.GetOptions{})
}

func (a *Apps) Get(name string) (*rbacv1.App, error) {
	t, ok := a.items.Get(name)
	if !ok {
		return nil, fmt.Errorf("app name:%s not found", name)
	}
	return t.(*rbacv1.App), nil
}

func (a *Apps) checkBeforeDelete(appId string) error {
	t, ok := a.accounts.Get(strings.ToLower(appId))
	if !ok {
		return fmt.Errorf("checkBeforeDelete error appid:%s not found", appId)
	}
	ga := t.(*app.GenericApp)
	if res, err := ga.Kind().List(); err != nil {
		return err
	} else {
		if len(res.Items) > 0 {
			return fmt.Errorf("error please delete app.Kind before delete app, there was %d remains", len(res.Items))
		}
	}
	if res, err := ga.Role().List(); err != nil {
		return err
	} else {
		if len(res.Items) > 0 {
			return fmt.Errorf("error please delete app.Role before delete app, there was %d remains", len(res.Items))
		}
	}
	if res, err := ga.ServiceAccount().List(); err != nil {
		return err
	} else {
		if len(res.Items) > 0 {
			return fmt.Errorf("error please delete app.ServiceAccount before delete app, there was %d remains", len(res.Items))
		}
	}
	return nil
}

func (a *Apps) Delete(ctx context.Context, name string) error {
	t, ok := a.items.Get(name)
	if !ok {
		return fmt.Errorf("app name:%s not found", name)
	}
	obj := t.(*rbacv1.App)
	if err := a.checkBeforeDelete(obj.Spec.Id); err != nil {
		return err
	}
	return a.clientSet.RbacV1().Apps(a.namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{
		Preconditions: &metav1.Preconditions{
			UID:             &obj.UID,
			ResourceVersion: &obj.ResourceVersion,
		},
	})
}

func (a *Apps) Update(ctx context.Context, app *rbacv1.App) error {
	t, ok := a.items.Get(app.Name)
	if !ok {
		return fmt.Errorf("error appname:%s not found", app.Name)
	}
	ori := t.(*rbacv1.App).DeepCopy()
	ori.Spec.Desc = app.Spec.Desc
	_, err := a.clientSet.RbacV1().Apps(a.namespace).Update(ctx, ori, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *Apps) Validate(appId, appSecret string) error {
	t, ok := a.accounts.Get(strings.ToLower(appId))
	if !ok {
		return fmt.Errorf("error appid:%s not found", appId)
	}
	ga := t.(*app.GenericApp)
	if ga.RbacV1App().Spec.Secret != appSecret {
		return fmt.Errorf("appid:%s error appSecret", appId)
	}
	return nil
}

func (a *Apps) GenericApp(appId string) (*app.GenericApp, error) {
	t, ok := a.accounts.Get(strings.ToLower(appId))
	if !ok {
		return nil, fmt.Errorf("error appid:%s not found", appId)
	}
	return t.(*app.GenericApp), nil
}
