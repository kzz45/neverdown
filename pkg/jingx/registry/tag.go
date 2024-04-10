package registry

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/utils/random"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Tag struct {
	clientSet    kubernetes.Interface
	ctx          context.Context
	namespace    string
	items        cmap.ConcurrentMap
	tags         cmap.ConcurrentMap
	watchChannel chan watch.Event
	recorder     RecordEventHandler
	mu           sync.Mutex
}

func NewTag(ctx context.Context, cfg kubernetes.Interface, namespace string) *Tag {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	t := &Tag{
		clientSet: cfg,
		ctx:       ctx,
		namespace: namespace,
		items:     cmap.New(),
		tags:      cmap.New(),
	}
	go t.watch()
	return t
}

func (t *Tag) watch() {
	timeout := int64(3600) * 24
	var opts = metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
rewatch:
	res, err := t.clientSet.JingxV1().Tags(t.namespace).Watch(t.ctx, opts)
	if err != nil {
		zaplogger.Sugar().Error(err)
		time.Sleep(time.Second * 5)
		goto rewatch
	}
	for {
		select {
		case <-t.ctx.Done():
			res.Stop()
			return
		case e, isClose := <-res.ResultChan():
			if !isClose {
				res.Stop()
				goto rewatch
			}
			if err = t.handle(e); err != nil {
				zaplogger.Sugar().Error(err)
				res.Stop()
				goto rewatch
			}
		}
	}
}

func (t *Tag) handle(e watch.Event) error {
	var err error
	switch e.Type {
	case watch.Modified, watch.Added:
		obj := e.Object.(*jingxv1.Tag)
		if ok := t.sync(obj); ok {
			t.triggerWatcher(e)
		}
	case watch.Deleted:
		obj := e.Object.(*jingxv1.Tag)
		t.del(obj)
		t.triggerWatcher(e)
	case watch.Error:
		err = fmt.Errorf("watch receive ERROR event obj:%#v", e.Object)
	}
	return err
}

func (t *Tag) sync(obj *jingxv1.Tag) bool {
	o, ok := t.items.Get(obj.Name)
	if ok {
		ori := o.(*jingxv1.Tag)
		if ori.UID == obj.UID {
			if ori.ResourceVersion == obj.ResourceVersion {
				return false
			}
		}
	}
	t.items.Set(obj.Name, obj)
	t.tags.Set(GenerateFullTagName(obj.Spec.RepositoryMeta, obj.Spec.Tag), obj)
	return true
}

func (t *Tag) del(obj *jingxv1.Tag) {
	t.items.Remove(obj.Name)
	t.tags.Remove(GenerateFullTagName(obj.Spec.RepositoryMeta, obj.Spec.Tag))
}

func (t *Tag) triggerWatcher(e watch.Event) {
	if t.watchChannel == nil {
		return
	}
	t.watchChannel <- e
}

func (t *Tag) Watcher() <-chan watch.Event {
	t.watchChannel = make(chan watch.Event, 102400)
	return t.watchChannel
}

func (t *Tag) AddRecordEventHandler(recorder RecordEventHandler) {
	t.recorder = recorder
}

func (t *Tag) recordEvent(ctx context.Context, author string, verb jingxv1.Verb, obj MarshalObject) error {
	if t.recorder == nil {
		return nil
	}
	return t.recorder(
		ctx,
		author,
		rbacv1.GroupVersionKind{
			Group:   jingxv1.GroupName,
			Version: jingxv1.SchemeGroupVersion.Version,
			Kind:    "Tag",
		},
		verb,
		obj)
}

func (t *Tag) List() (*jingxv1.TagList, error) {
	res := &jingxv1.TagList{
		ListMeta: metav1.ListMeta{},
		Items:    make([]jingxv1.Tag, 0),
	}
	all := t.items.Items()
	for _, v := range all {
		tmp := v.(*jingxv1.Tag)
		res.Items = append(res.Items, *tmp)
	}
	return res, nil
}

func (t *Tag) ListWithSelectors(ctx context.Context, labels map[string]string) (*jingxv1.TagList, error) {
	return t.clientSet.JingxV1().Tags(t.namespace).List(ctx, metav1.ListOptions{LabelSelector: GetLabelSelector(labels)})
}

func GenerateFullTagName(in jingxv1.RepositoryMeta, tag string) string {
	return fmt.Sprintf("%s-%s-%s", in.ProjectName, in.RepositoryName, tag)
}

