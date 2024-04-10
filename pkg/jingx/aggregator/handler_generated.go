package aggregator

import (
	"context"
	"fmt"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ErrGVKNotExist = "unable to match GroupVersionKind:%v"

func (a *Aggregator) Create(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &jingxv1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Event.Create(ctx, obj)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		obj := &jingxv1.Project{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Project.Create(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		obj := &jingxv1.Repository{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Repository.Create(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		obj := &jingxv1.Tag{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Tag.Create(ctx, obj, author)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (a *Aggregator) Delete(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &jingxv1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Event.Delete(ctx, obj)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		obj := &jingxv1.Project{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Project.Delete(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		obj := &jingxv1.Repository{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Repository.Delete(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		obj := &jingxv1.Tag{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Tag.Delete(ctx, obj, author)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (a *Aggregator) Update(gvk *schema.GroupVersionKind, namespace string, raw []byte, author string) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj := &jingxv1.Event{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Event.Update(ctx, obj)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		obj := &jingxv1.Project{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Project.Update(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		obj := &jingxv1.Repository{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Repository.Update(ctx, obj, author)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		obj := &jingxv1.Tag{}
		if err = obj.Unmarshal(raw); err != nil {
			break
		}
		obj.Namespace = namespace
		err = a.Tag.Update(ctx, obj, author)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	return code, res, err
}

func (a *Aggregator) List(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	var obj Object
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		obj, err = a.Event.List()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		obj, err = a.Project.List()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		obj, err = a.Repository.List()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		obj, err = a.Tag.List()
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		code = 1
	}
	res, err = obj.Marshal()
	if err != nil {
		return 1, nil, err
	}
	return code, res, err
}

func (a *Aggregator) ConvertObjectToRaw(gvk *schema.GroupVersionKind, obj runtime.Object) (namespace string, res []byte, err error) {
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		into := obj.(*jingxv1.Event)
		namespace = into.Namespace
		res, err = into.Marshal()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		into := obj.(*jingxv1.Project)
		namespace = into.Namespace
		res, err = into.Marshal()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		into := obj.(*jingxv1.Repository)
		namespace = into.Namespace
		res, err = into.Marshal()
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		into := obj.(*jingxv1.Tag)
		namespace = into.Namespace
		res, err = into.Marshal()
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	return namespace, res, err
}

func (a *Aggregator) Providers() map[schema.GroupVersionKind]WatchProvider {
	providers := make(map[schema.GroupVersionKind]WatchProvider, 0)
	providers[schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}] = a.Event
	providers[schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}] = a.Project
	providers[schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}] = a.Repository
	providers[schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}] = a.Tag
	return providers
}
