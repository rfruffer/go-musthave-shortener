package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/rfruffer/go-musthave-shortener/cmd/shortener/router"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

func ExampleURLHandler_CreateShortJSONURLHandler() {
	repo := repository.NewInFileStore()
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service, "")

	r := router.SetupRouter(router.Router{URLHandler: handler})
	server := httptest.NewServer(r)
	defer server.Close()

	// handler.SetResultHost(server.URL)
	handler.SetResultHost(server.URL)

	// Отправляем POST-запрос
	reqBody := `{"url":"https://practicum.yandex.ru"}`
	resp, err := http.Post(server.URL+"/api/shorten", "application/json", strings.NewReader(reqBody))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var jsonResp models.ShortenResponse
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)

	shortURL := jsonResp.Result
	parts := strings.SplitN(shortURL, "/", 4)
	if len(parts) >= 3 {
		host := parts[2]
		// обрезаем порт, если есть
		if idx := strings.Index(host, ":"); idx != -1 {
			host = host[:idx]
		}
		shortURL = parts[0] + "//" + host
	}

	fmt.Println("Short URL:", shortURL)
	// Output:
	// Short URL: http://127.0.0.1
}
