package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

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
	user, err := auth.GetUserFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	ctx = auth.PutUserToContext(ctx, user)

	return handler(ctx, req)
}

func interceptorStream(
	srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	log.Printf("Streaming %v", info.FullMethod)
	wrappedStream := &wrappedServerStream{
		ServerStream: ss,
	}

	return handler(srv, wrappedStream)
}

type wrappedServerStream struct {
	grpc.ServerStream
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("[SERVER] Intercepted server Send: %v\n", m)

	// Call the original Send method
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	// Call the original Rec method
	err := w.ServerStream.RecvMsg(m)
	log.Printf("[SERVER] Intercepted server Got: %v\n", m)

	return err
}
