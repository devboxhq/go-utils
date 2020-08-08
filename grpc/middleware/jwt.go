package middleware

import (
	"context"
	"fmt"

	"github.com/devboxhq/go-utils/auth/jwt"
	jwt_go "github.com/dgrijalva/jwt-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

// ValidateJwtFunc is a function that is provided to validate JWT request
type ValidateJwtFunc func(context context.Context, server interface{}, fullMethod string) (context.Context, error)

// JwtAuthService is an interface that is implemented to signal middleware to authenticate the request
type JwtAuthService interface {
	IsProtected() bool
}

// JwtAuthOverride is an interface that is implemented to overrride JWT authentication logic
type JwtAuthOverride interface {
	AuthImpl(manager *jwt.Manager, token string) (jwt_go.Claims, error)
}

// JwtValidator validates JWT requests when used inside JWT middleware
func JwtMiddlewareValidator(manager *jwt.Manager) ValidateJwtFunc {
	var newCtx context.Context
	var claims jwt_go.Claims
	var err error

	return func(ctx context.Context, server interface{}, fullMethod string) (context.Context, error) {
		if serviceAuth, ok := server.(JwtAuthService); ok {
			if serviceAuth.IsProtected() {
				// get token from the authentication header
				authToken, err := grpc_auth.AuthFromMD(ctx, manager.GetHeaderScheme())
				if err != nil {
					return ctx, nil
				}

				if authOverride, ok := server.(JwtAuthOverride); ok {
					claims, err = authOverride.AuthImpl(manager, authToken)
				} else {
					claims, err = manager.Verify(authToken)
				}

				if err != nil {
					return ctx, err
				}
				newCtx = context.WithValue(ctx, "user", claims)
			}
		} else {
			newCtx = ctx
		}

		return newCtx, err
	}
}

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

func (m jwtMiddleware) GetInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor, error) {
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
