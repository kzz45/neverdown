package openx

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/hashicorp/go-version"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/registry"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"github.com/kzz45/neverdown/pkg/openx/controller"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type meta struct {
	key    string
	image  string
	policy openxv1.WatchPolicy

	tagName        string
	repositoryMeta jingxv1.RepositoryMeta
}

func (m *meta) validateIsMatch(ctx context.Context, jingx aggregator.Aggregator) (int, error) {
	t1 := strings.Split(m.image, "/")
	if len(t1) != 3 {
		return 1, fmt.Errorf("invalid meta image")
	}
	t2 := strings.Split(t1[2], ":")
	if len(t2) != 2 {
		return 2, fmt.Errorf("invalid meta image")
	}
	domain, project, repository, tag := t1[0], t1[1], t2[0], t2[1]
	p, err := jingx.Project.Get(ctx, project)
	if err != nil {
		return 3, err
	}
	// validate domain
	exist := false
	for _, v := range p.Spec.Domains {
		if v == domain {
			exist = true
			break
		}
	}
	if !exist {
		return 4, fmt.Errorf("project:%s domains don't match supported:%v target:%s", project, p.Spec.Domains, domain)
	}
	m.repositoryMeta = jingxv1.RepositoryMeta{
		ProjectName:    project,
		RepositoryName: repository,
	}
	_, err = jingx.Repository.Get(ctx, registry.GenRepositoryFullName(m.repositoryMeta))
	if err != nil {
		return 5, err
	}
	tagObj, err := jingx.Tag.Get(ctx, registry.GenerateFullTagName(m.repositoryMeta, tag))
	if err != nil {
		return 6, err
	}
	m.tagName = tagObj.Name
	return 0, nil
}

func (tc *OpenxController) watch() {
	res := tc.jingx.Tag.Watcher()
	for {
		select {
		//case <-tc.ctx.Done():
		//	return
		case e, isClosed := <-res:
			if !isClosed {
				return
			}
			tag := e.Object.(*jingxv1.Tag)
			zaplogger.Sugar().Infof("watch type:%s tag:%#v", e.Type, *tag)
			switch e.Type {
			case watch.Added:
				tc.rollingUpdate(tag)
			case watch.Modified:
				tc.inPlaceUpdate(tag)
			case watch.Deleted:
			default:
				zaplogger.Sugar().Infof("openx watchers unknown type: %s", e.Type)
			}
		}
	}
}

func (tc *OpenxController) inPlaceUpdate(tag *jingxv1.Tag) {
	// change replicas
	labels := map[string]string{
		controller.LabelJingxProject:    tag.Spec.RepositoryMeta.ProjectName,
		controller.LabelJingxRepository: tag.Spec.RepositoryMeta.RepositoryName,
		controller.LabelJingxTag:        tag.Spec.Tag,
		controller.LabelWatchPolicy:     string(openxv1.WatchPolicyInPlaceUpgrade),
	}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: labels,
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	dpList, err := tc.dpLister.List(selector)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	var zero int32 = 0
	for _, dp := range dpList {
		if dp.DeletionTimestamp != nil {
			continue
		}
		up := dp.DeepCopy()
		// todo is it necessary to use retry.OnError for conflict updating ?
		up.Spec.Replicas = &zero
		if _, err = tc.dpControl.Update(up); err != nil {
			zaplogger.Sugar().Error(err)
		}
	}
}

func validateTag(tag string) bool {
	reg, err := regexp.Compile(`^v?[0-9]+(.*?)$`)
	if err != nil {
		zaplogger.Sugar().Panic(err)
	}
	return reg.Match([]byte(tag))
}

func compareTags(tags []jingxv1.Tag) *jingxv1.Tag {
	versions := make([]*version.Version, 0)
	maps := make(map[string]jingxv1.Tag, 0)
	for _, raw := range tags {
		if ok := validateTag(raw.Spec.Tag); !ok {
			zaplogger.Sugar().Infof("tag:%s is invalid", raw.Spec.Tag)
			continue
		}
		v, err := version.NewVersion(raw.Spec.Tag)
		if err != nil {
			fmt.Println(err)
			continue
		}
		maps[v.String()] = raw
		versions = append(versions, v)
	}
	if len(versions) == 0 {
		return nil
	}
	sort.Sort(version.Collection(versions))
	tag := versions[len(versions)-1]
	a := maps[tag.String()]
	return a.DeepCopy()
}

type TagMeta struct {
	tag *jingxv1.Tag

	supportedImages map[string]bool
}

type projectHandler func(ctx context.Context, name string) (*jingxv1.Project, error)

