package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

const (
	keepaliveTime    = 50 * time.Second
	keepaliveTimeout = 10 * time.Second
	keepaliveMinTime = 30 * time.Second
)

type server struct {
	pb.UnimplementedNoteAPIServer
}

func (s *server) GetNote(_ context.Context, req *pb.NoteRequest) (*pb.NoteResponse, error) {
	log.Printf("Received note request for id: %s", req.Id)

	now := time.Now()

	createdAt := &datetime.DateTime{
		Year:    int32(now.Year()),
		Month:   int32(now.Month()),
		Day:     int32(now.Day()),
		Hours:   int32(now.Hour()),
		Minutes: int32(now.Minute()),
		Seconds: int32(now.Second()),
	}

	return &pb.NoteResponse{
		CreatedAt: createdAt,
		Id:        req.Id,
		Text:      "some_test",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создание gRPC сервера с параметрами
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{ //nolint:exhaustruct
				Time:    keepaliveTime,
				Timeout: keepaliveTimeout,
			},
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepaliveMinTime,
			PermitWithoutStream: true,
		}),
		// Создаем интерсепторы
		grpc.ChainUnaryInterceptor(
			interceptorStat,
			interceptorAuth,
		),
	)
	pb.RegisterNoteAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
