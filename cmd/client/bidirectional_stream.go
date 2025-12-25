package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
)

func streamBidirectional(ctx context.Context, c pb.NoteAPIClient) {
	log.Printf("Streaming bidirectional server")
	cntSent := 0

	streamClient, err := c.StreamNotesBidirectional(ctx)
	if err != nil {
		log.Fatalf("failed to create async bidirectional stream: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Sender goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer streamClient.CloseSend()

		for {
			if ctx.Err() != nil {
				log.Printf("streaming cancelled")
				return
			}

			msg := fmt.Sprintf("sent %d", cntSent)
			if err := streamClient.Send(&pb.NoteRequest{Id: msg}); err != nil {
				errCh <- fmt.Errorf("failed to send async message %d: %w", cntSent, err)
				return
			}

			log.Printf("[Client] Sent async: %s", msg)
			time.Sleep(800 * time.Millisecond)
		}
	}()

	// Receiver goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			resp, err := streamClient.Recv()
			if err == io.EOF {
				log.Printf("[Client] Bidirectional async stream finished")
				return
			}
			if err != nil {
				errCh <- fmt.Errorf("failed to receive from async stream: %w", err)
				return
			}

			log.Printf("[Client] Async response: %s", resp.Text)
		}
	}()

	// Wait for completion or error
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case err := <-errCh:
		log.Fatal(err)
	case <-ctx.Done():
		log.Println("Bidirectional stream finished")
	}
}
