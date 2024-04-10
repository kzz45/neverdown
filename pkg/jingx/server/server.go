package server

import (
	"context"
	"net"
	"time"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/apps"

	"k8s.io/apiserver/pkg/server"
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

	AuthorityClientSet kubernetes.Interface
	// Api was the aggregator
	Api aggregator.Api

	StopChan <-chan struct{}
}

type Server struct {
	members     *apps.Members
	connections *websocket.Connections
	ctx         context.Context
	cancel      context.CancelFunc
}

func Run(opt *Options) (<-chan struct{}, <-chan struct{}, error) {
	subCtx, cancel := context.WithCancel(opt.Context)
	s := &Server{
		members:     apps.NewMembers(subCtx, opt.AuthorityClientSet, opt.Api),
		connections: websocket.NewConnections(subCtx),
		ctx:         subCtx,
		cancel:      cancel,
	}
	router := gin.New()
	df := cors.DefaultConfig()
	df.AllowAllOrigins = true
	df.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", jwttoken.TokenKey}
	router.Use(cors.New(df))
	s.registerRouters(router)
	sec := server.SecureServingInfo{
		Listener: opt.Listener,
		SNICerts: opt.SNICerts,
	}
	return sec.Serve(router, time.Second*5, opt.StopChan)
}

func (s *Server) Shutdown() {
	s.cancel()
}
