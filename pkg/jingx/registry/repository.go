package registry

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Repository struct {
	clientSet    kubernetes.Interface
	ctx          context.Context
	namespace    string
	items        cmap.ConcurrentMap
	watchChannel chan watch.Event
	recorder     RecordEventHandler
}

func NewRepository(ctx context.Context, cfg kubernetes.Interface, namespace string) *Repository {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	r := &Repository{
		clientSet: cfg,
		ctx:       ctx,
		namespace: namespace,
		items:     cmap.New(),
	}
	go r.watch()
	return r
}

func (r *Repository) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := r.clientSet.JingxV1().Repositories(r.namespace).Watch(r.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-r.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = r.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (r *Repository) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Modified, watch.Added:
		obj := e.Object.(*jingxv1.Repository)
		if ok := r.sync(obj); ok {
			r.triggerWatcher(e)
		}
	case watch.Deleted:
		obj := e.Object.(*jingxv1.Repository)
		r.del(obj)
		r.triggerWatcher(e)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (r *Repository) sync(obj *jingxv1.Repository) bool {
	t, ok := r.items.Get(obj.Name)
	if ok {
		ori := t.(*jingxv1.Repository)
		if ori.UID == obj.UID {
			if ori.ResourceVersion == obj.ResourceVersion {
				return false
			}
		}
	}
	r.items.Set(obj.Name, obj)
	return true
}

func (r *Repository) del(obj *jingxv1.Repository) {
	r.items.Remove(obj.Name)
}

func (r *Repository) triggerWatcher(e watch.Event) {
	if r.watchChannel == nil {
		return
	}
	r.watchChannel <- e
}

func (r *Repository) Watcher() <-chan watch.Event {
	if r.watchChannel != nil {
		return r.watchChannel
	}
	r.watchChannel = make(chan watch.Event, 10240)
	return r.watchChannel
}

func (r *Repository) AddRecordEventHandler(recorder RecordEventHandler) {
	r.recorder = recorder
}

func (r *Repository) recordEvent(ctx context.Context, author string, verb jingxv1.Verb, obj MarshalObject) error {
	if r.recorder == nil {
		return nil
	}
	return r.recorder(
		ctx,
		author,
		rbacv1.GroupVersionKind{
			Group:   jingxv1.GroupName,
			Version: jingxv1.SchemeGroupVersion.Version,
			Kind:    "Repository",
		},
		verb,
		obj)
}

func (r *Repository) List() (*jingxv1.RepositoryList, error) {
	res := &jingxv1.RepositoryList{
		ListMeta: metav1.ListMeta{},
		Items:    make([]jingxv1.Repository, 0),
	}
	all := r.items.Items()
	for _, v := range all {
		tmp := v.(*jingxv1.Repository)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func (r *Repository) Get(ctx context.Context, name string) (*jingxv1.Repository, error) {
	t, ok := r.items.Get(name)
	if !ok {
		return nil, fmt.Errorf("repository name:%s not found", name)
	}
	return t.(*jingxv1.Repository), nil
}

func GenRepositoryFullName(meta jingxv1.RepositoryMeta) string {
	return meta.ProjectName + "-" + meta.RepositoryName
}

func (r *Repository) validateProjectName(ctx context.Context, projectName string) error {
	_, err := r.clientSet.JingxV1().Projects(r.namespace).Get(ctx, projectName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, in *jingxv1.Repository, author string) error {
	if err := r.validateProjectName(ctx, in.Spec.RepositoryMeta.ProjectName); err != nil {
		return err
	}
	if err := ValidateName(in.Name); err != nil {
		return err
	}
	in.Labels = map[string]string{
		JingxProject: in.Spec.RepositoryMeta.ProjectName,
	}
	in.Namespace = r.namespace
	in.Name = GenRepositoryFullName(in.Spec.RepositoryMeta)
	_, err := r.clientSet.JingxV1().Repositories(r.namespace).Create(ctx, in, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return r.recordEvent(ctx, author, jingxv1.VerbCreate, in)
}

func (r *Repository) checkBeforeDelete(ctx context.Context, meta jingxv1.RepositoryMeta) error {
	selector := map[string]string{
		JingxProject:    meta.ProjectName,
		JingxRepository: meta.RepositoryName,
	}
	items, err := r.clientSet.JingxV1().Tags(r.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: GetLabelSelector(selector),
	})
	if err != nil {
		return err
	}
	if len(items.Items) > 0 {
		return fmt.Errorf("error please delete linked jingx.v1.Tag before delete jingx.v1.Repository, there was %d remains", len(items.Items))
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, in *jingxv1.Repository, author string) error {
	t, ok := r.items.Get(in.Name)
	if !ok {
		return fmt.Errorf("repository name:%s not found", in.Name)
	}
	obj := t.(*jingxv1.Repository)
	if err := r.checkBeforeDelete(ctx, in.Spec.RepositoryMeta); err != nil {
		return err
	}
	err := r.clientSet.JingxV1().Repositories(r.namespace).Delete(ctx, in.Name, metav1.DeleteOptions{Preconditions: &metav1.Preconditions{
		UID:             &obj.UID,
		ResourceVersion: &obj.ResourceVersion,
	}})
	if err != nil {
		return err
	}
	return r.recordEvent(ctx, author, jingxv1.VerbDelete, obj)
}

func (r *Repository) Update(ctx context.Context, in *jingxv1.Repository, author string) error {
	return fmt.Errorf("repository couldn't been modified")
}
