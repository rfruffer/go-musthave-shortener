package tests

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpcServer "github.com/rfruffer/go-musthave-shortener/internal/grpc"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	
	// Создаем test repository и services
	repo := repository.NewInFileStore()
	urlService := services.NewURLService(repo)
	commonService := services.NewCommonURLService(urlService, "http://localhost:8080")
	grpcURLServer := grpcServer.NewURLShortenerServer(commonService)
	
	grpcServer.RegisterURLShortenerServiceServer(s, grpcURLServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPCCreateShortURL(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcServer.NewURLShortenerServiceClient(conn)

	// Тест создания короткой ссылки
	req := &grpcServer.CreateShortURLRequest{
		Url:    "https://example.com",
		UserId: "test-user-123",
	}

	resp, err := client.CreateShortURL(ctx, req)
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}

	if resp.ShortUrl == "" {
		t.Error("Expected non-empty short URL")
	}

	if resp.AlreadyExists {
		t.Error("Expected AlreadyExists to be false for new URL")
	}

	t.Logf("Created short URL: %s", resp.ShortUrl)
}

func TestGRPCPing(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcServer.NewURLShortenerServiceClient(conn)

	// Тест ping
	req := &grpcServer.PingRequest{}

	resp, err := client.Ping(ctx, req)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	if !resp.Success {
		t.Error("Expected ping to be successful")
	}

	t.Log("Ping successful")
}

func TestGRPCBatch(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcServer.NewURLShortenerServiceClient(conn)

	// Тест пакетного создания
	req := &grpcServer.BatchRequest{
		UserId: "test-user-batch",
		Urls: []*grpcServer.BatchOriginalURL{
			{
				CorrelationId: "1",
				OriginalUrl:   "https://example1.com",
			},
			{
				CorrelationId: "2", 
				OriginalUrl:   "https://example2.com",
			},
		},
	}

	resp, err := client.Batch(ctx, req)
	if err != nil {
		t.Fatalf("Batch failed: %v", err)
	}

	if len(resp.Urls) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(resp.Urls))
	}

	for i, url := range resp.Urls {
		if url.ShortUrl == "" {
			t.Errorf("Expected non-empty short URL for item %d", i)
		}
		if url.CorrelationId == "" {
			t.Errorf("Expected non-empty correlation ID for item %d", i)
		}
		t.Logf("Batch item %d: %s -> %s", i, url.CorrelationId, url.ShortUrl)
	}
}