package middleware

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Zap grpc middleware
type zapMiddleware struct {
	logger  *zap.Logger
	verbose bool
}

// NewZapMiddleware creates middleware that logs server requests using zap
func NewZapMiddleware(logger *zap.Logger, verbose bool) zapMiddleware {
	return zapMiddleware{logger: logger, verbose: verbose}
}

func (m zapMiddleware) GetInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor, error) {
	if m.verbose {
		grpc_zap.ReplaceGrpcLogger(m.logger)
	}

	unary := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(m.logger),
	}

	stream := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(m.logger),
	}

	return unary, stream, nil
}
