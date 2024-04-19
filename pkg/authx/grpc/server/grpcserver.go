package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"k8s.io/klog/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpcapi "github.com/kzz45/neverdown/pkg/authx/grpc"
	pb "github.com/kzz45/neverdown/pkg/authx/grpc/proto"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

type GrpcApiOption struct {
	Certificate *tls.Certificate
	ListenPort  int32
	Secret      string
	Aggregator  grpcapi.Api
}

type Server struct {
	pb.AuthxApiServer
	hostname   string
	addr       string
	secret     string
	grpcServer *grpc.Server
	aggregator grpcapi.Api
}

func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}

func Init(opt *GrpcApiOption) *Server {
	grpcAddr := fmt.Sprintf(":%d", opt.ListenPort)
	s := &Server{
		// hostname:   env.GetHostNameMustSpecified(),
		hostname:   "authx-apiserver",
		addr:       grpcAddr,
		secret:     opt.Secret,
		aggregator: opt.Aggregator,
	}
	grpcOpts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(s.ensureValidToken),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.Creds(credentials.NewServerTLSFromCert(opt.Certificate)),
	}
	s.grpcServer = grpc.NewServer(grpcOpts...)
	return s
}

func (s *Server) DryRun() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		klog.Fatal("Failed to listen on addr: ", s.addr)
	}
	pb.RegisterAuthxApiServer(s.grpcServer, s)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			klog.Error(err)
		}
	}()
	klog.Info("Authority Grpc API is running...")
}

// valid validates the authorization.
func (s *Server) valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == s.secret
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func (s *Server) ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !s.valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
