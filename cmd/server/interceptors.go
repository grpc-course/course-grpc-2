package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
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
