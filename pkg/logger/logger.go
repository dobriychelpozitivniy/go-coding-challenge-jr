package middlewareLogger

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func NewUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return NewUnaryServerInterceptorWithLogger(&log.Logger)
}

// Logger middleware for grpc. Logging request before handler and logging response after handler.
func NewUnaryServerInterceptorWithLogger(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		data := make([]byte, 10)

		_, _ = rand.Read(data)

		id := fmt.Sprintf("%x", sha1.Sum(data))

		log.Info().
			Interface("Request", req).
			Str("ID", id).
			Str("Handler", info.FullMethod).
			Send()

		resp, err := handler(ctx, req)
		if err != nil {
			err, _ := status.FromError(err)
			log.Error().
				Interface("Error", err.Message()).
				Str("ID", id).
				Str("Handler", info.FullMethod).
				Send()
		} else {
			log.Info().
				Interface("Response", resp).
				Str("ID", id).
				Str("Handler", info.FullMethod).
				Send()
		}

		return resp, err
	}
}