func NewTagMeta(ctx context.Context, tag *jingxv1.Tag, handler projectHandler) (*TagMeta, error) {
	p, err := handler(ctx, tag.Spec.RepositoryMeta.ProjectName)
	if err != nil {
		return nil, err
	}
	tm := &TagMeta{
		tag:             tag,
		supportedImages: make(map[string]bool, 0),
	}
	for _, v := range p.Spec.Domains {
		tm.supportedImages[fmt.Sprintf("%s/%s/%s", v, tag.Spec.RepositoryMeta.ProjectName, tag.Spec.RepositoryMeta.RepositoryName)] = true
	}
	return tm, nil
}

func (tm *TagMeta) ValidateImage(image string) bool {
	t := strings.Split(image, ":")
	_, ok := tm.supportedImages[t[0]]
	return ok
}

func (tm *TagMeta) ValidateReplaceImage(image string) (string, bool) {
	t1 := strings.Split(image, ":")
	_, ok := tm.supportedImages[t1[0]]
	if !ok {
		return "", false
	}
	t2 := strings.Split(image, ":")
	return t2[0] + ":" + tm.tag.Spec.Tag, ok
}

func (tc *OpenxController) rollingUpdate(tag *jingxv1.Tag) {
	tags, err := tc.jingx.Tag.ListWithSelectors(context.TODO(), map[string]string{
		registry.JingxProject:    tag.Spec.RepositoryMeta.ProjectName,
		registry.JingxRepository: tag.Spec.RepositoryMeta.RepositoryName,
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	latest := compareTags(tags.Items)
	if latest == nil {
		return
	}
	tagMeta, err := NewTagMeta(context.TODO(), latest, func(ctx context.Context, name string) (*jingxv1.Project, error) {
		return tc.jingx.Project.Get(ctx, name)
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	// change image
	labels := map[string]string{
		controller.LabelJingxProject:    tag.Spec.RepositoryMeta.ProjectName,
		controller.LabelJingxRepository: tag.Spec.RepositoryMeta.RepositoryName,
		controller.LabelWatchPolicy:     string(openxv1.WatchPolicyRollingUpgrade),
	}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: labels,
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	dpList, err := tc.dpLister.List(selector)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	kinds := make(map[string]bool, 0)
	for _, dp := range dpList {
		if dp.DeletionTimestamp != nil {
			continue
		}
		if len(dp.Spec.Template.Spec.Containers) != 1 {
			continue
		}
		if ok := tagMeta.ValidateImage(dp.Spec.Template.Spec.Containers[0].Image); !ok {
			continue
		}
		controllerRef := metav1.GetControllerOf(dp)
		if controllerRef == nil {
			zaplogger.Sugar().Infof("Orphan Deployment:%s namespace:%s", dp.Name, dp.Namespace)
			continue
		}
		key := fmt.Sprintf("%s-%s", dp.Namespace, controllerRef.Kind)
		if _, ok := kinds[key]; ok {
			continue
		}
		// todo use retry.OnError for conflict update
		openx := tc.resolveControllerRef(dp.Namespace, controllerRef)
		if openx == nil {
			continue
		}
		kinds[key] = true
		if err := tc.updateAppImage(openx, tagMeta); err != nil {
			zaplogger.Sugar().Error(err)
		}
	}
}

func (tc *OpenxController) updateAppImage(openx *openxv1.Openx, tagMeta *TagMeta) error {
	apps := make([]openxv1.App, 0)
	changed := false
	for _, app := range openx.Spec.Applications {
		if app.WatchPolicy != openxv1.WatchPolicyRollingUpgrade {
			apps = append(apps, app)
			continue
		}
		if len(app.Pod.Spec.Containers) != 1 {
			apps = append(apps, app)
			zaplogger.Sugar().Infof("updateAppImage Containers != 1 namespace:%s name:%s app:%s", openx.Namespace, openx.Name, app.AppName)
			continue
		}
		image, ok := tagMeta.ValidateReplaceImage(app.Pod.Spec.Containers[0].Image)
		if !ok {
			apps = append(apps, app)
			zaplogger.Sugar().Infof("updateAppImage ValidateReplaceImage failed namespace:%s name:%s app:%s image:%s", openx.Namespace, openx.Name, app.AppName, app.Pod.Spec.Containers[0].Image)
			continue
		}
		apps = append(apps, app)
		zaplogger.Sugar().Infow("updateAppImage",
			zap.String("openx", openx.Name),
			zap.String("namespace", openx.Namespace),
			zap.String("app", app.AppName),
			zap.String("ori-image", app.Pod.Spec.Containers[0].Image),
			zap.String("update-image", image))
		app.Pod.Spec.Containers[0].Image = image
		changed = true
	}
	if changed {
		openx.Spec.Applications = apps
		if err := tc.forceSyncOpenx(openx); err != nil {
			zaplogger.Sugar().Error(err)
			return err
		}
	}
	return nil
}
