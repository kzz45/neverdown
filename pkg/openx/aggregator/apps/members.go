package apps

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/websocket/env"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/apps/player"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/proto"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/resources"
)

func NewMembers(ctx context.Context, resources *resources.Resources) *Members {
	hostname, err := env.GetHostName()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	m := &Members{
		incr:        0,
		removedChan: make(chan int64, 4096),
		players:     cmap.New(),
		resources:   resources,
		watchEvent:  resources.AddEventHandler(),
		hostname:    hostname,
		ctx:         ctx,
	}
	go m.LoopRemoved()
	go m.Watch()
	return m
}

type Members struct {
	incr        int64
	removedChan chan int64
	players     cmap.ConcurrentMap
	resources   *resources.Resources
	watchEvent  <-chan *proto.Response
	hostname    string
	ctx         context.Context
}

func (m *Members) Resource() *resources.Resources {
	return m.resources
}

func (m *Members) LoopRemoved() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case pid, isClose := <-m.removedChan:
			if !isClose {
				return
			}
			key := fmt.Sprintf("%d", pid)
			m.players.Remove(key)
			m.resources.MetricsCollector().DecOnlinePlayer()
		}
	}
}

func (m *Members) Watch() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case msg, isClose := <-m.watchEvent:
			if !isClose {
				return
			}
			data, err := msg.Marshal()
			if err != nil {
				zaplogger.Sugar().Error(err)
				continue
			}
			items := m.players.Items()
			var wg sync.WaitGroup
			wg.Add(len(items))
			for _, v := range items {
				go func(v interface{}) {
					defer wg.Done()
					p := v.(*player.Player)
					bo, err := m.resources.ValidateAccess(p.Name(), msg.Namespace, msg.GroupVersionKind, string(msg.Verb))
					if err != nil {
						zaplogger.Sugar().Error(err)
						return
					}
					if !bo {
						zaplogger.Sugar().Debugw("watch noRoot",
							zap.Any("namespace:", msg.Namespace),
							zap.Any("gvk", msg.GroupVersionKind),
							zap.Any("verb", msg.Verb))
						return
					}
					if err := p.TransferPushNotify(data); err != nil {
						zaplogger.Sugar().Error(err)
					}
				}(v)
			}
			wg.Wait()
		}
	}
}

// ValidToken
func (m *Members) ValidToken(token string) (*player.Player, bool) {
	claims, err := jwttoken.Parse(token)
	if err != nil {
		zaplogger.Sugar().Debugw("Members ValidToken err",
			zap.String("token", token),
			zap.Error(err))
		return nil, false
	}
	p, err := m.loadPlayer(claims)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, false
	}
	return p, true
}

// loadPlayer
func (m *Members) loadPlayer(claims *jwttoken.Claims) (*player.Player, error) {
	pid := atomic.AddInt64(&m.incr, 1)
	key := fmt.Sprintf("%d", pid)
	p, err := player.New(m.ctx, pid, claims, m.resources)
	if err != nil {
		return nil, err
	}
	m.players.Set(key, p)
	p.RegisterRemoveChan(m.removedChan)
	m.resources.MetricsCollector().IncrOnlinePlayer()
	return p, nil
}
