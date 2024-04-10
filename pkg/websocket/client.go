package websocket

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	ConnectionTimeoutInSec      = 10
	WriteToChannelTimeoutInMS   = 1000
	RetryConnectionDurationInMS = 3000
)

func NewOption(ctx context.Context, url url.URL, pingData []byte, dialer *websocket.Dialer) *Option {
	sub, cancel := context.WithCancel(ctx)
	return &Option{
		Dialer:            dialer,
		Url:               url,
		Status:            OptionInActive,
		pingData:          pingData,
		Maintain:          OptionMaintainRetry,
		MaxRetryCount:     -1,
		RetryDuration:     RetryConnectionDurationInMS,
		registerFunctions: make([]OptionRegisterFunction, 0),
		readHandlerChan:   make(chan []byte, 4096),
		writeHandlerChan:  make(chan []byte, 4096),
		ctx:               sub,
		cancelFunc:        cancel,
		writeTimeout:      WriteToChannelTimeoutInMS,
	}
}

type OptionStatus string

const (
	OptionActive   OptionStatus = "active"
	OptionInActive OptionStatus = "inactive"
	OptionClosed   OptionStatus = "closed"
)

type OptionMaintainType string

const (
	OptionMaintainOnce  OptionMaintainType = "once"
	OptionMaintainRetry OptionMaintainType = "retry"
)

type OptionRegisterFunction func() error

type Option struct {
	rw   sync.RWMutex
	once sync.Once

	Dialer *websocket.Dialer
	Url    url.URL
	Status OptionStatus

	Maintain      OptionMaintainType
	MaxRetryCount int64
	// RetryDuration was the wait time in Millisecond
	RetryDuration int64

	pingData          []byte
	registerFunctions []OptionRegisterFunction
	readHandlerChan   chan []byte
	writeHandlerChan  chan []byte
	ctx               context.Context
	cancelFunc        context.CancelFunc

	writeTimeout int64
}

func (o *Option) SetMaxRetryCount(in int64) {
	if in <= 0 {
		o.MaxRetryCount = -1
	} else {
		o.MaxRetryCount = in
	}
}

func (o *Option) Next() bool {
	if o.Maintain == OptionMaintainOnce {
		return false
	}
	switch o.MaxRetryCount {
	case -1:
		return true
	case 0:
		return false
	default:
		o.MaxRetryCount -= 1
		return true
	}
}

func (o *Option) RegisterFunc(do ...OptionRegisterFunction) {
	for _, v := range do {
		o.registerFunctions = append(o.registerFunctions, v)
	}
}

func (o *Option) Prepare() error {
	for _, v := range o.registerFunctions {
		if err := v(); err != nil {
			zaplogger.Sugar().Errorw("websocket Option Prepare do func failed",
				zap.Error(err))
			return err
		}
	}
	return nil
}

func (o *Option) Read() <-chan []byte {
	return o.readHandlerChan
}

func (o *Option) Send(in []byte) error {
	o.rw.RLock()
	defer o.rw.RUnlock()
	switch o.Status {
	case OptionActive:
		select {
		case <-time.After(time.Millisecond * time.Duration(o.writeTimeout)):
			return fmt.Errorf("option send message failed due to write to channel timeout:%v ms", o.writeTimeout)
		case o.writeHandlerChan <- in:
		}
	case OptionInActive, OptionClosed:
		return fmt.Errorf("option skip sending message due to the status was:%v", o.Status)
	}
	return nil
}

func (o *Option) ChangeStatus(s OptionStatus) {
	o.rw.Lock()
	defer o.rw.Unlock()
	switch o.Status {
	case OptionClosed:
	case OptionActive, OptionInActive:
		o.Status = s
	}
}

func (o *Option) Done() <-chan struct{} {
	return o.ctx.Done()
}

func (o *Option) Cancel() {
	o.once.Do(func() {
		o.ChangeStatus(OptionClosed)
		o.cancelFunc()
	})
}

func NewClient(opt *Option) (*Client, error) {
	dialer := websocket.DefaultDialer
	if opt.Dialer != nil {
		dialer = opt.Dialer
	}
	a, _, err := dialer.Dial(opt.Url.String(), nil)
	if err != nil {
		zaplogger.Sugar().Errorw("websocket.DefaultDialer.Dial err",
			zap.String("path", opt.Url.String()),
			zap.Error(err))
		return nil, err
	}
	sub, cancel := context.WithCancel(opt.ctx)
	c := &Client{
		opt:        opt,
		conn:       a,
		Context:    sub,
		CancelFunc: cancel,
	}
	go c.readPump()
	go c.writePump()
	go c.ping(opt.pingData)
	return c, nil
}

type Client struct {
	sync.Once
	opt        *Option
	conn       *websocket.Conn
	Context    context.Context
	CancelFunc context.CancelFunc
}

func (c *Client) Option() *Option {
	return c.opt
}

func (c *Client) Close() {
	c.Once.Do(func() {
		c.CancelFunc()
	})
}

func (c *Client) ping(in []byte) {
	defer c.Close()
	tick := time.NewTicker(time.Second * time.Duration(ConnectionTimeoutInSec/2))
	defer tick.Stop()
	for {
		select {
		case <-c.Context.Done():
			return
		case <-tick.C:
			c.opt.writeHandlerChan <- in
		}
	}
}

func (c *Client) readPump() {
	defer c.Close()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			zaplogger.Sugar().Debug(err)
			return
		}
		c.opt.readHandlerChan <- message
	}
}

func (c *Client) writePump() {
	defer c.Close()
	for {
		select {
		case msg, isClose := <-c.opt.writeHandlerChan:
			if !isClose {
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				zaplogger.Sugar().Error(err)
				return
			}
		case <-c.Context.Done():
			return
		}
	}
}

func NewClientWithReconnect(opt *Option, stop chan<- struct{}) {
	defer func() {
		opt.Cancel()
		if stop != nil {
			stop <- struct{}{}
		}
	}()
	for {
		switch opt.Status {
		case OptionClosed:
			return
		case OptionActive, OptionInActive:
			client, err := NewClient(opt)
			if err != nil {
				if ok := opt.Next(); !ok {
					return
				}
				retryWaitInMillisecond(opt.RetryDuration)
				continue
			}
			opt.ChangeStatus(OptionActive)
			if err := opt.Prepare(); err != nil {
				return
			}
			select {
			case <-opt.Done():
				return
			case <-client.Context.Done():
				if ok := opt.Next(); !ok {
					return
				}
				retryWaitInMillisecond(opt.RetryDuration)
			}
		}
	}
}

func retryWaitInMillisecond(sleep int64) {
	if sleep == 0 {
		return
	}
	time.Sleep(time.Millisecond * time.Duration(sleep))
}
