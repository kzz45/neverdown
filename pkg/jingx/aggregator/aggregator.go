package aggregator

import (
	"context"

	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/jingx/registry"
)

type Aggregator struct {
	ctx        context.Context
	cfg        kubernetes.Interface
	Event      *registry.Event
	Project    *registry.Project
	Repository *registry.Repository
	Tag        *registry.Tag
}

func New(ctx context.Context, cfg kubernetes.Interface, namespace string) *Aggregator {
	ag := &Aggregator{
		ctx:        ctx,
		cfg:        cfg,
		Event:      registry.NewEvent(ctx, cfg, namespace),
		Project:    registry.NewProject(ctx, cfg, namespace),
		Repository: registry.NewRepository(ctx, cfg, namespace),
		Tag:        registry.NewTag(ctx, cfg, namespace),
	}
	ag.Project.AddRecordEventHandler(ag.Event.Record)
	ag.Repository.AddRecordEventHandler(ag.Event.Record)
	ag.Tag.AddRecordEventHandler(ag.Event.Record)
	return ag
}

func (a *Aggregator) ClientSet() kubernetes.Interface {
	return a.cfg
}
