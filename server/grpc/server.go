package grpc

import (
	"context"
	"crypto/tls"
	"github.com/CodeFish-xiao/go-fish-kit/internal/endpoint"
	"github.com/CodeFish-xiao/go-fish-kit/internal/host"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"time"

	"github.com/CodeFish-xiao/go-fish-kit/log"
	"github.com/CodeFish-xiao/go-fish-kit/server"
	"net"
	"net/url"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var (
	_ server.Server = (*Server)(nil)
)

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.ints = in
	}
}

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	baseCtx  context.Context
	tlsConf  *tls.Config
	lis      net.Listener
	once     sync.Once
	err      error
	network  string
	address  string
	endpoint *url.URL
	timeout  time.Duration
	health   *health.Server
	ints     []grpc.UnaryServerInterceptor
	log      *log.Helper
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		health:  health.NewServer(),
		log:     log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	grpcOpts := []grpc.ServerOption{}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)

	reflection.Register(srv.Server)
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.baseCtx = ctx
	s.log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.health.Shutdown()
	s.log.Info("[gRPC] server stopping")
	return nil
}
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, lis)
		if err != nil {
			_ = lis.Close()
			s.err = err
			return
		}
		s.lis = lis
		s.endpoint = endpoint.NewEndpoint("grpc", addr, s.tlsConf != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}