func (t *Tag) Get(ctx context.Context, name string) (*jingxv1.Tag, error) {
	obj, ok := t.tags.Get(name)
	if !ok {
		return nil, fmt.Errorf("tag full-name:%s not found", name)
	}
	return obj.(*jingxv1.Tag), nil
}

func (t *Tag) genTagUid() string {
retry:
	uid := "tag" + "-" + strings.ToLower(random.GenRandomString(6))
	if _, ok := t.items.Get(uid); ok {
		goto retry
	}
	return uid
}

func (t *Tag) validateRepositoryMeta(ctx context.Context, meta jingxv1.RepositoryMeta) error {
	name := GenRepositoryFullName(meta)
	_, err := t.clientSet.JingxV1().Repositories(t.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ValidateTag(tag string) error {
	if tag == "latest" {
		return nil
	}
	reg, err := regexp.Compile(`^v?[0-9]+(.*?)$`)
	if err != nil {
		return err
	}
	if bo := reg.Match([]byte(tag)); !bo {
		return fmt.Errorf("invalid tag: %s doesn't match rule", tag)
	}
	return nil
}

func (t *Tag) Create(ctx context.Context, in *jingxv1.Tag, author string) error {
	if err := t.validateRepositoryMeta(ctx, in.Spec.RepositoryMeta); err != nil {
		return err
	}
	if err := ValidateTag(in.Spec.Tag); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tagName := GenerateFullTagName(in.Spec.RepositoryMeta, in.Spec.Tag)
	if _, ok := t.tags.Get(tagName); ok {
		return fmt.Errorf("repeated tag name:%s", in.Spec.Tag)
	}
	t.tags.Set(tagName, &jingxv1.Tag{})
	in.Labels = map[string]string{
		JingxProject:    in.Spec.RepositoryMeta.ProjectName,
		JingxRepository: in.Spec.RepositoryMeta.RepositoryName,
	}
	in.Namespace = t.namespace
	in.Name = t.genTagUid()
	in.Spec.DockerImage.LastModifiedTime = time.Now().Unix()
	if _, err := t.clientSet.JingxV1().Tags(t.namespace).Create(ctx, in, metav1.CreateOptions{}); err != nil {
		t.tags.Remove(tagName)
		return err
	}
	return t.recordEvent(ctx, author, jingxv1.VerbCreate, in)
}

func (t *Tag) Delete(ctx context.Context, in *jingxv1.Tag, author string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	o, ok := t.items.Get(in.Name)
	if !ok {
		return fmt.Errorf("tag name:%s not found", in.Name)
	}
	obj := o.(*jingxv1.Tag)
	err := t.clientSet.JingxV1().Tags(t.namespace).Delete(ctx, in.Name, metav1.DeleteOptions{Preconditions: &metav1.Preconditions{
		UID:             &obj.UID,
		ResourceVersion: &obj.ResourceVersion,
	}})
	if err != nil {
		return err
	}
	return t.recordEvent(ctx, author, jingxv1.VerbDelete, obj)
}

func compareRepositoryMeta(ori, in jingxv1.RepositoryMeta) error {
	if ori.RepositoryName == in.RepositoryName && ori.ProjectName == in.ProjectName {
		return nil
	}
	return fmt.Errorf("tag update couldn't modify RepositoryMeta")
}

func (t *Tag) Update(ctx context.Context, in *jingxv1.Tag, author string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	o, ok := t.items.Get(in.Name)
	if !ok {
		return fmt.Errorf("tag name:%s not found", in.Name)
	}
	obj := o.(*jingxv1.Tag).DeepCopy()
	if obj.Spec.Tag != in.Spec.Tag {
		return fmt.Errorf("tag update couldn't modify Tag")
	}
	if obj.Spec.DockerImage.Sha256 == in.Spec.DockerImage.Sha256 {
		zaplogger.Sugar().Errorf("tag update error DockerImage.Sha256 doesn't changed ori:%#v up:%#v", *obj, *in)
		return fmt.Errorf("tag update error DockerImage.Sha256 doesn't changed")
	}
	if err := compareRepositoryMeta(obj.Spec.RepositoryMeta, in.Spec.RepositoryMeta); err != nil {
		return err
	}
	obj.Spec.GitReference = in.Spec.GitReference
	obj.Spec.DockerImage = in.Spec.DockerImage
	obj.Spec.DockerImage.LastModifiedTime = time.Now().Unix()
	_, err := t.clientSet.JingxV1().Tags(t.namespace).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return t.recordEvent(ctx, author, jingxv1.VerbDelete, obj)
}
