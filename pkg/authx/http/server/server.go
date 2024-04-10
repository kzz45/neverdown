package server

import (
	"context"
	"net"
	"time"

	"github.com/kzz45/neverdown/pkg/apiserver/server"
	"github.com/kzz45/neverdown/pkg/authx/rbac/handler"

	"github.com/kzz45/neverdown/pkg/jwttoken"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
)

type Options struct {
	// Listener is the secure server network listener.
	Listener net.Listener

	// Cert is the main server cert which is used if SNI does not match. Cert must be non-nil and is
	// allowed to be in SNICerts.
	Cert dynamiccertificates.CertKeyContentProvider

	// SNICerts are the TLS certificates used for SNI.
	SNICerts []dynamiccertificates.SNICertKeyContentProvider

	Context context.Context

	StopChan <-chan struct{}

	Handler *handler.Handler
}

type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	option  *Options
	handler *handler.Handler
}

func Init(opt *Options) *Server {
	subCtx, cancel := context.WithCancel(opt.Context)
	s := &Server{
		ctx:     subCtx,
		cancel:  cancel,
		option:  opt,
		handler: opt.Handler,
	}
	return s
}

func (s *Server) DryRun() (<-chan struct{}, <-chan struct{}, error) {
	router := gin.New()
	df := cors.DefaultConfig()
	df.AllowAllOrigins = true
	df.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", jwttoken.TokenKey}
	router.Use(cors.New(df))
	s.registerRouters(router)
	sec := server.SecureServingInfo{
		Listener: s.option.Listener,
		SNICerts: s.option.SNICerts,
		// Cert: s.option.Cert,
	}
	return sec.Serve(router, time.Second*5, s.option.StopChan)
}

func (s *Server) Shutdown() {
	s.cancel()
}
