package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
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

	s := grpc.NewServer()
	pb.RegisterNoteAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
