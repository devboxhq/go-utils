package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// Middleware is a grpc middleware
type Middleware interface {
	GetInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor, error)
}

// Manager is a grpc middleware manager that builds server options
type Manager struct {
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (m *Manager) AddMiddleware(middleware Middleware) error {
	unary, stream, err := middleware.GetInterceptors()
	if err != nil {
		return err
	}

	m.addUnaryInterceptor(unary...)
	m.addStreamInterceptor(stream...)
	return nil
}

func (m *Manager) addUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	m.unaryInterceptors = append(m.unaryInterceptors, interceptors...)
}

func (m *Manager) addStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) {
	m.streamInterceptors = append(m.streamInterceptors, interceptor...)
}

func (m *Manager) BuildServerOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}

	opts = append(opts, grpc_middleware.WithUnaryServerChain(m.unaryInterceptors...))
	opts = append(opts, grpc_middleware.WithStreamServerChain(m.streamInterceptors...))
	return opts
}
