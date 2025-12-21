package main

import (
	"context"
	"log"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

func (s *server) CreateNote(ctx context.Context, req *pb.NoteCreateRequest) (*pb.NoteCreateResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("CreateNote: %v", req)
	return &pb.NoteCreateResponse{}, nil
}
