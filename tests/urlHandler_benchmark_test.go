package tests

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func BenchmarkShortURLHandler(b *testing.B) {
	repo := repository.NewInFileStore()
	service := services.NewURLService(repo)
	shortURLhandler := handlers.NewURLHandler(service, "")

	r := router.SetupRouter(router.Router{
		URLHandler: shortURLhandler,
	})
	server := httptest.NewServer(r)
	defer server.Close()

	shortURLhandler.SetResultHost(server.URL)

	client := resty.New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		body := "https://practicum.yandex.ru/" + strings.Repeat("x", i%10)
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			SetBody(body).
			Post(server.URL + "/")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExpandURLHandler(b *testing.B) {
	repo := repository.NewInFileStore()
	service := services.NewURLService(repo)
	shortURLhandler := handlers.NewURLHandler(service, "")

	r := router.SetupRouter(router.Router{
		URLHandler: shortURLhandler,
	})
	server := httptest.NewServer(r)
	defer server.Close()

	shortURLhandler.SetResultHost(server.URL)

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody("https://practicum.yandex.ru/").
		Post(server.URL + "/")
	if err != nil {
		b.Fatal(err)
	}
	id := strings.TrimPrefix(string(resp.Body()), server.URL+"/")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.R().
			Get(server.URL + "/" + id)
		if err != nil {
			b.Fatal(err)
		}
	}
}
