package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
	"github.com/easyp-tech/grpc-cource-2/pkg/auth"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewNoteAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	md := auth.CreateClientMD()
	ctx = metadata.NewOutgoingContext(ctx, md)

	r, err := c.GetNote(ctx, &pb.NoteRequest{Id: "note_1"})
	if err != nil {
		log.Fatalf("could not get note: %v", err)
	}
	log.Printf("Response: %s", r)
}
