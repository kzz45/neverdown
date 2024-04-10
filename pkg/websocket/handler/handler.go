package handler

import (
	"context"
	"sync"
)

type clientHandler func(in []byte, h Interface) (res []byte, err error)

type handler struct {
	sync.Once
	id             int64
	router         string
	managerHandler clientHandler
	removeChan     chan<- int64
	outputChan     chan<- []byte
	pingFunc       func()
	closeFunc      func()
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewHandler(ctx context.Context, router string, mh clientHandler, removeChan chan<- int64) *handler {
	sub, cancel := context.WithCancel(ctx)
	c := &handler{
		id:             0,
		router:         router,
		managerHandler: mh,
		removeChan:     removeChan,
		ctx:            sub,
		cancel:         cancel,
	}
	return c
}

func (h *handler) RegisterId(id int64) {
	h.id = id
}

func (h *handler) Id() int64 {
	return h.id
}

func (h *handler) RegisterRemoveChan(ch chan<- int64) {
	h.removeChan = ch
}

func (h *handler) RegisterConnWriteChan(ch chan<- []byte) {
	h.outputChan = ch
}

func (h *handler) RegisterConnClose(do func()) {
	h.closeFunc = do
}

func (h *handler) RegisterConnPing(do func()) {
	h.pingFunc = do
}

func (h *handler) Ping() {
	h.pingFunc()
}

func (h *handler) Handler(in []byte) (res []byte, err error) {
	return h.managerHandler(in, h)
}

func (h *handler) Run() {}

func (h *handler) Close() {
	h.Once.Do(func() {
		h.removeChan <- h.id
		h.cancel()
	})
}
