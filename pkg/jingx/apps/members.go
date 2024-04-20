package apps

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	cmap "github.com/orcaman/concurrent-map"

	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/apps/player"
	"github.com/kzz45/neverdown/pkg/jingx/proto"
	"github.com/kzz45/neverdown/pkg/jingx/resources"
)

func NewMembers(ctx context.Context, authorityClientSet kubernetes.Interface, api aggregator.Api) *Members {
	// hostname, err := env.GetHostName()
	// if err != nil {
	// zaplogger.Sugar().Fatal(err)
	// }
	m := &Members{
		incr:        0,
		removedChan: make(chan int64, 4096),
		players:     cmap.New(),
		resource:    resources.New(ctx, authorityClientSet, api),
		watchEvent:  make(chan *proto.Response, 4096),
		// hostname:    hostname,
		hostname: "hostname",
		ctx:      ctx,
	}
	go m.LoopRemoved()
	go m.Watch()
	go func() {
		_ = m.resource.Watch(m.watchEvent)
	}()
	return m
}

type Members struct {
	incr        int64
	removedChan chan int64
	players     cmap.ConcurrentMap
	resource    *resources.Resource
	watchEvent  chan *proto.Response
	hostname    string
	ctx         context.Context
}

func (m *Members) Resource() *resources.Resource {
	return m.resource
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
			zaplogger.Sugar().Info("remove player:", pid)
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
					bo, err := m.resource.ValidateAccess(p.Name(), msg.Namespace, msg.GroupVersionKind, string(msg.Verb))
					if err != nil {
						zaplogger.Sugar().Error(err)
						return
					}
					if !bo {
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
	zaplogger.Sugar().Info("ValidToken", zap.String("token", token))
	claims, err := jwttoken.Parse(token)
	if err != nil {
		zaplogger.Sugar().Error("ValidToken err",
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
	p, err := player.New(m.ctx, pid, claims, m.resource)
	if err != nil {
		return nil, err
	}
	m.players.Set(key, p)
	p.RegisterRemoveChan(m.removedChan)
	return p, nil
}
