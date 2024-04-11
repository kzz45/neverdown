package pods

import (
	"context"
	"net/http"

	"github.com/kzz45/neverdown/pkg/execpod"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/resources"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Handler struct {
	resources  *resources.Resources
	ClientSet  kubernetes.Interface
	RestConfig *rest.Config
	sessionHub execpod.SessionHub
}

func New(r *resources.Resources) *Handler {
	clientset := r.ClientBuilder().ClientOrDie("openx-exec-pod")
	cfg := r.ClientBuilder().ConfigOrDie("openx-exec-pod")
	h := &Handler{
		resources:  r,
		ClientSet:  clientset,
		RestConfig: cfg,
		sessionHub: execpod.NewSessionHub(clientset, cfg),
	}
	return h
}

func (h *Handler) SSH(w http.ResponseWriter, r *http.Request, options *execpod.ExecOptions) {
	session, err := h.sessionHub.New(options)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	proxy, err := execpod.NewProxy(context.Background(), w, r)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	session.HandleSSH(proxy)
}

func (h *Handler) LogStreaming(w http.ResponseWriter, r *http.Request, options *execpod.ExecOptions) {
	session, err := h.sessionHub.New(options)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	proxy, err := execpod.NewProxy(context.Background(), w, r)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return
	}
	go session.HandleLog(proxy)
}
