package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (us *URLHandler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Empty or invalid body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	id, err := us.service.GenerateShortURL(originalURL)
	if err != nil {
		http.Error(w, "Failed to create a short url", http.StatusBadRequest)
		return
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	//shortURL := "http://localhost:8080/" + id
	shortURL := baseURL + "/" + id
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(shortURL))
}

func (us *URLHandler) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	originalURL, err := us.service.RedirectURL(id)
	if err != nil {
		http.Error(w, "Cant find id in store", http.StatusBadRequest)
		return
	}
	// w.Header().Set("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusTemporaryRedirect)
	// w.Write([]byte("Location: " + originalURL))
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
