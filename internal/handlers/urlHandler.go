package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type URLHandler struct {
	service *services.URLService
	baseURL string
}

func NewURLHandler(service *services.URLService, baseURL string) *URLHandler {
	return &URLHandler{service: service, baseURL: baseURL}
}

func (us *URLHandler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "empty or invalid body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	id, err := us.service.GenerateShortURL(originalURL)
	if err != nil {
		http.Error(w, "failed to create a short url", http.StatusNotFound)
		return
	}

	shortURL := us.baseURL + "/" + id
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(shortURL))
}

func (us *URLHandler) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing ID", http.StatusBadRequest)
		return
	}

	originalURL, err := us.service.RedirectURL(id)
	if err != nil {
		http.Error(w, "cant find id in store", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (us *URLHandler) SetResultHost(host string) {
	us.baseURL = host
}
