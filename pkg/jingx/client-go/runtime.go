package clientgo

import (
	"fmt"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	"github.com/kzz45/neverdown/pkg/jingx/proto"

	cmap "github.com/orcaman/concurrent-map"
)

type runtime struct {
	repositories cmap.ConcurrentMap
	tags         cmap.ConcurrentMap
}

func newRuntime() *runtime {
	return &runtime{
		repositories: cmap.New(),
		tags:         cmap.New(),
	}
}

func generateFullTagName(in *jingxv1.Tag) string {
	return fmt.Sprintf("%s-%s-%s", in.Spec.RepositoryMeta.ProjectName, in.Spec.RepositoryMeta.RepositoryName, in.Spec.Tag)
}

func (r *runtime) repository(eventType proto.EventType, in *jingxv1.Repository) {
	name := in.Name
	switch eventType {
	case proto.EventAdded, proto.EventModified:
		t, ok := r.repositories.Get(name)
		if !ok {
			r.repositories.Set(name, in)
			return
		}
		ori := t.(*jingxv1.Tag)
		if ori.UID == in.UID {
			if ori.ResourceVersion >= in.ResourceVersion {
				return
			}
		}
		r.repositories.Set(name, in)
	case proto.EventDeleted:
		r.repositories.Remove(name)
	}
}

func (r *runtime) checkRepositoryExist(in *jingxv1.Repository) bool {
	_, ok := r.repositories.Get(in.Name)
	return ok
}

func (r *runtime) tag(eventType proto.EventType, in *jingxv1.Tag) {
	name := generateFullTagName(in)
	switch eventType {
	case proto.EventAdded, proto.EventModified:
		t, ok := r.tags.Get(name)
		if !ok {
			r.tags.Set(name, in)
			return
		}
		ori := t.(*jingxv1.Tag)
		if ori.UID == in.UID {
			if ori.ResourceVersion >= in.ResourceVersion {
				return
			}
		}
		r.tags.Set(name, in)
	case proto.EventDeleted:
		r.tags.Remove(name)
	}
}

func (r *runtime) checkTagExist(in *jingxv1.Tag) (*jingxv1.Tag, bool) {
	name := generateFullTagName(in)
	t, ok := r.tags.Get(name)
	if !ok {
		return nil, false
	}
	return t.(*jingxv1.Tag), true
}
