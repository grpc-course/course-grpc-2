package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	openapi "github.com/easyp-tech/grpc-cource-2/docs/api/notes/v1"
	pb "github.com/easyp-tech/grpc-cource-2/pkg/api/notes/v1"
	"github.com/easyp-tech/grpc-cource-2/pkg/auth"
)

const (
	keepaliveTime    = 50 * time.Second
	keepaliveTimeout = 10 * time.Second
	keepaliveMinTime = 30 * time.Second
)

type server struct {
	ctx context.Context
	cnt int // dummy
	pb.UnimplementedNoteAPIServer
}

func (s *server) GetNote(ctx context.Context, req *pb.NoteRequest) (*pb.NoteResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		//return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	log.Printf("Received note request for id: %s; user: %v", req.Id, user)

	now := time.Now()
	s.cnt++
	if s.cnt%2 == 0 {
		return s.WithError(ctx, req)
	}

	if s.cnt%3 == 0 {
		return nil, status.Error(codes.NotFound, "user has no notes")
	}

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

func (s *server) WithError(ctx context.Context, in *pb.NoteRequest) (*pb.NoteResponse, error) {
	// формируем кастомную ошибку
	st := status.New(codes.FailedPrecondition, "Custom error")
	errMsg := &pb.CustomError{Reason: pb.ErrorCode_ERROR_CODE_INVALID_TEXT}

	var err error
	// дополняем ее деталями: которые содержат структуру сообщения из proto файла.
	st, err = st.WithDetails(errMsg)
	if err != nil {
		return nil, err
	}

	return nil, st.Err()
}

func main() {
	wg := sync.WaitGroup{}

	ctx := context.Background()

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serverContext, cancel := context.WithCancel(ctx)

	ser := &server{
		ctx: serverContext,
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
		grpc.ChainStreamInterceptor(
			//interceptorStream,
		),
	)
	pb.RegisterNoteAPIServer(s, ser)

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Printf("grpc server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// init http
	//mux := runtime.NewServeMux()
	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = pb.RegisterNoteAPIHandlerFromEndpoint(ctx, gwMux, lis.Addr().String(), opts)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	serveSwagger(mux, openapi.Content)

	serverHttp := &http.Server{
		Addr:    ":8081",
		Handler: wsproxy.WebsocketProxy(mux),
	}

	wg.Add(1)
	mux.Handle("/api/", corsMiddleware(gwMux))
	go func() {
		defer wg.Done()

		log.Printf("http server listening at :8081")
		if err := serverHttp.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("failed to serve: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")

	cancel()

	if err := serverHttp.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}
	log.Printf("Http server is closed")

	s.GracefulStop()
	//s.Stop()

	log.Printf("Server exited properly")

	wg.Wait()
}
