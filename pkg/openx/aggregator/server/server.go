package server

import (
	"context"
	"net"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/kzz45/neverdown/pkg/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/websocket"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/apps"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/pods"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/resources"
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

	Resources *resources.Resources
}

type Server struct {
	resources   *resources.Resources
	members     *apps.Members
	connections *websocket.Connections
	podshandler *pods.Handler
	ctx         context.Context
	cancel      context.CancelFunc
}

func Run(opt *Options) (<-chan struct{}, <-chan struct{}, error) {
	subCtx, cancel := context.WithCancel(opt.Context)
	s := &Server{
		resources:   opt.Resources,
		members:     apps.NewMembers(subCtx, opt.Resources),
		connections: websocket.NewConnections(subCtx),
		podshandler: pods.New(opt.Resources),
		ctx:         subCtx,
		cancel:      cancel,
	}
	router := gin.New()
	df := cors.DefaultConfig()
	df.AllowAllOrigins = true
	df.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", jwttoken.TokenKey}
	router.Use(cors.New(df))
	s.registerRouters(router)
	if ok := env.PProfDebug(); ok {
		pprof.Register(router)
	}
	sec := server.SecureServingInfo{
		Listener: opt.Listener,
		SNICerts: opt.SNICerts,
	}
	return sec.Serve(router, time.Second*5, opt.StopChan)
}

func (s *Server) Shutdown() {
	s.cancel()
}
