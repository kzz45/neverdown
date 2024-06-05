package registry

import (
	"context"
	"fmt"
	"time"

	rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"

	"github.com/kzz45/discovery/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"

	jingxv1 "github.com/kzz45/discovery/pkg/apis/jingx/v1"
	"github.com/kzz45/discovery/pkg/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	DefaultNamespace = "discovery-jingx"
	LabelKey         = "discovery-jingx-appid"
)

type Project struct {
	clientSet    kubernetes.Interface
	ctx          context.Context
	namespace    string
	items        cmap.ConcurrentMap
	watchChannel chan watch.Event
	recorder     RecordEventHandler
}

func NewProject(ctx context.Context, cfg kubernetes.Interface, namespace string) *Project {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	p := &Project{
		clientSet: cfg,
		ctx:       ctx,
		namespace: namespace,
		items:     cmap.New(),
	}
	go p.watch()
	return p
}

func (p *Project) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := p.clientSet.JingxV1().Projects(p.namespace).Watch(p.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-p.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = p.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (p *Project) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Modified, watch.Added:
		obj := e.Object.(*jingxv1.Project)
		if ok := p.sync(obj); ok {
			p.triggerWatcher(e)
		}
	case watch.Deleted:
		obj := e.Object.(*jingxv1.Project)
		p.del(obj)
		p.triggerWatcher(e)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (p *Project) sync(obj *jingxv1.Project) bool {
	t, ok := p.items.Get(obj.Name)
	if ok {
		ori := t.(*jingxv1.Project)
		if ori.UID == obj.UID {
			if ori.ResourceVersion == obj.ResourceVersion {
				return false
			}
		}
	}
	p.items.Set(obj.Name, obj)
	return true
}

func (p *Project) del(obj *jingxv1.Project) {
	p.items.Remove(obj.Name)
}

func (p *Project) triggerWatcher(e watch.Event) {
	if p.watchChannel == nil {
		return
	}
	p.watchChannel <- e
}

func (p *Project) Watcher() <-chan watch.Event {
	if p.watchChannel != nil {
		return p.watchChannel
	}
	p.watchChannel = make(chan watch.Event, 10240)
	return p.watchChannel
}

func (p *Project) AddRecordEventHandler(recorder RecordEventHandler) {
	p.recorder = recorder
}

func (p *Project) recordEvent(ctx context.Context, author string, verb jingxv1.Verb, obj MarshalObject) error {
	if p.recorder == nil {
		return nil
	}
	return p.recorder(
		ctx,
		author,
		rbacv1.GroupVersionKind{
			Group:   jingxv1.GroupName,
			Version: jingxv1.SchemeGroupVersion.Version,
			Kind:    "Project",
		},
		verb,
		obj)
}

func (p *Project) Namespace() string {
	return p.namespace
}

func (p *Project) List() (*jingxv1.ProjectList, error) {
	res := &jingxv1.ProjectList{
		ListMeta: metav1.ListMeta{},
		Items:    make([]jingxv1.Project, 0),
	}
	all := p.items.Items()
	for _, v := range all {
		tmp := v.(*jingxv1.Project)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func (p *Project) Get(ctx context.Context, name string) (*jingxv1.Project, error) {
	t, ok := p.items.Get(name)
	if !ok {
		return nil, fmt.Errorf("project name:%s not found", name)
	}
	return t.(*jingxv1.Project), nil
}

func checkDuplicatedDomain(in []string) error {
	t := make(map[string]bool)
	for _, v := range in {
		if _, ok := t[v]; ok {
			return fmt.Errorf("dulicated domain:%s", v)
		}
		t[v] = true
	}
	return nil
}

func (p *Project) Create(ctx context.Context, in *jingxv1.Project, author string) error {
	if err := checkDuplicatedDomain(in.Spec.Domains); err != nil {
		return err
	}
	if err := ValidateName(in.Name); err != nil {
		return err
	}
	in.Namespace = p.namespace
	_, err := p.clientSet.JingxV1().Projects(p.namespace).Create(ctx, in, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return p.recordEvent(ctx, author, jingxv1.VerbCreate, in)
}

func (p *Project) checkBeforeDelete(ctx context.Context, projectName string) error {
	selector := map[string]string{
		JingxProject: projectName,
	}
	items, err := p.clientSet.JingxV1().Repositories(p.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: GetLabelSelector(selector),
	})
	if err != nil {
		return err
	}
	if len(items.Items) > 0 {
		return fmt.Errorf("error please delete linked jingx.v1.Repository before delete jingx.v1.Project, there was %d remains", len(items.Items))
	}
	return nil
}

func (p *Project) Delete(ctx context.Context, in *jingxv1.Project, author string) error {
	t, ok := p.items.Get(in.Name)
	if !ok {
		return fmt.Errorf("project name:%s not found", in.Name)
	}
	obj := t.(*jingxv1.Project)
	if err := p.checkBeforeDelete(ctx, in.Name); err != nil {
		return err
	}
	err := p.clientSet.JingxV1().Projects(p.namespace).Delete(ctx, in.Name, metav1.DeleteOptions{Preconditions: &metav1.Preconditions{
		UID:             &obj.UID,
		ResourceVersion: &obj.ResourceVersion,
	}})
	if err != nil {
		return err
	}
	return p.recordEvent(ctx, author, jingxv1.VerbDelete, obj)
}

func compareDomains(ori []string, up []string) bool {
	if len(ori) != len(up) {
		return true
	}
	for i := 0; i < len(ori); i++ {
		if ori[i] != up[i] {
			return true
		}
	}
	return false
}

func (p *Project) Update(ctx context.Context, in *jingxv1.Project, author string) error {
	// if err := checkDuplicatedDomain(in.Spec.Domains); err != nil {
	// 	return err
	// }
	t, ok := p.items.Get(in.Name)
	if !ok {
		return fmt.Errorf("error project name:%s not found", in.Name)
	}
	ori := t.(*jingxv1.Project).DeepCopy()

	if bo := compareDomains(ori.Spec.Domains, in.Spec.Domains); !bo {
		return fmt.Errorf("error project name:%s nothing changed", in.Name)
	}
	ori.Spec.Domains = in.Spec.Domains
	zaplogger.Sugar().Infof("project changed %v", ori)
	_, err := p.clientSet.JingxV1().Projects(p.namespace).Update(ctx, ori, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return p.recordEvent(ctx, author, jingxv1.VerbUpdate, ori)
}
