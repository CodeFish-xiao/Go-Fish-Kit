package http

import (
	"context"
	"crypto/tls"
	"github.com/CodeFish-xiao/go-fish-kit/internal/endpoint"
	"github.com/CodeFish-xiao/go-fish-kit/internal/host"
	"github.com/CodeFish-xiao/go-fish-kit/log"
	"github.com/CodeFish-xiao/go-fish-kit/server"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

//为了易于扩展（不是），包装原生的http，所有的拦截器或者中间件由外界决定，舒服了

var (
	_ server.Server = (*Server)(nil)
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

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

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.ReadTimeout = timeout
		s.WriteTimeout = timeout
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.TLSConfig = c
	}
}

func Handler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.Handler = handler
	}
}

type Server struct {
	*http.Server

	once    sync.Once
	err     error
	network string
	address string

	endpoint *url.URL

	log *log.Helper

	lis net.Listener
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		log:     log.NewHelper(log.DefaultLogger),
	}
	srv.Server = &http.Server{}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	s.log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	var err error
	if s.TLSConfig != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}

// Endpoint return a real address to registry endpoint.
// examples:
//   http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, lis)
		if err != nil {
			lis.Close()
			s.err = err
			return
		}
		s.lis = lis

		s.endpoint = endpoint.NewEndpoint("http", addr, s.TLSConfig != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}
