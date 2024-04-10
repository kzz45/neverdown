package player

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/kzz45/neverdown/pkg/jingx/proto"
	"github.com/kzz45/neverdown/pkg/jingx/resources"
)

func New(ctx context.Context, pid int64, claims *jwttoken.Claims, resource *resources.Resource) (*Player, error) {
	sub, cancel := context.WithCancel(ctx)
	p := &Player{
		id:       pid,
		ctx:      sub,
		cancel:   cancel,
		claims:   claims,
		resource: resource,
	}
	return p, nil
}

type Player struct {
	id         int64
	clearChan  chan<- int64
	outputChan chan<- []byte
	pingFunc   func()
	closeFunc  func()
	once       sync.Once
	ctx        context.Context
	cancel     context.CancelFunc

	claims   *jwttoken.Claims
	resource *resources.Resource
}

func (p *Player) Name() string {
	return p.claims.Username
}

func (p *Player) Id() int64 {
	return p.id
}

func (p *Player) RegisterId(id int64) {
	p.id = id
}

func (p *Player) RegisterRemoveChan(ch chan<- int64) {
	p.clearChan = ch
}

func (p *Player) RegisterConnWriteChan(ch chan<- []byte) {
	p.outputChan = ch
}

func (p *Player) RegisterConnClose(do func()) {
	p.closeFunc = do
}

func (p *Player) RegisterConnPing(do func()) {
	p.pingFunc = do
}

func (p *Player) Ping() {
	p.pingFunc()
}

func (p *Player) Close() {
	p.once.Do(func() {
		if p.closeFunc != nil {
			go p.closeFunc()
		}
		if p.cancel != nil {
			p.cancel()
		}
		if p.clearChan != nil {
			p.clearChan <- p.Id()
		}
	})
}

func (p *Player) Run() {}

func (p *Player) Handler(in []byte) (res []byte, err error) {
	req := &proto.Request{}
	if err = req.Unmarshal(in); err != nil {
		return res, err
	}
	var code int32
	switch req.Verb {
	case proto.VerbPing:
		p.Ping()
		res = []byte("ping success")
		if time.Now().Unix() >= p.claims.ExpiresAt {
			return res, fmt.Errorf("token was expired")
		}
	default:
		code, res, err = p.resource.Handler(p.claims.Username, req)
	}
	response := &proto.Response{
		Code:             code,
		GroupVersionKind: req.GroupVersionKind,
		Namespace:        req.Namespace,
		Verb:             req.Verb,
		Raw:              res,
	}
	if err != nil {
		zaplogger.Sugar().Error(err)
		response.Raw = []byte(err.Error())
	}
	return response.Marshal()
}

// TransferPushNotify
func (p *Player) TransferPushNotify(in []byte) error {
	zaplogger.Sugar().Debugf("player:%d TransferPushNotify message:%s", p.Id(), string(in))
	p.outputChan <- in
	return nil
}
