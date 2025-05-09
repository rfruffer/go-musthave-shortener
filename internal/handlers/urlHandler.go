package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type URLHandler struct {
	service *services.URLService
	baseURL string
}

func NewURLHandler(service *services.URLService, baseURL string) *URLHandler {
	return &URLHandler{service: service, baseURL: baseURL}
}

func (us *URLHandler) CreateShortJSONURLHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "iempty or invalid body", http.StatusBadRequest)
		return
	}
	id, err := us.service.GenerateShortURL(req.URL)
	if err != nil {
		http.Error(w, "failed to create a short url", http.StatusNotFound)
		return
	}

	resp := models.ShortenResponse{
		Result: us.baseURL + "/" + id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
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
