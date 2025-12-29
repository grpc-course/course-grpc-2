package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewNoteAPIClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancel()
	ctx := context.Background()

	streamServer(ctx, c)
	//streamBidirectional(ctx, c)
}

func streamServer(ctx context.Context, c pb.NoteAPIClient) {
	log.Printf("Streaming server")

	req := &pb.Empty{}

	streamer, err := c.StreamNotes(ctx, req)
	if err != nil {
		log.Fatalf("could not create stream: %v", err)
	}

	for {
		if ctx.Err() != nil {
			log.Printf("context canceled")
			return
		}

		resp, err := streamer.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("server closed stream")
				return
			}
		}
		log.Printf("note: %v", resp)
	}
}
