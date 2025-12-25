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
	log.Println("EchoBidirectionalStreamAsync: Starting bidirectional stream (async)")
	cntSent := 0

	ctx := streamServer.Context()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			req, err := streamServer.Recv()
			if err == io.EOF {
				log.Println("EchoBidirectionalStreamAsync: Client closed connection")
				return
			}
			if err != nil {
				log.Printf("EchoBidirectionalStreamAsync: Error receiving message: %v", err)
				return
			}

			log.Printf("EchoBidirectionalStreamAsync: Received message: %s", req.Id)

			select {
			case <-ctx.Done():
				return
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
				Text: fmt.Sprintf("Note #%v", cntSent),
			}
			cntSent++
			log.Printf("EchoBidirectionalStreamAsync: Sending message: %v", toSend)
			if err := streamServer.Send(toSend); err != nil {
				log.Printf("EchoBidirectionalStreamAsync: Error sending message: %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
		}
	}()

	wg.Wait()
	log.Println("EchoBidirectionalStreamAsync: Stream finished")
	return nil
}
