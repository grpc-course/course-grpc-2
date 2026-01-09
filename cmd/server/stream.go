package main

import (
	"fmt"
	"log"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

func (s *server) StreamNotes(req *pb.Empty, stream pb.NoteAPI_StreamNotesServer) error {
	log.Printf("StreamNotes called with req: %v", req)

	ctx := stream.Context()

	//for i := range 5 {
	i := 0
	for {
		i++
		if ctx.Err() != nil {
			log.Printf("StreamNotes cancelled: %v", ctx.Err())
			return ctx.Err()
		}

		if s.ctx.Err() != nil {
			log.Printf("StreamNotes cancelled: %v", s.ctx.Err())
			return ctx.Err()
		}

		//log.Printf("Sending note: %d", i)
		noteResp := &pb.NoteResponse{
			Id: fmt.Sprintf("note: %d", i),
		}
		if err := stream.Send(noteResp); err != nil {
			return err
		}
	}

	return nil
}
