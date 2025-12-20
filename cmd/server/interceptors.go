package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/easyp-tech/grpc-cource-2/pkg/auth"
)

func interceptorStat(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	// Pre-processing
	start := time.Now()

	// Вызов обработчика
	resp, err := handler(ctx, req)

	// Вычисляыем время выполенния
	duration := time.Since(start)
	if err != nil {
		log.Printf("[INTERCEPTOR STAT] %s failed after %v: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("[INTERCEPTOR STAT] %s completed in %v", info.FullMethod, duration)
	}

	return resp, err
}

func interceptorAuth(
	ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	keys := md.Get(auth.MDKeyName)
	if len(keys) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization")
	}

	if err := auth.ValidateAuthToken(keys[0]); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return handler(ctx, req)
}
