package main

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

func (s *server) StreamNotesBidirectional(streamServer pb.NoteAPI_StreamNotesBidirectionalServer) error {
	log.Println("[SERVER] Starting bidirectional stream (async)")
	cntSent := 0

	ctx := streamServer.Context()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			req, err := streamServer.Recv()
			if err == io.EOF {
				log.Println("[SERVER] Client closed connection")
				return
			}
			if err != nil {
				log.Printf("[SERVER] Error receiving message: %v", err)
				return
			}

			log.Printf("[SERVER] Received message: %s", req.Id)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if ctx.Err() != nil {
				return
			}

			toSend := &pb.NoteResponse{
				Text: fmt.Sprintf("FROM SERVER: %v", cntSent),
			}
			cntSent++
			log.Printf("[SERVER] Sending message: %v", toSend)
			if err := streamServer.Send(toSend); err != nil {
				log.Printf("[SERVER] Error sending message: %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
		}
	}()

	wg.Wait()
	log.Println("[SERVER] Stream finished")
	return nil
}
