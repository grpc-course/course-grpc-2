package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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
		log.Printf("could not get note: %v", err)

		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("status.FromError: %v", err)
		}
		log.Printf("Code: %s", st.Code().String())

		for _, d := range st.Details() {
			switch t := d.(type) {
			case *pb.CustomError:
				log.Printf("Reason: %v", t.Reason)
			}
		}

	}
	log.Printf("Response: %s", r)
}
