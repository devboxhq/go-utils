package jwt

import (
	"context"

	jwt_go "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/uptimize/gg-execution/pkg/auth/jwt"
	"github.com/uptimize/gg-execution/pkg/grpc/middleware"
)

// JwtAuthService is an interface that is implemented to signal middleware to authenticate the request
type JwtAuthService interface {
	IsProtected() bool
}

// JwtAuthOverride is an interface that is implemented to overrride JWT authentication logic
type JwtAuthOverride interface {
	AuthImpl(manager jwt.Manager, token string) (jwt_go.Claims, error)
}

// JwtValidator validates JWT requests when used inside JWT middleware
func JwtMiddlewareValidator(manager jwt.Manager) middleware.ValidateJwtFunc {
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

				newCtx = context.WithValue(ctx, "user", claims)
			}
		} else {
			newCtx = ctx
		}

		return newCtx, err
	}
}
