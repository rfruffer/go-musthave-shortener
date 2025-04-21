package handlers

import (
	"io"
	"net/http"

	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type UrlHandler struct {
	service *services.UrlService
}

func NewUrlHandler(service *services.UrlService) *UrlHandler {
	return &UrlHandler{service: service}
}

func (us *UrlHandler) ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Empty or invalid body", http.StatusBadRequest)
			return
		}
		originalURL := string(body)
		id, err := us.service.GenerateShortUrl(originalURL)
		if err != nil {
			http.Error(w, "Failed to create a short url", http.StatusBadRequest)
			return
		}
		shortURL := "http://localhost:8080/" + id
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(shortURL))
	} else if r.Method == http.MethodGet {
		id := r.URL.Path[len("/"):]

		if id == "" {
			http.Error(w, "Missing ID", http.StatusBadRequest)
			return
		}
		originalURL, err := us.service.RedirectUrl(id)
		if err != nil {
			http.Error(w, "Cant find id in store", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)

		w.Write([]byte("Location: " + originalURL))

		return
	} else {
		http.Error(w, "Unsupport method", http.StatusUnauthorized)
		//w.WriteHeader(http.StatusUnauthorized)
	}
}
