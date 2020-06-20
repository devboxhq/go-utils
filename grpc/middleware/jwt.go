package middleware

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// ValidateJwtFunc is a function that is provided to validate JWT request
type ValidateJwtFunc func(context context.Context, server interface{}, fullMethod string) (context.Context, error)

// jwtMiddleware is grpc middleware for JWT authentication
type jwtMiddleware struct {
	validationFunc ValidateJwtFunc
}

// NewJwtMiddleware creates jwtMiddleware
func NewJwtMiddleware(validationFunc ValidateJwtFunc) jwtMiddleware {
	return jwtMiddleware{
		validationFunc: validationFunc,
	}
}

func (m *jwtMiddleware) GetInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor, error) {
	unary := []grpc.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			newCtx, err := m.validationFunc(ctx, info.Server, info.FullMethod)
			if err != nil {
				return nil, fmt.Errorf("jwt auth unary middleware failed: %v", err)
			}

			return handler(newCtx, req)
		},
	}

	stream := []grpc.StreamServerInterceptor{
		func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			newCtx, err := m.validationFunc(stream.Context(), srv, info.FullMethod)
			if err != nil {
				return fmt.Errorf("jwt auth stream middleware failed: %v", err)
			}

			wrapped := grpc_middleware.WrapServerStream(stream)
			wrapped.WrappedContext = newCtx
			return handler(srv, wrapped)
		},
	}

	return unary, stream, nil
}
