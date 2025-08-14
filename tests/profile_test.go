package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func loadURL(t *testing.T, serverURL, original string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, serverURL+"/", bytes.NewBufferString(original))
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("cannot send request: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func runProfileTest(t *testing.T, profilePath string) {
	repo := repository.NewInFileStore()
	service := services.NewURLService(repo)
	shortURLHandler := handlers.NewURLHandler(service, "")

	r := router.SetupRouter(router.Router{
		URLHandler: shortURLHandler,
	})
	server := httptest.NewServer(r)
	defer server.Close()

	shortURLHandler.SetResultHost(server.URL)

	for i := 0; i < 5000; i++ {
		original := fmt.Sprintf("http://example.com/%s", strconv.Itoa(i))
		loadURL(t, server.URL, original)
	}

	time.Sleep(1 * time.Second)

	if err := os.MkdirAll("../profiles", 0755); err != nil {
		t.Fatalf("mkdir profiles: %v", err)
	}
	f, err := os.Create(profilePath)
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		t.Fatalf("write heap profile: %v", err)
	}
}

func TestProfileBase(t *testing.T) {
	runProfileTest(t, "../profiles/base.pprof")
}

func TestProfileResult(t *testing.T) {
	runProfileTest(t, "../profiles/result.pprof")
}
